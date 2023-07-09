package common

import (
	"fmt"
	"net"

	"github.com/alexflint/go-arg"
)

type Config struct {
	VideoOut ConfigVideoOutputStream `arg:"group:VideoOut"`
	AudioOut ConfigAudioOutputStream `arg:"group:AudioOut"`
	AudioIn  ConfigAudioInputStream  `arg:"group:AudioIn"`
	Logging  ConfigLogging           `arg:"group:Logging"`
	Http     ConfigHTTP              `arg:"group:Http"`
}

type ConfigStream struct {
	VideoOut ConfigVideoOutputStream
	AudioOut ConfigAudioOutputStream
	AudioIn  ConfigAudioInputStream
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

type ConfigVideoOutputStream struct {
	Source string `arg:"--video-src,env:VIDEO_SRC" default:"v4l2src"`
	Device string `arg:"--video-device,env:VIDEO_DEVICE" default:"/dev/video0"`
	Codec  string `arg:"--video-codec,env:VIDEO_CODEC" default:"vp8"`
	Height uint   `arg:"--video-height,env:VIDEO_HEIGHT" default:"480"`
	Width  uint   `arg:"--video-width,env:VIDEO_WIDTH" default:"640"`
}

type ConfigAudioOutputStream struct {
	Source     string  `arg:"--audio-src,env:AUDIO_SRC" default:"alsasrc"`
	DeviceName string  `arg:"--audio-device-name,env:AUDIO_DEVICE" default:"default"`
	Device     *string `arg:"--audio-device,env:AUDIO_DEVICE"`
	Codec      string  `arg:"--audio-codec,env:AUDIO_CODEC" default:"opus"`
	Channels   uint    `arg:"--audio-channels,env:AUDIO_CHANNELS" default:"1"`
}

type ConfigAudioInputStream struct {
	Sink       string  `arg:"--audio-src,env:AUDIO_SINK" default:"alsasink"`
	DeviceName string  `arg:"--audio-device-name,env:AUDIO_DEVICE" default:"default"`
	Device     *string `arg:"--audio-device,env:AUDIO_DEVICE"`
	Codec      string  `arg:"--audio-codec,env:AUDIO_CODEC" default:"opus"`
	Channels   uint    `arg:"--audio-channels,env:AUDIO_CHANNELS" default:"1"`
}

func (c *Config) Stream() *ConfigStream {
	return &ConfigStream{
		VideoOut: c.VideoOut,
		AudioOut: c.AudioOut,
		AudioIn:  c.AudioIn,
	}
}

func (c *ConfigHTTP) Address() string {
	return net.JoinHostPort(c.Host, fmt.Sprint(c.Port))
}

func (c *Config) MustParse() {
	arg.MustParse(c)
}
