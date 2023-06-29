package common

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *ConfigLogging) (*zap.Logger, error) {
	var c zap.Config
	if cfg.Level == "debug" {
		c = zap.NewDevelopmentConfig()
	} else {
		c = zap.NewProductionConfig()
	}

	c.EncoderConfig.TimeKey = "time"
	c.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	lg, err := c.Build()
	if err != nil {
		return nil, err
	}

	return lg, nil
}
