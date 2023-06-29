package common

import (
	"fmt"
	"net"

	"github.com/alexflint/go-arg"
)

type Config struct {
	Logging ConfigLogging `arg:"-"`
	Http    ConfigHTTP    `arg:"-"`
}

type ConfigLogging struct {
	Level string `arg:"--log-level,env:LOG_LEVEL" default:"debug"`
}

type ConfigHTTP struct {
	Host             string  `arg:"--http-host,env:HTTP_HOST" default:"0.0.0.0"`
	Port             uint    `arg:"--http-port,env:HTTP_PORT" default:"8080"`
	PathGetLiveness  string  `arg:"env:HTTP_PATH_LIVENESS" default:"/healthz"`
	PathGetReadiness string  `arg:"env:HTTP_PATH_READINESS" default:"/readyz"`
	StaticPath       *string `arg:"--http-static,env:HTTP_STATIC"`
}

func (c *ConfigHTTP) Address() string {
	return net.JoinHostPort(c.Host, fmt.Sprint(c.Port))
}

func (c *Config) MustParse() {
	arg.MustParse(&c.Logging)
	arg.MustParse(&c.Http)
	arg.MustParse(c)
}
