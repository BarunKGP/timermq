/**
`listener` listens for messages over a protocol.
It validates incoming messages and calls the required functions
to parse messages into TimerMQ commands
*/

package clients

import (
	"log/slog"
	"net"

	"github.com/BarunKGP/timermq/internal/adapters"
	"github.com/BarunKGP/timermq/internal/core"
	"github.com/google/uuid"
)

type Client interface {
	Start()
	Close() error
}
