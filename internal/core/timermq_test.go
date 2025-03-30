package core

import (
	"bytes"
	"runtime"
	"testing"
	"time"
)

func TestLen(t *testing.T) {
	t.Log("TimerMQ: Testing message queue length")
	tmq := NewTimerMQ(5)
	len := Len(tmq)
	if len != 0 {
		t.Errorf("Expected length of empty TimerMQ to be 0, received %d", len)
	}

	msg := []byte("test message")
	tmq.Publish(msg, 0)

	len = Len(tmq)
	if len != 1 {
		t.Errorf("Expected length of empty TimerMQ to be 1, received %d", len)
	}
}

func contains(store [][]byte, item []byte) bool {
	for _, b := range store {
		if bytes.Equal(b, item) {
			return true
		}
	}
	return false
}

func TestPublish(t *testing.T) {
	t.Log("TimerMQ: Testing publish and listen")
	tmq := NewTimerMQ(2)
	msg1, msg2 := []byte("test message 1"), []byte("test message 2")

	// var rcv [][]byte
	go func() {
		id1 := tmq.Publish(msg1, time.Second)
		t.Logf("Published message with id: %d", id1)

		id2 := tmq.Publish(msg2, 0)
		t.Logf("Published message with id: %d", id2)

	}()

	// Wait for all messages to be sent
	time.Sleep(time.Second + 500*time.Millisecond)
	tmq.Close()

	rcv := tmq.Listen()
	// t.Logf("Received: %X", rcv)
	if len(rcv) != 2 {
		t.Errorf("Unexpected number of items returned: %d -> %s", len(rcv), rcv)
	}

	if exists := contains(rcv, msg1); !exists {
		t.Errorf("msg: %s missing in received slice", msg1)
	}

	if exists := contains(rcv, msg2); !exists {
		t.Errorf("msg: %s missing in received slice", msg2)
	}

	t.Logf("Currently running %d goroutines", runtime.NumGoroutine())
}
