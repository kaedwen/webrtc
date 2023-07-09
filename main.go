package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/kaedwen/webrtc/pkg/common"
	"github.com/kaedwen/webrtc/pkg/server"
	"github.com/kaedwen/webrtc/pkg/webrtc"
	"github.com/tinyzimmer/go-glib/glib"
	"go.uber.org/zap"
)

func main() {
	mainLoop := glib.NewMainLoop(glib.MainContextDefault(), true)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := common.Config{}
	cfg.MustParse()

	lg, err := common.NewLogger(&cfg.Logging)
	if err != nil {
		panic(err)
	}

	http := server.NewHttpServer(lg.With(zap.String("context", "server")), &cfg)

	err = webrtc.NewWebrtcHandler(ctx, lg.With(zap.String("context", "webrtc")), cfg.Stream(), http.Hndl)
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(ctx, cfg.Http.Address())

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	lg.Info("shutting down gracefully, press Ctrl+C again to force")

	// stop gst main loop
	mainLoop.Quit()

	if err := http.TearDown(); err != nil {
		lg.Error("timeout", zap.Error(err))
	}
}
