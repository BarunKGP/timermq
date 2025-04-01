package servers

import (
	"fmt"

	"github.com/BarunKGP/timermq/internal/adapters"
)

type Server interface {
	Start()

	Close() error
}

func NewServer(key adapters.ConnType, opts adapters.InitOpts) (Server, error) {
	switch key {
	case adapters.TCP:
		return NewTCPServer(opts), nil

	default:
		return nil, fmt.Errorf("Invalid key %+v", key)
	}
}
