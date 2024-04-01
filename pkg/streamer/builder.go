package streamer

import "fmt"

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
		pl = append(pl, fmt.Sprintf("%s=%v", k, v))
	}

	pb.parts = append(pb.parts, fmt.Sprint(element, pl))
	return pb
}
