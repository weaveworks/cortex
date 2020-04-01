package api

import (
	"errors"
	"flag"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/prometheus/common/route"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"
	v1 "github.com/prometheus/prometheus/web/api/v1"
	"github.com/weaveworks/common/middleware"
	"github.com/weaveworks/common/server"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/cortexproject/cortex/pkg/alertmanager"
	"github.com/cortexproject/cortex/pkg/chunk/purger"
	"github.com/cortexproject/cortex/pkg/compactor"
	"github.com/cortexproject/cortex/pkg/distributor"
	"github.com/cortexproject/cortex/pkg/ingester"
	"github.com/cortexproject/cortex/pkg/ingester/client"
	"github.com/cortexproject/cortex/pkg/querier"
	"github.com/cortexproject/cortex/pkg/querier/frontend"
	"github.com/cortexproject/cortex/pkg/ring"
	"github.com/cortexproject/cortex/pkg/ruler"
	"github.com/cortexproject/cortex/pkg/storegateway"
	"github.com/cortexproject/cortex/pkg/util/push"
)

type Config struct {
	AlertmanagerHTTPPrefix string `yaml:"alertmanager_http_prefix"`
	PrometheusHTTPPrefix   string `yaml:"prometheus_http_prefix"`

	// The following configs are injected by the upstream caller.
	ServerPrefix       string          `yaml:"-"`
	LegacyHTTPPrefix   string          `yaml:"-"`
	HTTPAuthMiddleware middleware.Func `yaml:"-"`
}

// RegisterFlags adds the flags required to config this to the given FlagSet.
func (cfg *Config) RegisterFlags(f *flag.FlagSet) {
	cfg.RegisterFlagsWithPrefix("", f)
}

// RegisterFlagsWithPrefix adds the flags required to config this to the given FlagSet with the set prefix.
func (cfg *Config) RegisterFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	f.StringVar(&cfg.AlertmanagerHTTPPrefix, prefix+"http.alertmanager-http-prefix", "/alertmanager", "Base path for data storage.")
	f.StringVar(&cfg.PrometheusHTTPPrefix, prefix+"http.prometheus-http-prefix", "/prometheus", "Base path for data storage.")
}

type API struct {
	cfg              Config
	authMiddleware   middleware.Func
	server           *server.Server
	prometheusRouter *mux.Router
	logger           log.Logger
}

func New(cfg Config, s *server.Server, logger log.Logger) (*API, error) {
	// Ensure the encoded path is used. Required for the rules API
	s.HTTP.UseEncodedPath()

	api := &API{
		cfg:            cfg,
		authMiddleware: cfg.HTTPAuthMiddleware,
		server:         s,
		logger:         logger,
	}

	// If no authentication middleware is present in the config, use the middlewar
	if cfg.HTTPAuthMiddleware == nil {
		api.authMiddleware = middleware.AuthenticateUser
	}

	return api, nil
}

func (a *API) registerRoute(path string, handler http.Handler, auth bool, methods ...string) {
	level.Debug(a.logger).Log("msg", "api: registering route", "methods", strings.Join(methods, ","), "path", path, "auth", auth)
	if auth {
		handler = a.authMiddleware.Wrap(handler)
	}
	if len(methods) == 0 {
		a.server.HTTP.Path(path).Handler(handler)
		return
	}
	a.server.HTTP.Path(path).Methods(methods...).Handler(handler)
}

// Latest Prometheus requires r.RemoteAddr to be set to addr:port, otherwise it reject the request.
// Requests to Querier sometimes doesn't have that (if they are fetched from Query-Frontend).
// Prometheus uses this when logging queries to QueryLogger, but Cortex doesn't call engine.SetQueryLogger to set one.
//
// Can be removed when (if) https://github.com/prometheus/prometheus/pull/6840 is merged.
func fakeRemoteAddr(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RemoteAddr == "" {
			r.RemoteAddr = "127.0.0.1:8888"
		}
		handler.ServeHTTP(w, r)
	})
}

// RegisterAlertmanager registers endpoints associated with the alertmanager. It will only
// serve endpoints using the legacy http-prefix if it is not run as a single binary.
func (a *API) RegisterAlertmanager(am *alertmanager.MultitenantAlertmanager, target bool) {
	// Ensure this route is registered before the prefixed AM route
	a.registerRoute("/multitenant-alertmanager/status", am.GetStatusHandler(), false)

	// UI components lead to a large number of routes to support, utilize a path prefix instead
	a.server.HTTP.PathPrefix(a.cfg.AlertmanagerHTTPPrefix).Handler(a.authMiddleware.Wrap(am))
	level.Debug(a.logger).Log("msg", "api: registering alertmanager", "path_prefix", a.cfg.AlertmanagerHTTPPrefix)

	// If the target is Alertmanager, enable the legacy behaviour. Otherwise only enable
	// the component routed API.
	if target {
		a.registerRoute("/status", am.GetStatusHandler(), false)
		a.server.HTTP.PathPrefix(a.cfg.LegacyHTTPPrefix).Handler(a.authMiddleware.Wrap(am))
	}
}

// RegisterAPI registers the standard endpoints associated with a running Cortex.
func (a *API) RegisterAPI(cfg interface{}) {
	a.registerRoute("/config", configHandler(cfg), false)
	a.registerRoute("/", http.HandlerFunc(indexHandler), false)
}

// RegisterDistributor registers the endpoints associated with the distributor.
func (a *API) RegisterDistributor(d *distributor.Distributor, pushConfig distributor.Config) {
	a.registerRoute("/api/v1/push", push.Handler(pushConfig, d.Push), true)
	a.registerRoute("/distributor/all_user_stats", http.HandlerFunc(d.AllUserStatsHandler), false)
	a.registerRoute("/distributor/ha-tracker", d.Replicas, false)

	// Legacy Routes
	a.registerRoute(a.cfg.LegacyHTTPPrefix+"/push", push.Handler(pushConfig, d.Push), true)
	a.registerRoute("/all_user_stats", http.HandlerFunc(d.AllUserStatsHandler), false)
	a.registerRoute("/ha-tracker", d.Replicas, false)
}

// RegisterIngester registers the ingesters HTTP and GRPC service
func (a *API) RegisterIngester(i *ingester.Ingester, pushConfig distributor.Config) {
	client.RegisterIngesterServer(a.server.GRPC, i)
	grpc_health_v1.RegisterHealthServer(a.server.GRPC, i)

	a.registerRoute("/ingester/flush", http.HandlerFunc(i.FlushHandler), false)
	a.registerRoute("/ingester/shutdown", http.HandlerFunc(i.ShutdownHandler), false)
	a.registerRoute("/ingester/push", push.Handler(pushConfig, i.Push), true) // For testing and debugging.

	// Legacy Routes
	a.registerRoute("/flush", http.HandlerFunc(i.FlushHandler), false)
	a.registerRoute("/shutdown", http.HandlerFunc(i.ShutdownHandler), false)
	a.registerRoute("/push", push.Handler(pushConfig, i.Push), true) // For testing and debugging.
}

// RegisterPurger registers the endpoints associated with the Purger/DeleteStore. They do not exacty
// match the Prometheus API but mirror it closely enough to justify their routing under the Prometheus
// component/
func (a *API) RegisterPurger(store *purger.DeleteStore) {
	deleteRequestHandler := purger.NewDeleteRequestHandler(store)

	a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/admin/tsdb/delete_series", http.HandlerFunc(deleteRequestHandler.AddDeleteRequestHandler), true, "PUT", "POST")
	a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/admin/tsdb/delete_series", http.HandlerFunc(deleteRequestHandler.GetAllDeleteRequestsHandler), true, "GET")

	// Legacy Routes
	a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/admin/tsdb/delete_series", http.HandlerFunc(deleteRequestHandler.AddDeleteRequestHandler), true, "PUT", "POST")
	a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/admin/tsdb/delete_series", http.HandlerFunc(deleteRequestHandler.GetAllDeleteRequestsHandler), true, "GET")
}

// RegisterRuler registers routes associated with the Ruler service. If the
// API is not enabled only the ring route is registered.
func (a *API) RegisterRuler(r *ruler.Ruler, apiEnabled bool) {
	a.registerRoute("/ruler/ring", r, false)

	if apiEnabled {
		a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/rules", http.HandlerFunc(r.PrometheusRules), true, "GET")
		a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/alerts", http.HandlerFunc(r.PrometheusAlerts), true, "GET")

		ruler.RegisterRulerServer(a.server.GRPC, r)

		a.registerRoute("/api/v1/rules", http.HandlerFunc(r.ListRules), true, "GET")
		a.registerRoute("/api/v1/rules/{namespace}", http.HandlerFunc(r.ListRules), true, "GET")
		a.registerRoute("/api/v1/rules/{namespace}/{groupName}", http.HandlerFunc(r.GetRuleGroup), true, "GET")
		a.registerRoute("/api/v1/rules/{namespace}/{groupName}", http.HandlerFunc(r.GetRuleGroup), true, "GET")
		a.registerRoute("/api/v1/rules/{namespace}", http.HandlerFunc(r.CreateRuleGroup), true, "POST")
		a.registerRoute("/api/v1/rules/{namespace}/{groupName}", http.HandlerFunc(r.DeleteRuleGroup), true, "DELETE")

		a.registerRoute(a.cfg.LegacyHTTPPrefix+"/api/v1/rules", http.HandlerFunc(r.PrometheusRules), true, "GET")
		a.registerRoute(a.cfg.LegacyHTTPPrefix+"/api/v1/alerts", http.HandlerFunc(r.PrometheusAlerts), true, "GET")

		a.registerRoute(a.cfg.LegacyHTTPPrefix+"/rules", http.HandlerFunc(r.ListRules), true, "GET")
		a.registerRoute(a.cfg.LegacyHTTPPrefix+"/rules/{namespace}", http.HandlerFunc(r.ListRules), true, "GET")
		a.registerRoute(a.cfg.LegacyHTTPPrefix+"/rules/{namespace}/{groupName}", http.HandlerFunc(r.GetRuleGroup), true, "GET")
		a.registerRoute(a.cfg.LegacyHTTPPrefix+"/rules/{namespace}/{groupName}", http.HandlerFunc(r.GetRuleGroup), true, "GET")
		a.registerRoute(a.cfg.LegacyHTTPPrefix+"/rules/{namespace}", http.HandlerFunc(r.CreateRuleGroup), true, "POST")
		a.registerRoute(a.cfg.LegacyHTTPPrefix+"/rules/{namespace}/{groupName}", http.HandlerFunc(r.DeleteRuleGroup), true, "DELETE")
	}
}

// // RegisterRing registers the ring UI page associated with the distributor for writes.
func (a *API) RegisterRing(r *ring.Ring) {
	a.registerRoute("/ring", r, false)
}

// RegisterStoreGateway registers the ring UI page associated with the store-gateway.
func (a *API) RegisterStoreGateway(s *storegateway.StoreGateway) {
	a.registerRoute("/store-gateway/ring", http.HandlerFunc(s.RingHandler), false)
}

// RegisterCompactor registers the ring UI page associated with the compactor.
func (a *API) RegisterCompactor(c *compactor.Compactor) {
	a.registerRoute("/compactor/ring", http.HandlerFunc(c.RingHandler), false)
}

// RegisterQuerier registers the Prometheus routes supported by the
// Cortex querier service. Currently this can not be registered simultaneously
// with the QueryFrontend.
func (a *API) RegisterQuerier(queryable storage.Queryable, engine *promql.Engine, distributor *distributor.Distributor) {
	api := v1.NewAPI(
		engine,
		queryable,
		querier.DummyTargetRetriever{},
		querier.DummyAlertmanagerRetriever{},
		func() config.Config { return config.Config{} },
		map[string]string{}, // TODO: include configuration flags
		func(f http.HandlerFunc) http.HandlerFunc { return f },
		func() v1.TSDBAdmin { return nil }, // Only needed for admin APIs.
		false,                              // Disable admin APIs.
		a.logger,
		querier.DummyRulesRetriever{},
		0, 0, 0, // Remote read samples and concurrency limit.
		regexp.MustCompile(".*"),
		func() (v1.RuntimeInfo, error) { return v1.RuntimeInfo{}, errors.New("not implemented") },
		&v1.PrometheusVersion{},
	)

	promRouter := route.New().WithPrefix(a.cfg.ServerPrefix + a.cfg.PrometheusHTTPPrefix + "/api/v1")
	api.Register(promRouter)
	promHandler := fakeRemoteAddr(promRouter)

	a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/read", querier.RemoteReadHandler(queryable), true)
	a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/query", promHandler, true)
	a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/query_range", promHandler, true)
	a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/labels", promHandler, true)
	a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/label/{name}/values", promHandler, true)
	a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/series", promHandler, true)
	a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/metadata", promHandler, true)

	a.registerRoute("/api/v1/user_stats", http.HandlerFunc(distributor.UserStatsHandler), true)
	a.registerRoute("/api/v1/chunks", querier.ChunksHandler(queryable), true)

	// Legacy Routes
	a.registerRoute("/user_stats", http.HandlerFunc(distributor.UserStatsHandler), true)
	a.registerRoute("/chunks", querier.ChunksHandler(queryable), true)

	legacyPromRouter := route.New().WithPrefix(a.cfg.ServerPrefix + a.cfg.LegacyHTTPPrefix + "/api/v1")
	api.Register(legacyPromRouter)
	legacyPromHandler := fakeRemoteAddr(legacyPromRouter)

	a.registerRoute(a.cfg.LegacyHTTPPrefix+"/api/v1/read", querier.RemoteReadHandler(queryable), true)
	a.registerRoute(a.cfg.LegacyHTTPPrefix+"/api/v1/query", legacyPromHandler, true)
	a.registerRoute(a.cfg.LegacyHTTPPrefix+"/api/v1/query_range", legacyPromHandler, true)
	a.registerRoute(a.cfg.LegacyHTTPPrefix+"/api/v1/labels", legacyPromHandler, true)
	a.registerRoute(a.cfg.LegacyHTTPPrefix+"/api/v1/label/{name}/values", legacyPromHandler, true)
	a.registerRoute(a.cfg.LegacyHTTPPrefix+"/api/v1/series", legacyPromHandler, true)
	a.registerRoute(a.cfg.LegacyHTTPPrefix+"/api/v1/metadata", legacyPromHandler, true)
}

// RegisterQueryFrontend registers the Prometheus routes supported by the
// Cortex querier service. Currently this can not be registered simultaneously
// with the Querier.
func (a *API) RegisterQueryFrontend(f *frontend.Frontend) {
	frontend.RegisterFrontendServer(a.server.GRPC, f)

	// Previously the frontend handled all calls to the provided prefix. Instead explicit
	// routing is used since it will be required to enable the frontend to be run as part
	// of a single binary in the future.
	a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/read", f.Handler(), true)
	a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/query", f.Handler(), true)
	a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/query_range", f.Handler(), true)
	a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/labels", f.Handler(), true)
	a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/label/{name}/values", f.Handler(), true)
	a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/series", f.Handler(), true)
	a.registerRoute(a.cfg.PrometheusHTTPPrefix+"/api/v1/metadata", f.Handler(), true)

	// Register Legacy Routers
	a.registerRoute(a.cfg.LegacyHTTPPrefix+"/api/v1/read", f.Handler(), true)
	a.registerRoute(a.cfg.LegacyHTTPPrefix+"/api/v1/query", f.Handler(), true)
	a.registerRoute(a.cfg.LegacyHTTPPrefix+"/api/v1/query_range", f.Handler(), true)
	a.registerRoute(a.cfg.LegacyHTTPPrefix+"/api/v1/labels", f.Handler(), true)
	a.registerRoute(a.cfg.LegacyHTTPPrefix+"/api/v1/label/{name}/values", f.Handler(), true)
	a.registerRoute(a.cfg.LegacyHTTPPrefix+"/api/v1/series", f.Handler(), true)
	a.registerRoute(a.cfg.LegacyHTTPPrefix+"/api/v1/metadata", f.Handler(), true)
}

// RegisterServiceMapHandler registers the Cortex structs service handler
// TODO: Refactor this code to be accomplished using the services.ServiceManager
// or a future module manager #2291
func (a *API) RegisterServiceMapHandler(handler http.Handler) {
	a.registerRoute("/services", handler, false)
}
