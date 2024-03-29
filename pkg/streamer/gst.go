package streamer

import (
	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
	"go.uber.org/zap"
)

func init() {
	gst.Init(nil)
}

func handleMessage(msg *gst.Message) error {
	switch msg.Type() {
	case gst.MessageEOS:
		return app.ErrEOS
	case gst.MessageError:
		return msg.ParseError()
	}

	return nil
}

func LoopBus(lg *zap.Logger, pipeline *gst.Pipeline) {
	// Retrieve the bus from the pipeline
	bus := pipeline.GetPipelineBus()

	// Loop over messsages from the pipeline
	go func() {
		for {
			msg := bus.TimedPop(gst.ClockTimeNone)
			if msg == nil {
				return
			}
			if err := handleMessage(msg); err != nil {
				lg.Error("failed to handle message", zap.Error(err))
			}
		}
	}()
}
