package clients

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/BarunKGP/timermq/internal/adapters"
	"github.com/BarunKGP/timermq/internal/core"
	"github.com/BarunKGP/timermq/internal/entities"
	"github.com/google/uuid"
)

type TCPClient struct {
	tmq     *core.TimerMQ
	idMap   map[uuid.UUID]core.MessageIndex
	closeCh chan bool
	adapters.Connector
	pConn net.Conn
}

func NewTCPClient(opts adapters.InitOpts) *TCPClient {
	return &TCPClient{
		tmq:     core.NewTimerMQ(opts.Capacity),
		idMap:   make(map[uuid.UUID]core.MessageIndex),
		closeCh: make(chan bool, 1),
		pConn:   nil,
		Connector: adapters.Connector{
			Port:      opts.Port,
			Addr:      opts.Addr,
			KeepAlive: opts.KeepAlive,
			Protocol:  adapters.TCPProtocol(),
		},
	}
}

func (s *TCPClient) Close() error {
	if s.Closed {
		return fmt.Errorf("Client already closed")
	}

	slog.Info("Closing TCPClient")
	s.tmq.Close()
	s.Closed = true
	s.closeCh <- true
	close(s.closeCh)
	if s.KeepAlive {
		s.pConn.Close()
	}

	return nil
}

func (s *TCPClient) Push(val string, args entities.OptionalArgs) (*entities.Message, error) {
	rawstr := fmt.Sprintf("PUSH %s delayMs=%d", val, args.Delay.Milliseconds())
	msg, err := entities.NewMessage(rawstr).WithPush()
	if err != nil {
		return &entities.Message{}, err
	}

	msg.SetValue(val)
	msg.SetArgs(args)
	return msg, nil
}

func (s *TCPClient) Start() {
	slog.Info("Starting server", "address", s.GetFullAddress())
	conn, err := net.Dial("tcp", s.GetFullAddress())
	if err != nil {
		slog.Error("Failed to connect to server", "error", err)
		return
	}

	if s.KeepAlive {
		s.pConn = conn
	}

	defer conn.Close()

	for <-s.closeCh {
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
