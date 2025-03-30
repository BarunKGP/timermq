package values

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type CommandMethod string
type MessageIndex uuid.UUID

const (
	Push   CommandMethod = "PUSH"
	Get                  = "GET"
	Delay                = "DELAY"
	Cancel               = "CANCEL"
	Ping                 = "PING"
)

var MinimumRequiredArgs = map[CommandMethod]int{
	Push:   1,
	Get:    1,
	Cancel: 1,
	Delay:  2,
	Ping:   0,
}
var (
	ErrMsgTooShort = errors.New("Invalid message: message missing essential parameters")
)

var commandMap = map[string]CommandMethod{
	"PUSH":   Push,
	"GET":    Get,
	"CANCEL": Cancel,
}

func CmdFromString(s string) (CommandMethod, error) {
	c, ok := commandMap[s]
	if !ok {
		return CommandMethod(""), fmt.Errorf("Unrecognized command %s", s)
	}
	return c, nil
}

func ParseValidateCommand(words []string) (CommandMethod, error) {
	if len(words) < 1 {
		return CommandMethod(""), ErrMsgTooShort
	}

	cmd, err := CmdFromString(words[0])
	if err != nil {
		return cmd, err
	}

	if len(words) < MinimumRequiredArgs[cmd]+1 {
		return CommandMethod(""), ErrMsgTooShort
	}

	return cmd, nil
}

type ContentType string

const (
	PlainText ContentType = "text/plain"
	JsonText  ContentType = "text/json"
)

type CommandHeader struct {
	Type          ContentType `json:"Content-Type"`
	ClientId      string      `json:"Client-Id"`
	ContentLength int         `json:"Content-Length"`
}

// func parseCmd(cmdString string) (Response, error) {
// 	words := strings.SplitN(cmdString, " ", 2)
// 	if len(words) == 0 {
// 		return Response{}, newMessageParsingErr(1, "invalid command")
// 	}
// 	cmd := words[0]
// 	switch commandMap[cmd] {
// 	case Push:
// 		a := strings.Split(words[1], ", ")
// 		if len(a) == 0 {
// 			return Response{}, newMessageParsingErr(1, "empty item pushed")
// 		}
// 		if len(a) == 1 {
// 			return Response{Cmd: Push, Body: a[0]}, nil
// 		}
// 		if len(a) == 2 {
// 			return Response{}, newMessageParsingErr(1, "unrecognized argument %s for Push")
// 		}
// 		return Response{}, newMessageParsingErr(1, fmt.Sprintf("invalid command %s", cmd))
//
// 	case Get:
// 		return Response{Cmd: Get}, nil
// 	default:
// 		return Response{}, newMessageParsingErr(1, fmt.Sprintf("invalid command %s", cmd))
// 	}
// }
