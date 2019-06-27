package actiontracker_test

import (
	"sync"
	"testing"

	"github.com/DanMaples/ActionTracker/actiontracker"
)

func TestActionTrackerKeepsCorrectAverage(t *testing.T) {
	//TODO: check report is correct
	// tracker := actiontracker.New()
	// tracker.AddAction("jump", 100)
	// tracker.AddAction("jump", 200)
	// tracker.AddAction("jump", 300)
	// report := tracker.GetStats()
}

func TestActionTrackerConcurencey(t *testing.T) {
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
	//TODO: check report is correct
	t.Logf(tracker.GetStats())
}

func TestMaxActionsHaveBeenAddedReturnsError(t *testing.T) {
	actionName := "jump"
	tracker := actiontracker.NewMaxedCountActionTracker(actionName)
	err := tracker.AddAction(actionName, 1)
	if err != nil {
		t.Fatalf("recieved unexpected err: %+v", err)
	}
	err = tracker.AddAction(actionName, 1)
	if actiontracker.TooManyValuesError != err.Error() {
		t.Fatalf("expected error of '%s', but recieved '%s' instead", actiontracker.TooManyValuesError, err.Error())
	}
}
