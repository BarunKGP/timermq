/**
`listener` listens for messages over a protocol.
It validates incoming messages and calls the required functions
to parse messages into TimerMQ commands
*/

package listeners

import "github.com/BarunKGP/timermq/internal/adapters"

type Listener struct {
	protocol adapters.Protocol
}
