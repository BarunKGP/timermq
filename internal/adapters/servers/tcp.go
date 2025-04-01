package servers

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"

	"log/slog"

	"github.com/BarunKGP/timermq/internal/adapters"
	"github.com/BarunKGP/timermq/internal/core"
	"github.com/BarunKGP/timermq/internal/values"
	"github.com/google/uuid"
)

type TCPServer struct {
	tmq   *core.TimerMQ
	idMap map[uuid.UUID]core.MessageIndex
	adapters.Connector
}

func NewTCPServer(opts adapters.InitOpts) *TCPServer {
	return &TCPServer{
		tmq:   core.NewTimerMQ(opts.Capacity),
		idMap: make(map[uuid.UUID]core.MessageIndex),
		Connector: adapters.Connector{
			Port:      opts.Port,
			Addr:      opts.Addr,
			KeepAlive: opts.KeepAlive,
			Protocol:  adapters.TCPProtocol(),
		},
	}
}

func (t *TCPServer) Close() error {
	if t.Closed {
		return fmt.Errorf("Server is already closed!")
	}
	t.Closed = true
	t.tmq.Close()
	return nil
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	slog.Info("New connection created")
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		// buf := make([]byte, 2048)
		// n, err := conn.Read(buf)
		str, err := reader.ReadString(s.Protocol.Delim)
		if err != nil {
			if errors.Is(err, io.EOF) {
				slog.Info("Connection closed")
				continue
			}
			slog.Error("Connection error", "error", err)
			continue
		}
		slog.Debug("Received data from client", "data", str)

		msg, err := s.Protocol.DecodeMessage(str)
		if err != nil {
			slog.Error("Unable to parse message", "error", err, "message", str)
		}

		switch msg.CommandType() {
		case values.Push:
			id := s.tmq.Publish(msg.GetValueBytes(), msg.GetDelay())
			s.idMap[msg.GetId()] = id

			slog.Info("Published message", "messageId", msg.GetId(), "timermqId", id, "delayMs", msg.GetDelay().Milliseconds())
		case values.Ping:
			res := s.tmq.Ping()
			if res == "pong" {
				slog.Info("Ping returned: " + res)
			} else {
				slog.Info("Ping failed")
				slog.Warn("TimerMQ ping failed", "res", res)
			}
		default:
			slog.Error("Unrecognized command type", "cmd", msg.CommandType())
		}
	}
}

func (s *TCPServer) Start() {
	slog.Info("Starting server", "address", s.GetFullAddress())
	listener, err := net.Listen("tcp", s.GetFullAddress())
	if err != nil {
		slog.Error("Failed to start server", "error", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, io.EOF) {
				slog.Info("Connection closed")
				continue
			}
			slog.Error("Connection error", "error", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

var _ Server = &TCPServer{}
