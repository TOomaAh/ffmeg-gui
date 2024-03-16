package core

import (
	"context"
	"time"

	"github.com/TOomaAh/media_tools/internal/component"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type Filter struct {
	path       string
	conditions *[]*component.ConditionalWidget
}

func NewFilter(path string, conditions *[]*component.ConditionalWidget) *Filter {
	return &Filter{path: path, conditions: conditions}
}

func (f *Filter) Apply() bool {
	for _, c := range *f.conditions {
		data := f.GetFfprobeInfo()
		if !c.Choice.Check(data) {
			return false
		}
	}
	return true
}

func (f *Filter) GetFfprobeInfo() *ffprobe.ProbeData {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	data, err := ffprobe.ProbeURL(ctx, f.path)

	if err != nil {
		return &ffprobe.ProbeData{}
	}

	return data
}
