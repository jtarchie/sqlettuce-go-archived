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
	listener net.Listener
	poolSize uint
	port     uint
}

func NewServer(
	port uint,
	poolSize uint,
) (*Server, error) {
	listener, err := reuseport.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return nil, fmt.Errorf("could not listen for tcp: %w", err)
	}

	return &Server{
		port:     port,
		poolSize: poolSize,
		listener: listener,
	}, nil
}

func (s *Server) Listen(handler Handler) error {
	var totalConnections atomic.Uint64

	workerPool := worker.New(
		int(s.poolSize),
		int(s.poolSize),
		func(worker int, conn net.Conn) {
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

	slog.Info("started server",
		slog.Uint64("port", uint64(s.port)),
		slog.Uint64("pool", uint64(s.poolSize)),
	)

	for {
		conn, err := s.listener.Accept()

		//nolint:errorlint
		if opErr, ok := err.(*net.OpError); ok && !opErr.Temporary() {
			return nil
		}

		if err != nil {
			return fmt.Errorf("could not accept connection for tcp: %w", err)
		}

		workerPool.Enqueue(conn)
	}
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("could not close server: %w", err)
	}

	return nil
}
