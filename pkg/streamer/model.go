package streamer

import (
	"github.com/go-gst/go-gst/gst"
	"github.com/kaedwen/webrtc/pkg/common"
)

type StreamElement struct {
	Kind       string
	Codec      common.StreamCodec
	Properties map[string]interface{}
	Caps       *CapsBuilder
	Queue      bool
}

type CapsBuilder struct {
	caps Caps
}

type Caps struct {
	mime   string
	values []CapsValue
}

type CapsValue struct {
	name  string
	value any
}

func NewCapsBuilder(mime string) *CapsBuilder {
	return &CapsBuilder{
		caps: Caps{mime, nil},
	}
}

func (b *CapsBuilder) Height(v int) *CapsBuilder {
	b.caps.values = append(b.caps.values, CapsValue{
		name:  "height",
		value: v,
	})
	return b
}

func (b *CapsBuilder) Width(v int) *CapsBuilder {
	b.caps.values = append(b.caps.values, CapsValue{
		name:  "width",
		value: v,
	})
	return b
}

func (b *CapsBuilder) Channels(v int) *CapsBuilder {
	b.caps.values = append(b.caps.values, CapsValue{
		name:  "channels",
		value: v,
	})
	return b
}

func (b *CapsBuilder) Format(v string) *CapsBuilder {
	b.caps.values = append(b.caps.values, CapsValue{
		name:  "format",
		value: v,
	})
	return b
}

func (b *CapsBuilder) Rate(v int) *CapsBuilder {
	b.caps.values = append(b.caps.values, CapsValue{
		name:  "rate",
		value: v,
	})
	return b
}

func (b *CapsBuilder) Build() *gst.Caps {
	caps := gst.NewEmptySimpleCaps(b.caps.mime)
	for _, v := range b.caps.values {
		caps.SetValue(v.name, v.value)
	}

	return caps
}
