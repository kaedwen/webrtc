package streamer

import (
	"fmt"
	"strings"

	"github.com/kaedwen/webrtc/pkg/common"
)

type StreamElement struct {
	Kind       string
	Codec      common.StreamCodec
	Bitrate    uint
	Properties map[string]interface{}
	SrcCaps    *Caps
	EnvCaps    *Caps
	Queue      bool
}

type Caps struct {
	mime   string
	filter map[string]any
}

func NewCaps(mime string, filter map[string]any) *Caps {
	return &Caps{mime, filter}
}

func (c *Caps) Build() string {
	fv := make([]string, 0, len(c.filter))
	for k, v := range c.filter {
		fv = append(fv, fmt.Sprint(k, "=", v))
	}

	return fmt.Sprint(c.mime, ",", strings.Join(fv, ","))
}
