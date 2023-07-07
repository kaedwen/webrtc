package server

import "io"

type StaticHandler interface {
	Get(path string) io.ReadCloser
}
