package common

import (
	"fmt"
	"io"
	"io/fs"
	"mime"
	"path/filepath"
	"strings"
)

type ContentEncoding int

const (
	PLAIN ContentEncoding = iota
	GZ
	BR
)

type StaticFileInfo struct {
	Size         int64
	Reader       io.ReadCloser
	ExtraHeaders map[string]string
	Mime         string
}

type StaticSource interface {
	Open(path string) (StaticSourceFile, error)
}

type StaticSourceFile interface {
	Reader() io.ReadCloser
	Stat() (fs.FileInfo, error)
}

type StaticHandler struct {
	source StaticSource
	index  string
}

func NewStaticHandler(s StaticSource, i string) *StaticHandler {
	return &StaticHandler{s, i}
}

func (h *StaticHandler) get(p string, encoding string) (*StaticFileInfo, error) {
	mt := mime.TypeByExtension(filepath.Ext(p))

	switch true {
	case strings.Contains(encoding, "br"):
		lt := p + ".br"
		if o, err := h.source.Open(lt); err == nil {
			if i, err := o.Stat(); err == nil && !i.IsDir() {
				return &StaticFileInfo{
					ExtraHeaders: map[string]string{
						"Content-Encoding": "br",
						"Vary":             "Accept-Encoding",
					},
					Size:   i.Size(),
					Mime:   mt,
					Reader: o.Reader(),
				}, nil
			}
		}
	case strings.Contains(encoding, "gzip"):
		lt := p + ".gz"
		if o, err := h.source.Open(lt); err == nil {
			if i, err := o.Stat(); err == nil && !i.IsDir() {
				return &StaticFileInfo{
					ExtraHeaders: map[string]string{
						"Content-Encoding": "gzip",
					},
					Size:   i.Size(),
					Mime:   mt,
					Reader: o.Reader(),
				}, nil
			}
		}
	}

	if o, err := h.source.Open(p); err == nil {
		if i, err := o.Stat(); err == nil && !i.IsDir() {
			return &StaticFileInfo{
				Size:   i.Size(),
				Mime:   mt,
				Reader: o.Reader(),
			}, nil
		}
	}

	return nil, fmt.Errorf("path %s not found", p)
}

func (h *StaticHandler) Get(p string, encoding string) (*StaticFileInfo, error) {

	if i, err := h.get(p, encoding); err == nil {
		return i, nil
	}

	if i, err := h.get(h.index, encoding); err == nil {
		return i, nil
	}

	return nil, fmt.Errorf("path %s and index not found", p)
}
