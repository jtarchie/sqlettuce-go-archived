package tcp

import (
	"context"
	"io"
)

type Handler interface {
	OnConnection(context.Context, io.ReadWriter) error
}
