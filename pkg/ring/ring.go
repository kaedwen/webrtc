package ring

import (
	"context"
	"fmt"
	"net/url"

	"github.com/holoplot/go-evdev"
	"github.com/kaedwen/webrtc/pkg/common"
	"github.com/kaedwen/webrtc/pkg/ring/sonos"
	"go.uber.org/zap"
)

type Player interface {
	Play(uri *url.URL) error
}

type RingHandler struct {
	lg      *zap.Logger
	cfg     *common.ConfigRing
	players []Player
}

func NewRingHandler(lg *zap.Logger, cfg *common.ConfigRing) (*RingHandler, error) {
	spl, err := sonos.NewSonosPlayers(lg.With(zap.String("context", "sonos")), cfg.SonosTarget)
	if err != nil {
		return nil, err
	}

	pl := make([]Player, 0, len(spl))
	for _, sp := range spl {
		pl = append(pl, sp)
	}

	return &RingHandler{lg, cfg, pl}, nil
}

func (h *RingHandler) Watch(ctx context.Context) error {
	d, err := evdev.Open(h.cfg.Device)
	if err != nil {
		return err
	}

	vMajor, vMinor, vMicro := d.DriverVersion()
	h.lg.Info("input driver running", zap.String("version", fmt.Sprintf("%d.%d.%d", vMajor, vMinor, vMicro)))

	key := evdev.KEYFromString[h.cfg.Key]

	bu := url.URL{
		Scheme: "http",
		Host:   "192.168.60.141:9099",
	}

	go func() {
		err := d.NonBlock()
		if err != nil {
			h.lg.Fatal("failed to set non block", zap.Error(err))
		}

		for {
			e, err := d.ReadOne()
			if err == nil && e.Code == key && e.Value == 0 {
				tu := bu.JoinPath(h.cfg.JingleName)
				for _, p := range h.players {
					if err := p.Play(tu); err != nil {
						h.lg.Error("failed to play", zap.Error(err))
					}
				}
			}
		}
	}()

	go func() {
		<-ctx.Done()
		_ = d.Close()
	}()

	return nil
}
