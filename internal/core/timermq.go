package core

import (
	"fmt"
	"sync"
	"time"
)

type MessageIndex = int

type TimerMQ struct {
	store    *Store[[]byte]
	dlq      map[MessageIndex][]byte
	capacity int

	mCh    chan []byte
	timers map[MessageIndex]*time.Timer
	mu     sync.Mutex
}

func NewTimerMQ(cap int) *TimerMQ {
	return &TimerMQ{
		store:    NewStore[[]byte](),
		dlq:      map[MessageIndex][]byte{},
		capacity: cap,

		mCh:    make(chan []byte, cap),
		timers: map[int]*time.Timer{},
		mu:     sync.Mutex{},
	}
}

func Len(tmq *TimerMQ) int {
	return tmq.store.Len()
}

func (tmq *TimerMQ) IsArchived(index MessageIndex) bool {
	tmq.mu.Lock()
	defer tmq.mu.Unlock()
	_, exists := tmq.dlq[index]
	return exists
}

func (tmq *TimerMQ) Ping() string {
	return "pong"
}

func (tmq *TimerMQ) Archive(index MessageIndex) error {
	tmq.mu.Lock()
	defer tmq.mu.Unlock()
	if _, exists := tmq.dlq[index]; exists {
		return nil
	}

	data, err := tmq.store.Get(StoreIndex(index))
	if err != nil {
		return err
	}

	tmq.dlq[index] = data
	return nil
}

func (tmq *TimerMQ) AddTimer(index MessageIndex, timer *time.Timer) {
	tmq.mu.Lock()
	defer tmq.mu.Unlock()
	if _, exists := tmq.timers[index]; exists {
		return
	}
	tmq.timers[index] = timer
}

func (tmq *TimerMQ) Publish(data []byte, delay time.Duration) MessageIndex {
	newIndex := tmq.store.Consume(data)

	timer := time.AfterFunc(delay, func() {
		if tmq.IsArchived(newIndex) {
			return
		}
		tmq.mu.Lock()
		delete(tmq.timers, newIndex)
		tmq.mu.Unlock()
		tmq.mCh <- data
		close(tmq.mCh)
	})

	tmq.AddTimer(newIndex, timer)
	return newIndex
}

func (tmq *TimerMQ) Listen() [][]byte {
	res := [][]byte{}
	select {
	case data := <-tmq.mCh:
		res = append(res, data)
	}
	return res
}

func (tmq *TimerMQ) CancelSend(index MessageIndex) error {
	tmq.mu.Lock()
	if timer, exists := tmq.timers[index]; exists {
		timer.Stop()
		tmq.Archive(index)
		delete(tmq.timers, index)
		return nil
	}
	return fmt.Errorf("index %+v does not exist", index)
}
