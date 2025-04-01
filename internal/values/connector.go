package values

import "fmt"

type Connector struct {
	Port      uint16
	Addr      string
	KeepAlive bool
	Closed    bool
	Protocol  Protocol
}

func (s *Connector) GetFullAddress() string {
	return fmt.Sprintf("%s:%d", s.Addr, s.Port)
}

type ConnType int

const (
	TCP ConnType = iota
	HTTP
	AMQP
)

type InitOpts struct {
	Addr      string   `json:"addr"`
	Port      uint16   `json:"port"`
	Protocol  ConnType `json:"protocol"`
	KeepAlive bool     `json:"persistent,omitempty"`
	Capacity  int      `json:"capacity"`
}
