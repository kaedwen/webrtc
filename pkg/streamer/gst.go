package streamer

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
	"github.com/kaedwen/webrtc/pkg/common"
	"go.uber.org/zap"
)

const capsTagName = "caps"

func init() {
	gst.Init(nil)
}

type StreamElementCaps struct {
	Mime     string `caps:"mime"`
	Height   uint   `caps:"height"`
	Width    uint   `caps:"width"`
	Format   string `caps:"format"`
	Channels uint   `caps:"channels"`
	Rate     uint   `caps:"rate"`
}

type StreamElement struct {
	Kind       string
	Codec      common.StreamCodec
	Properties map[string]interface{}
	Caps       *StreamElementCaps
	Queue      bool
}

func (se StreamElement) ToGstCaps() *gst.Caps {
	if se.Caps == nil {
		return nil
	}

	t := reflect.ValueOf(*se.Caps)

	var mime string
	var caps []string
	for i := 0; i < t.NumField(); i++ {
		valueField := t.Field(i)
		tagField := t.Type().Field(i).Tag.Get(capsTagName)

		if tagField == "mime" {
			mime = valueField.String()
		} else {
			if !valueField.IsZero() {
				caps = append(caps, fmt.Sprintf("%s=%v", tagField, valueField.Interface()))
			}
		}
	}

	return gst.NewCapsFromString(strings.Join(append([]string{mime}, caps...), ","))
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
