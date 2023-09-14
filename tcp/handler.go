package tcp

import (
	"io"
)

type Handler interface {
	OnConnection(io.ReadWriter) error
}
