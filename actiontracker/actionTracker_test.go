package actiontracker_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/DanMaples/ActionTracker/actiontracker"
)

func TestActionTrackerKeepsCorrectAverage(t *testing.T) {
	tracker := actiontracker.New()
	tracker.AddAction(`{"action":"jump", "time":100}`)
	tracker.AddAction(`{"action":"run", "time":100}`)
	tracker.AddAction(`{"action":"jump", "time":200}`)
	tracker.AddAction(`{"action":"run", "time":0}`)
	tracker.AddAction(`{"action":"run", "time":50}`)
	tracker.AddAction(`{"action":"jump", "time":300}`)

	actualStats := tracker.GetStats()
	const expectedStats = `[{"action":"jump","avg":200},{"action":"run","avg":50}]`
	if expectedStats != actualStats {
		t.Fatalf("expected stats : %s did not match actual stats : %s", expectedStats, actualStats)
	}
}

func TestActionTrackerConcurencey(t *testing.T) {
	tracker := actiontracker.New()
	wg := sync.WaitGroup{}
	for count := 0; count < 10; count++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tracker.AddAction(`{"action":"jump", "time":100}`)
			tracker.GetStats()
			tracker.AddAction(`{"action":"run", "time":75}`)
			tracker.GetStats()
			tracker.AddAction(`{"action":"jump", "time":200}`)
		}()
	}
	wg.Wait()
	//TODO: check report is correct - need to figure out rounding first though...
	t.Logf(tracker.GetStats())
}

func TestInputValidation(t *testing.T) {
	var testCases = []struct {
		name        string
		input       string
		expectError bool
	}{
		{"properlyFormatedJSONReturnsNoError", `{"action":"jump", "time":100}`, false},
		{"improperlyFormatedJSONReturnsAnError", `"action":"jump", "time":100}`, true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tracker := actiontracker.New()
			err := tracker.AddAction(testCase.input)
			if testCase.expectError {
				if err == nil {
					t.Fatal("Expected an error to be returned, but none was.")
				}
			} else {
				if err != nil {
					t.Fatalf("No error was expected, but got error: %+v", err)
				}
			}
		})
	}
}

func TestGetStatsPanicsIfJSONIsUnmarshallable(t *testing.T) {
	oldMarshalJSON := *actiontracker.MarshalJSON
	defer func() { *actiontracker.MarshalJSON = oldMarshalJSON }()
	*actiontracker.MarshalJSON = func(interface{}) ([]byte, error) { return nil, fmt.Errorf("I had an error") }
	defer func() {
		recoveryValue := recover()
		if recoveryValue == nil {
			t.Fatal("GetStats did not convert error to panic")
		}
	}()
	tracker := actiontracker.New()
	tracker.GetStats()
}

func TestMaxActionsHaveBeenAddedReturnsError(t *testing.T) {
	actionName := "jump"
	tracker := actiontracker.NewMaxedCountActionTracker(actionName)
	err := tracker.AddAction(`{"action":"jump", "time":100}`)
	if err != nil {
		t.Fatalf("recieved unexpected err: %+v", err)
	}
	err = tracker.AddAction(`{"action":"jump", "time":100}`)
	if actiontracker.TooManyValuesError != err.Error() {
		t.Fatalf("expected error of '%s', but recieved '%s' instead",
			actiontracker.TooManyValuesError,
			err.Error())
	}
}
