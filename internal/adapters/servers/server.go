package servers

import (
	"fmt"
)

type Server interface {
	Start()

	Close() error
}

type ServerType int

const (
	TCP ServerType = iota
	HTTP
	AMQP
)

type InitOpts struct {
	Addr      string     `json:"addr"`
	Port      uint16     `json:"port"`
	Protocol  ServerType `json:"protocol"`
	KeepAlive bool       `json:"persistent,omitempty"`
	Capacity  int        `json:"capacity"`
}

func NewServer(key ServerType, opts InitOpts) (Server, error) {
	switch key {
	case TCP:
		return NewTCPServer(opts), nil

	default:
		return nil, fmt.Errorf("Invalid key %+v", key)
	}
}
