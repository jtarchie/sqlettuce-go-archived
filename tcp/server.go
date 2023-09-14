package tcp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/jtarchie/worker"
	"github.com/libp2p/go-reuseport"
	"go.uber.org/atomic"
)

type Server struct {
	port uint64
}

func NewServer(
	port uint64,
) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Listen(handler Handler) error {
	var totalConnections atomic.Uint64

	workerPool := worker.New(100, 100, func(worker int, conn net.Conn) {
		currentConnection := totalConnections.Add(1)

		slog.Info("accepted new connection",
			slog.Int("worker", worker),
			slog.Uint64("connection", currentConnection),
		)

		err := handler.OnConnection(conn)
		if err != nil {
			slog.Error("connection errored",
				slog.Int("worker", worker),
				slog.Uint64("connection", currentConnection),
				slog.String("error", err.Error()),
			)
		}

		err = conn.Close()
		if err != nil {
			slog.Error("connection closed",
				slog.Int("worker", worker),
				slog.Uint64("connection", currentConnection),
				slog.String("error", err.Error()),
			)
		}

		slog.Info("connection closed",
			slog.Int("worker", worker),
			slog.Uint64("connection", currentConnection),
		)
	})

	listener, err := reuseport.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", s.port))
	if err != nil {
		return fmt.Errorf("could not listen for tcp: %w", err)
	}

	slog.Info("started server", slog.Uint64("port", s.port))

	defer func() {
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("could not accept connection for tcp: %w", err)
		}

		workerPool.Enqueue(conn)
	}
}
