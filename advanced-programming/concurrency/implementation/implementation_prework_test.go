package main

import (
	"fmt"
	"testing"

	"golang.org/x/sync/errgroup"
)

type testCase struct {
	name    string
	service func() (idService, func())
}

func setup() []testCase {
	goroutineService := MakeGoroutineIdService()
	goroutineService.Start()

	return []testCase{
		// {"no-sync", func() (idService, func()) {
		// 	service := &noSyncIdService{}
		// 	teardown := func() {}
		// 	return service, teardown
		// }},
		{"atomic", func() (idService, func()) {
			service := &atomicIdService{}
			teardown := func() {}
			return service, teardown
		}},
		{"mutex", func() (idService, func()) {
			service := &mutexIdService{}
			teardown := func() {}
			return service, teardown
		}},
		{"goroutines", func() (idService, func()) {
			service := MakeGoroutineIdService()
			service.Start()
			teardown := func() { service.Stop() }
			return service, teardown
		}},
	}
}

func RunService(t testing.TB, service idService, numWorkers, numCalls int) {
	t.Helper()

	var eg errgroup.Group
	idChan := make(chan uint64, numWorkers*numCalls)

	for i := 0; i < numWorkers; i++ {
		eg.Go(func() error {
			lastId := uint64(0)
			for j := 0; j < numCalls; j++ {
				id := service.getNext()
				if id < lastId {
					return fmt.Errorf("Ids not monotonically increasing: got %d after %d", id, lastId)
				}
				idChan <- id
			}
			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		t.Fatalf(err.Error())
	}

	close(idChan)

	expectedMax := numWorkers * numCalls
	maxId := uint64(0)
	for id := range idChan {
		if maxId < id {
			maxId = id
		}
	}
	if maxId != uint64(expectedMax) {
		t.Fatalf("Max id across workers incorrect: expected %d, got %d", expectedMax, maxId)
	}
}

func TestServices(t *testing.T) {
	cases := setup()
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			service, teardown := testCase.service()
			defer teardown()
			RunService(t, service, 10, 10000)
		})
	}
}

func BenchmarkServices(b *testing.B) {
	cases := setup()
	for _, testCase := range cases {
		b.Run(testCase.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				service, teardown := testCase.service()
				defer teardown()
				RunService(b, service, 10, 10000)
			}
		})
	}
}
