package entities

import (
	"fmt"
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

	// Args
	delay   *time.Duration
	ttl     *time.Duration
	durable *bool
}

func (m Message) GetId() uuid.UUID {
	return m.id
}

func (m Message) ToString() string {
	s := fmt.Sprintf("%s %s", m.cmd, m.val)

	if m.delay != nil {
		s = s + fmt.Sprintf(" delayMs=%d,", m.delay.Milliseconds())
	}
	if m.ttl != nil {
		s = s + fmt.Sprintf(" ttlMs=%d,", m.ttl.Milliseconds())
	}
	if m.durable != nil {
		s = s + fmt.Sprintf(" durable=%t,", *m.durable)
	}

	return strings.TrimRight(s, ", ")
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
