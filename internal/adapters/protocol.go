package adapters

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/BarunKGP/timermq/internal/entities"
	"github.com/BarunKGP/timermq/internal/values"
)

var (
	ErrMsgParse            = errors.New("Unable to parse message")
	ErrMsgTooShort         = errors.New("Invalid message: message missing essential parameters")
	ErrInvalidCommand      = errors.New("Invalid command")
	ErrInvalidCommandArgs  = errors.New("Invalid args for command")
	ErrUnsupportedConnType = errors.New("Unsupported connection type")
)

type Protocol struct {
	ConnectionType ConnType
	Delim          byte
}

func TCPProtocol() Protocol {
	return Protocol{ConnectionType: TCP, Delim: byte('\n')}
}

func (p *Protocol) DecodeMessage(msg string) (*entities.Message, error) {
	words := strings.Split(msg, string(' '))
	cmd, err := values.ParseValidateCommand(words)
	if err != nil {
		return &entities.Message{}, err
	}

	switch p.ConnectionType {
	case TCP:
		switch cmd {
		case "PING":
			return handlePing(words)
		case "PUSH":
			return handlePush(words)
		default:
			return &entities.Message{}, ErrInvalidCommand
		}
	default:
		return &entities.Message{}, ErrUnsupportedConnType
	}
}

func (p *Protocol) EncodeMessage(msg *entities.Message) (string, error) {
	switch p.ConnectionType {
	case TCP:
		return msg.ToString(), nil
	default:
		return "", ErrUnsupportedConnType

	}
}

func handlePush(tokens []string) (*entities.Message,
	error) {
	msg, err := entities.NewMessageFromTokens(tokens).WithPush()
	if err != nil {
		return &entities.Message{}, ErrInvalidCommand
	}
	msg.SetValue(tokens[1])

	if len(tokens) == 2 {
		return msg, nil
	}

	args := entities.OptionalArgs{}
	for _, tok := range tokens[2:] {
		parts := strings.Split(tok, string("="))
		if len(parts) != 2 {
			return &entities.Message{}, ErrInvalidCommandArgs
		}

		switch strings.TrimSpace(parts[0]) {
		case "delay":
			delayMs, err := strconv.Atoi(parts[1])
			if err != nil {
				return &entities.Message{}, ErrInvalidCommandArgs
			}
			args.Delay = time.Duration(delayMs) * time.Millisecond
		case "durable":
			durable, err := strconv.ParseBool(parts[1])
			if err != nil {
				return &entities.Message{}, ErrInvalidCommandArgs
			}
			args.Durable = durable
		default:
			return &entities.Message{}, ErrInvalidCommandArgs
		}
	}

	msg.SetArgs(args)
	return msg, nil
}

func handlePing(tokens []string) (*entities.Message, error) {
	if len(tokens) > 1 {
		return &entities.Message{}, ErrInvalidCommand
	}
	msg, err := entities.NewMessageFromTokens(tokens).WithPing()
	if err != nil {
		return &entities.Message{}, ErrInvalidCommand
	}

	return msg, nil
}
