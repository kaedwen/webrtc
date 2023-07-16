package sonos

import (
	"errors"
	"net/url"
	"time"

	so "github.com/szatmary/sonos"
	"go.uber.org/zap"
)

type SonosPlayer struct {
	lg *zap.Logger
	zp *so.ZonePlayer
}

func NewSonosPlayer(lg *zap.Logger, target string) (*SonosPlayer, error) {
	zp, err := SearchTarget(lg, target)
	if err != nil {
		return nil, err
	}

	return &SonosPlayer{lg, zp}, nil
}

func (p *SonosPlayer) Play(uri *url.URL) error {
	err := p.zp.SetAVTransportURI(uri.String())
	if err != nil {
		return err
	}

	return p.zp.Play()
}

func SearchTarget(lg *zap.Logger, target string) (*so.ZonePlayer, error) {
	son, err := so.NewSonos()
	if err != nil {
		return nil, err
	}
	defer son.Close()

	found, _ := son.Search()
	to := time.After(10 * time.Second)
	for {
		select {
		case <-to:
			return nil, errors.New("timeout")
		case zp := <-found:
			lg.Info("found player", zap.String("name", zp.RoomName()))
			if zp.RoomName() == target {
				return zp, nil
			}
		}
	}
}
