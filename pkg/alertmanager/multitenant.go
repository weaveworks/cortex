package alertmanager

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/prometheus/alertmanager/cluster"
	amconfig "github.com/prometheus/alertmanager/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/weaveworks/common/user"

	"github.com/cortexproject/cortex/pkg/alertmanager/alerts"
	"github.com/cortexproject/cortex/pkg/util"
	"github.com/cortexproject/cortex/pkg/util/flagext"
)

var backoffConfig = util.BackoffConfig{
	// Backoff for loading initial configuration set.
	MinBackoff: 100 * time.Millisecond,
	MaxBackoff: 2 * time.Second,
}

const (
	// If a config sets the webhook URL to this, it will be rewritten to
	// a URL derived from Config.AutoWebhookRoot
	autoWebhookURL = "http://internal.monitor"

	statusPage = `
<!doctype html>
<html>
	<head><title>Cortex Alertmanager Status</title></head>
	<body>
		<h1>Cortex Alertmanager Status</h1>
		<h2>Node</h2>
		<dl>
			<dt>Name</dt><dd>{{.self.Name}}</dd>
			<dt>Addr</dt><dd>{{.self.Addr}}</dd>
			<dt>Port</dt><dd>{{.self.Port}}</dd>
		</dl>
		<h3>Members</h3>
		{{ with .members }}
		<table>
		<tr><th>Name</th><th>Addr</th></tr>
		{{ range . }}
		<tr><td>{{ .Name }}</td><td>{{ .Addr }}</td></tr>
		{{ end }}
		</table>
		{{ else }}
		<p>No peers</p>
		{{ end }}
	</body>
</html>
`
)

var (
	totalConfigs = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "cortex",
		Name:      "alertmanager_configs",
		Help:      "How many configs the multitenant alertmanager knows about.",
	})
	statusTemplate *template.Template
)

func init() {
	prometheus.MustRegister(totalConfigs)
	statusTemplate = template.Must(template.New("statusPage").Funcs(map[string]interface{}{
		"state": func(enabled bool) string {
			if enabled {
				return "enabled"
			}
			return "disabled"
		},
	}).Parse(statusPage))
}

// MultitenantAlertmanagerConfig is the configuration for a multitenant Alertmanager.
type MultitenantAlertmanagerConfig struct {
	DataDir      string
	Retention    time.Duration
	ExternalURL  flagext.URLValue
	PollInterval time.Duration

	ClusterBindAddr      string
	ClusterAdvertiseAddr string
	Peers                flagext.StringSlice
	PeerTimeout          time.Duration

	FallbackConfigFile string
	AutoWebhookRoot    string

	Store AlertStoreConfig
}

const defaultClusterAddr = "0.0.0.0:9094"

// RegisterFlags adds the flags required to config this to the given FlagSet.
func (cfg *MultitenantAlertmanagerConfig) RegisterFlags(f *flag.FlagSet) {
	f.StringVar(&cfg.DataDir, "alertmanager.storage.path", "data/", "Base path for data storage.")
	f.DurationVar(&cfg.Retention, "alertmanager.storage.retention", 5*24*time.Hour, "How long to keep data for.")

	f.Var(&cfg.ExternalURL, "alertmanager.web.external-url", "The URL under which Alertmanager is externally reachable (for example, if Alertmanager is served via a reverse proxy). Used for generating relative and absolute links back to Alertmanager itself. If the URL has a path portion, it will be used to prefix all HTTP endpoints served by Alertmanager. If omitted, relevant URL components will be derived automatically.")

	f.StringVar(&cfg.FallbackConfigFile, "alertmanager.configs.fallback", "", "Filename of fallback config to use if none specified for instance.")
	f.StringVar(&cfg.AutoWebhookRoot, "alertmanager.configs.auto-webhook-root", "", "Root of URL to generate if config is "+autoWebhookURL)
	f.DurationVar(&cfg.PollInterval, "alertmanager.configs.poll-interval", 15*time.Second, "How frequently to poll Cortex configs")

	f.StringVar(&cfg.ClusterBindAddr, "cluster.listen-address", defaultClusterAddr, "Listen address for cluster.")
	f.StringVar(&cfg.ClusterAdvertiseAddr, "cluster.advertise-address", "", "Explicit address to advertise in cluster.")
	f.Var(&cfg.Peers, "cluster.peer", "Initial peers (may be repeated).")
	f.DurationVar(&cfg.PeerTimeout, "cluster.peer-timeout", time.Second*15, "Time to wait between peers to send notifications.")

	cfg.Store.RegisterFlags(f)
}

// A MultitenantAlertmanager manages Alertmanager instances for multiple
// organizations.
type MultitenantAlertmanager struct {
	cfg *MultitenantAlertmanagerConfig

	store alerts.AlertStore

	// The fallback config is stored as a string and parsed every time it's needed
	// because we mutate the parsed results and don't want those changes to take
	// effect here.
	fallbackConfig string

	// All the organization configurations that we have. Only used for instrumentation.
	cfgs map[string]alerts.AlertConfigDesc

	alertmanagersMtx sync.Mutex
	alertmanagers    map[string]*Alertmanager

	peer *cluster.Peer

	stop chan struct{}
	done chan struct{}
}

// NewMultitenantAlertmanager creates a new MultitenantAlertmanager.
func NewMultitenantAlertmanager(cfg *MultitenantAlertmanagerConfig) (*MultitenantAlertmanager, error) {
	err := os.MkdirAll(cfg.DataDir, 0777)
	if err != nil {
		return nil, fmt.Errorf("unable to create Alertmanager data directory %q: %s", cfg.DataDir, err)
	}

	var fallbackConfig []byte
	if cfg.FallbackConfigFile != "" {
		fallbackConfig, err = ioutil.ReadFile(cfg.FallbackConfigFile)
		if err != nil {
			return nil, fmt.Errorf("unable to read fallback config %q: %s", cfg.FallbackConfigFile, err)
		}
		_, err = amconfig.LoadFile(cfg.FallbackConfigFile)
		if err != nil {
			return nil, fmt.Errorf("unable to load fallback config %q: %s", cfg.FallbackConfigFile, err)
		}
	}

	var peer *cluster.Peer
	if cfg.ClusterBindAddr != "" {
		peer, err = cluster.Create(
			log.With(util.Logger, "component", "cluster"),
			prometheus.DefaultRegisterer,
			cfg.ClusterBindAddr,
			cfg.ClusterAdvertiseAddr,
			cfg.Peers,
			true,
			cluster.DefaultPushPullInterval,
			cluster.DefaultGossipInterval,
			cluster.DefaultTcpTimeout,
			cluster.DefaultProbeTimeout,
			cluster.DefaultProbeInterval,
		)
		if err != nil {
			return nil, errors.Wrap(err, "unable to initialize gossip mesh")
		}
		err = peer.Join(cluster.DefaultReconnectInterval, cluster.DefaultReconnectTimeout)
		if err != nil {
			level.Warn(util.Logger).Log("msg", "unable to join gossip mesh", "err", err)
		}
		go peer.Settle(context.Background(), cluster.DefaultGossipInterval)
	}

	store, err := NewAlertStore(cfg.Store)
	if err != nil {
		return nil, err
	}

	am := &MultitenantAlertmanager{
		cfg:            cfg,
		fallbackConfig: string(fallbackConfig),
		cfgs:           map[string]alerts.AlertConfigDesc{},
		alertmanagers:  map[string]*Alertmanager{},
		peer:           peer,
		store:          store,
		stop:           make(chan struct{}),
		done:           make(chan struct{}),
	}
	return am, nil
}

// Run the MultitenantAlertmanager.
func (am *MultitenantAlertmanager) Run() {
	defer close(am.done)

	// Load initial set of all configurations before polling for new ones.
	am.addNewConfigs(am.loadAllConfigs())
	ticker := time.NewTicker(am.cfg.PollInterval)
	for {
		select {
		case now := <-ticker.C:
			err := am.updateConfigs(now)
			if err != nil {
				level.Warn(util.Logger).Log("msg", "MultitenantAlertmanager: error updating configs", "err", err)
			}
		case <-am.stop:
			ticker.Stop()
			return
		}
	}
}

// Stop stops the MultitenantAlertmanager.
func (am *MultitenantAlertmanager) Stop() {
	close(am.stop)
	<-am.done
	am.alertmanagersMtx.Lock()
	for _, am := range am.alertmanagers {
		am.Stop()
	}
	am.alertmanagersMtx.Unlock()
	err := am.peer.Leave(am.cfg.PeerTimeout)
	if err != nil {
		level.Warn(util.Logger).Log("msg", "MultitenantAlertmanager: failed to leave the cluster", "err", err)
	}
	level.Debug(util.Logger).Log("msg", "MultitenantAlertmanager stopped")
}

// Load the full set of configurations from the server, retrying with backoff
// until we can get them.
func (am *MultitenantAlertmanager) loadAllConfigs() map[string]alerts.AlertConfigDesc {
	backoff := util.NewBackoff(context.Background(), backoffConfig)
	for {
		cfgs, err := am.poll()
		if err == nil {
			level.Debug(util.Logger).Log("msg", "MultitenantAlertmanager: initial configuration load", "num_configs", len(cfgs))
			return cfgs
		}
		level.Warn(util.Logger).Log("msg", "MultitenantAlertmanager: error fetching all configurations, backing off", "err", err)
		backoff.Wait()
	}
}

func (am *MultitenantAlertmanager) updateConfigs(now time.Time) error {
	cfgs, err := am.poll()
	if err != nil {
		return err
	}
	am.addNewConfigs(cfgs)
	return nil
}

// poll the configuration server. Not re-entrant.
func (am *MultitenantAlertmanager) poll() (map[string]alerts.AlertConfigDesc, error) {
	cfgs, err := am.store.ListAlertConfigs(context.Background())
	if err != nil {
		level.Warn(util.Logger).Log("msg", "MultitenantAlertmanager: configs server poll failed", "err", err)
		return nil, err
	}
	return cfgs, nil
}

func (am *MultitenantAlertmanager) addNewConfigs(cfgs map[string]alerts.AlertConfigDesc) {
	// TODO: instrument how many configs we have, both valid & invalid.
	level.Debug(util.Logger).Log("msg", "adding configurations", "num_configs", len(cfgs))
	for _, cfg := range cfgs {
		err := am.setConfig(cfg)
		if err != nil {
			level.Warn(util.Logger).Log("msg", "MultitenantAlertmanager: error applying config", "err", err)
			continue
		}
	}

	am.alertmanagersMtx.Lock()
	defer am.alertmanagersMtx.Unlock()
	for user, userAM := range am.alertmanagers {
		if _, exists := am.alertmanagers[user]; !exists {
			go userAM.Stop()
			delete(am.alertmanagers, user)
			delete(am.cfgs, user)
			level.Info(util.Logger).Log("msg", "deleting alertmanager", "user", user)
		}
	}
	totalConfigs.Set(float64(len(am.alertmanagers)))
}

func (am *MultitenantAlertmanager) transformConfig(userID string, amConfig *amconfig.Config) (*amconfig.Config, error) {
	if amConfig == nil { // shouldn't happen, but check just in case
		return nil, fmt.Errorf("no usable Cortex configuration for %v", userID)
	}
	if am.cfg.AutoWebhookRoot != "" {
		for _, r := range amConfig.Receivers {
			for _, w := range r.WebhookConfigs {
				if w.URL.String() == autoWebhookURL {
					u, err := url.Parse(am.cfg.AutoWebhookRoot + "/" + userID + "/monitor")
					if err != nil {
						return nil, err
					}
					w.URL = &amconfig.URL{URL: u}
				}
			}
		}
	}

	return amConfig, nil
}

func (am *MultitenantAlertmanager) createTemplatesFile(userID, fn, content string) (bool, error) {
	dir := filepath.Join(am.cfg.DataDir, "templates", userID, filepath.Dir(fn))
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return false, fmt.Errorf("unable to create Alertmanager templates directory %q: %s", dir, err)
	}

	file := filepath.Join(dir, fn)
	// Check if the template file already exists and if it has changed
	if tmpl, err := ioutil.ReadFile(file); err == nil && string(tmpl) == content {
		return false, nil
	}

	if err := ioutil.WriteFile(file, []byte(content), 0644); err != nil {
		return false, fmt.Errorf("unable to create Alertmanager template file %q: %s", file, err)
	}

	return true, nil
}

// setConfig applies the given configuration to the alertmanager for `userID`,
// creating an alertmanager if it doesn't already exist.
func (am *MultitenantAlertmanager) setConfig(cfg alerts.AlertConfigDesc) error {
	am.alertmanagersMtx.Lock()
	existing, hasExisting := am.alertmanagers[cfg.User]
	am.alertmanagersMtx.Unlock()
	var userAmConfig *amconfig.Config
	var err error
	var hasTemplateChanges bool

	for _, tmpl := range cfg.Templates {
		hasChanged, err := am.createTemplatesFile(cfg.User, tmpl.Filename, tmpl.Body)
		if err != nil {
			return err
		}

		if hasChanged {
			hasTemplateChanges = true
		}
	}

	level.Debug(util.Logger).Log("msg", "MultitenantAlertmanager: setting config", "user", cfg.User)

	if cfg.RawConfig == "" {
		if am.fallbackConfig == "" {
			return fmt.Errorf("blank Alertmanager configuration for %v", cfg.User)
		}
		level.Info(util.Logger).Log("msg", "blank Alertmanager configuration; using fallback", "user_id", cfg.User)
		userAmConfig, err = amconfig.Load(am.fallbackConfig)
		if err != nil {
			return fmt.Errorf("unable to load fallback configuration for %v: %v", cfg.User, err)
		}
	} else {
		userAmConfig, err = amconfig.Load(cfg.RawConfig)
		if err != nil && hasExisting {
			// XXX: This means that if a user has a working configuration and
			// they submit a broken one, we'll keep processing the last known
			// working configuration, and they'll never know.
			// TODO: Provide a way of communicating this to the user and for removing
			// Alertmanager instances.
			return fmt.Errorf("invalid Cortex configuration for %v: %v", cfg.User, err)
		}
	}

	if userAmConfig, err = am.transformConfig(cfg.User, userAmConfig); err != nil {
		return err
	}

	// If no Alertmanager instance exists for this user yet, start one.
	if !hasExisting {
		level.Debug(util.Logger).Log("msg", "MultitenantAlertmanager: initializing new alertmanager tenant", "user", cfg.User)
		newAM, err := am.newAlertmanager(cfg.User, userAmConfig)
		if err != nil {
			return err
		}
		am.alertmanagersMtx.Lock()
		am.alertmanagers[cfg.User] = newAM
		am.alertmanagersMtx.Unlock()
	} else if am.cfgs[cfg.User].RawConfig != cfg.RawConfig || hasTemplateChanges {
		level.Debug(util.Logger).Log("msg", "MultitenantAlertmanager: updating new alertmanager tenant", "user", cfg.User)
		// If the config changed, apply the new one.
		err := existing.ApplyConfig(cfg.User, userAmConfig)
		if err != nil {
			return fmt.Errorf("unable to apply Alertmanager config for user %v: %v", cfg.User, err)
		}
	}
	am.cfgs[cfg.User] = cfg
	return nil
}

func (am *MultitenantAlertmanager) newAlertmanager(userID string, amConfig *amconfig.Config) (*Alertmanager, error) {
	newAM, err := New(&Config{
		UserID:      userID,
		DataDir:     am.cfg.DataDir,
		Logger:      util.Logger,
		Peer:        am.peer,
		PeerTimeout: am.cfg.PeerTimeout,
		Retention:   am.cfg.Retention,
		ExternalURL: am.cfg.ExternalURL.URL,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to start Alertmanager for user %v: %v", userID, err)
	}

	if err := newAM.ApplyConfig(userID, amConfig); err != nil {
		return nil, fmt.Errorf("unable to apply initial config for user %v: %v", userID, err)
	}
	return newAM, nil
}

// ServeHTTP serves the Alertmanager's web UI and API.
func (am *MultitenantAlertmanager) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	userID, _, err := user.ExtractOrgIDFromHTTPRequest(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	am.alertmanagersMtx.Lock()
	userAM, ok := am.alertmanagers[userID]
	am.alertmanagersMtx.Unlock()
	if !ok {
		http.Error(w, fmt.Sprintf("no Alertmanager for this user ID"), http.StatusNotFound)
		return
	}
	userAM.mux.ServeHTTP(w, req)
}

// GetStatusHandler returns the status handler for this multi-tenant
// alertmanager.
func (am *MultitenantAlertmanager) GetStatusHandler() StatusHandler {
	return StatusHandler{
		am: am,
	}
}

// StatusHandler shows the status of the alertmanager.
type StatusHandler struct {
	am *MultitenantAlertmanager
}

// ServeHTTP serves the status of the alertmanager.
func (s StatusHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := statusTemplate.Execute(w, s.am.peer.Info())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
