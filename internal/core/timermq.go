package core

import (
	"fmt"
	"log/slog"
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

func (t *TimerMQ) Close() {
	slog.Info("Closing TimerMQ")
	close(t.mCh)
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
	return tmq.archive(index)
}

func (tmq *TimerMQ) archive(index MessageIndex) error {
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

func (tmq *TimerMQ) NumActiveTimers() int {
	tmq.mu.Lock()
	defer tmq.mu.Unlock()

	return len(tmq.timers)
}

func (tmq *TimerMQ) ActiveTimerKeys() []MessageIndex {
	keys := []int{}
	for k := range tmq.timers {
		keys = append(keys, k)
	}
	return keys
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
		slog.Debug("Pushing to mCh", "data", data)
		tmq.mCh <- data
		// close(tmq.mCh)
	})

	tmq.AddTimer(newIndex, timer)
	return newIndex
}

func (tmq *TimerMQ) Listen() [][]byte {
	res := [][]byte{}
	for data := range tmq.mCh {
		res = append(res, data)
		slog.Debug("Reading from mCh", "data", data)
	}
	slog.Debug("Created output", "output", res)
	return res
}

func (tmq *TimerMQ) CancelSend(index MessageIndex) error {
	tmq.mu.Lock()
	defer tmq.mu.Unlock()

	timer, exists := tmq.timers[index]
	if !exists {
		return fmt.Errorf("index %+v does not exist", index)
	}

	slog.Debug("Found timer", "index", index, "timer", timer)
	tmq.archive(index)
	delete(tmq.timers, index)
	timer.Stop()
	return nil
}
