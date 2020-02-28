package cortex

import (
	"context"

	"github.com/go-kit/kit/log/level"
	"github.com/weaveworks/common/server"

	"github.com/cortexproject/cortex/pkg/util"
	"github.com/cortexproject/cortex/pkg/util/services"
)

// NewServerService constructs service from Server component.
// servicesToWaitFor is called when server is stopping, and should return all
// services that need to terminate before server actually stops.
func NewServerService(serv *server.Server, servicesToWaitFor func() []services.Service) services.Service {
	serverDone := make(chan error, 1)

	runFn := func(ctx context.Context) error {
		go func() {
			defer close(serverDone)
			serverDone <- serv.Run()
		}()

		select {
		case <-ctx.Done():
			return nil
		case err := <-serverDone:
			if err != nil {
				level.Error(util.Logger).Log("msg", "server failed", "err", err)
			}
			return err
		}
	}

	stoppingFn := func() error {
		// wait until all modules are done, and then shutdown server.
		for _, s := range servicesToWaitFor() {
			_ = s.AwaitTerminated(context.Background())
		}

		// unblock Run, if it's still running (e.g. service was asked to stop via StopAsync)
		serv.Stop()

		// shutdown HTTP and gRPC servers
		serv.Shutdown()

		// if not closed yet, wait until server stops.
		<-serverDone
		level.Info(util.Logger).Log("msg", "server stopped")
		return nil
	}

	return services.NewBasicService(nil, runFn, stoppingFn)
}
