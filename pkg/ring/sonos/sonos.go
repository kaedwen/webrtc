package sonos

import (
	"net/url"
	"time"

	so "github.com/szatmary/sonos"
	"go.uber.org/zap"
)

type SonosPlayer struct {
	lg *zap.Logger
	zp *so.ZonePlayer
}

func NewSonosPlayers(lg *zap.Logger, target string) ([]*SonosPlayer, error) {
	spl, err := SearchTargets(lg, target)
	if err != nil {
		return nil, err
	}

	return spl, nil
}

func (p *SonosPlayer) Play(uri *url.URL) error {
	err := p.zp.SetAVTransportURI(uri.String())
	if err != nil {
		return err
	}

	return p.zp.Play()
}

func SearchTargets(lg *zap.Logger, target string) ([]*SonosPlayer, error) {
	players := make([]*SonosPlayer, 0)

	son, err := so.NewSonos()
	if err != nil {
		return nil, err
	}
	defer son.Close()

	found, _ := son.Search()
	to := time.After(5 * time.Second)
	for {
		select {
		case <-to:
			return players, nil
		case zp := <-found:
			lg.Info("found player", zap.String("name", zp.RoomName()))
			if target == "-" || zp.RoomName() == target {
				players = append(players, &SonosPlayer{lg, zp})
			}
		}
	}
}
