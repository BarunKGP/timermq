package entities

import (
	"strings"
	"time"

	"github.com/BarunKGP/timermq/internal/values"
	"github.com/google/uuid"
)

type OptionalArgs struct {
	Delay    time.Duration
	Ttl      time.Duration
	Durable  bool
	Loggable bool
}

type Message struct {
	id        uuid.UUID
	rawstring string

	cmd  values.CommandMethod
	val  string
	args OptionalArgs
}

func (m Message) GetId() uuid.UUID {
	return m.id
}

func NewMessage(data string) *Message {
	return &Message{id: uuid.New(), rawstring: data}
}

func NewMessageFromTokens(tokens []string) *Message {
	data := strings.Join(tokens, string(' '))
	return NewMessage(data)
}

func (m *Message) WithPush() (*Message, error) {
	m.cmd = values.Push
	return m, nil
}

func (m *Message) WithPing() (*Message, error) {
	m.cmd = values.Ping
	return m, nil
}

func (m *Message) WithGet() (*Message, error) {
	m.cmd = values.Get
	return m, nil
}

func (m *Message) WithCancel() (*Message, error) {
	m.cmd = values.Cancel
	return m, nil
}

func (m *Message) SetValue(val string) {
	m.val = val
}

func (m *Message) SetArgs(args OptionalArgs) {
	m.args = args
}

func (m *Message) CommandType() values.CommandMethod {
	return m.cmd
}

func (m *Message) GetValue() string {
	return m.val
}

func (m *Message) GetValueBytes() []byte {
	return []byte(m.val)
}

func (m *Message) GetDelay() time.Duration {
	return m.args.Delay
}
