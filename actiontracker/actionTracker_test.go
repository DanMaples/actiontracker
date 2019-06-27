package actiontracker_test

import (
	"sync"
	"testing"

	"github.com/DanMaples/ActionTracker/actiontracker"
)

func TestActionTracker(t *testing.T) {

	tracker := actiontracker.New()
	wg := sync.WaitGroup{}

	for count := 0; count < 10; count++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tracker.AddAction("jump", 100)
			tracker.GetStats()
			tracker.AddAction("run", 75)
			tracker.GetStats()
			tracker.AddAction("jump", 200)
		}()
	}
	wg.Wait()

	t.Logf(tracker.GetStats())
}
