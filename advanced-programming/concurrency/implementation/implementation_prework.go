package main

import (
	"sync"
	"sync/atomic"
)

type idService interface {
	// Returns values in ascending order; it should be safe to call
	// getNext() concurrently without any additional synchronization.
	getNext() uint64
}
type noSyncIdService struct {
	id uint64
}

func (i *noSyncIdService) getNext() uint64 {
	i.id++
	return i.id
}

type atomicIdService struct {
	id uint64
}

func (i *atomicIdService) getNext() uint64 {
	return atomic.AddUint64(&i.id, 1)
}

type mutexIdService struct {
	sync.Mutex
	id uint64
}

func (i *mutexIdService) getNext() uint64 {
	i.Lock()
	defer i.Unlock()
	i.id += 1
	return i.id
}

type goroutineIdService struct {
	requests  chan struct{}
	responses chan uint64
}

func MakeGoroutineIdService() *goroutineIdService {
	service := goroutineIdService{
		requests:  make(chan struct{}),
		responses: make(chan uint64),
	}
	service.Start()
	return &service
}

func (s *goroutineIdService) Start() {
	go func() {
		id := uint64(0)
		for range s.requests {
			id++
			s.responses <- id
		}
	}()
}

func (s *goroutineIdService) Stop() {
	close(s.requests)
}

func (s *goroutineIdService) getNext() uint64 {
	s.requests <- struct{}{}
	return <-s.responses
}
