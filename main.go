package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello from timermq - a timer-based message queue, written in Golang")

	ports := []string{"tcp", "http", "amqp"}
	_ = ports[1]
}
