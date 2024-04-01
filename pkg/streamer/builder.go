package streamer

import (
	"fmt"
	"strings"
)

type PipelineBuilder struct {
	parts []string
}

func NewPipelineBuilder() *PipelineBuilder {
	return &PipelineBuilder{}
}

func (pb *PipelineBuilder) Add(element string) *PipelineBuilder {
	pb.parts = append(pb.parts, element)
	return pb
}

func (pb *PipelineBuilder) AddWithProperties(element string, properties map[string]any) *PipelineBuilder {
	pl := make([]string, 0, len(properties))
	for k, v := range properties {
		pl = append(pl, fmt.Sprint(k, "=", v))
	}

	pb.parts = append(pb.parts, fmt.Sprint(element, " ", strings.Join(pl, " ")))
	return pb
}

func (pb *PipelineBuilder) AddFilter(c *Caps) {
	if c != nil {
		pb.parts = append(pb.parts, c.Build())
	}
}

func (pb *PipelineBuilder) Build() string {
	return strings.Join(pb.parts, " ! ")
}
