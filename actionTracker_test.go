package actiontracker_test

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"github.com/DanMaples/actiontracker"
)

func TestActionTrackerKeepsCorrectAverage(t *testing.T) {
	tracker := actiontracker.NewWithJSONActionFormatter()
	if err := tracker.AddAction(`{"action":"jump", "time":100}`); err != nil {
		t.Fatalf("Received unexpected error from AddAction: %+v", err)
	}
	if err := tracker.AddAction(`{"action":"run", "time":100}`); err != nil {
		t.Fatalf("Received unexpected error from AddAction: %+v", err)
	}
	if err := tracker.AddAction(`{"action":"jump", "time":200}`); err != nil {
		t.Fatalf("Received unexpected error from AddAction: %+v", err)
	}
	if err := tracker.AddAction(`{"action":"run", "time":0}`); err != nil {
		t.Fatalf("Received unexpected error from AddAction: %+v", err)
	}
	if err := tracker.AddAction(`{"action":"run", "time":50}`); err != nil {
		t.Fatalf("Received unexpected error from AddAction: %+v", err)
	}
	if err := tracker.AddAction(`{"action":"jump", "time":300}`); err != nil {
		t.Fatalf("Received unexpected error from AddAction: %+v", err)
	}
	if err := tracker.AddAction(`{"action":"swim", "time":3.69}`); err != nil {
		t.Fatalf("Received unexpected error from AddAction: %+v", err)
	}
	if err := tracker.AddAction(`{"action":"swim", "time":2}`); err != nil {
		t.Fatalf("Received unexpected error from AddAction: %+v", err)
	}

	actualStats := tracker.GetStats()
	const expectedStats = `[{"action":"jump","avg":200},{"action":"run","avg":50},{"action":"swim","avg":2.845}]`
	if expectedStats != actualStats {
		t.Fatalf("expected stats : %s did not match actual stats : %s", expectedStats, actualStats)
	}
}

func TestActionTrackerConcurencey(t *testing.T) {
	tracker := actiontracker.NewWithJSONActionFormatter()
	wg := sync.WaitGroup{}
	for count := 0; count < 1000; count++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := tracker.AddAction(`{"action":"jump", "time":100}`); err != nil {
				t.Fatalf("Received unexpected error from AddAction: %+v", err)
			}
			tracker.GetStats()
			if err := tracker.AddAction(`{"action":"run", "time":75}`); err != nil {
				t.Fatalf("Received unexpected error from AddAction: %+v", err)
			}
			tracker.GetStats()
			if err := tracker.AddAction(`{"action":"jump", "time":200}`); err != nil {
				t.Fatalf("Received unexpected error from AddAction: %+v", err)
			}
		}()
	}
	wg.Wait()
	actualStats := tracker.GetStats()
	const expectedStats = `[{"action":"jump","avg":150},{"action":"run","avg":75}]`
	if expectedStats != actualStats {
		t.Fatalf("expected stats : %s did not match actual stats : %s", expectedStats, actualStats)
	}
}

func TestAddActionJSONInputValidation(t *testing.T) {
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
			tracker := actiontracker.NewWithJSONActionFormatter()
			err := tracker.AddAction(testCase.input)
			if testCase.expectError && err == nil {
				t.Fatal("Expected an error to be returned, but none was.")
			}
			if !testCase.expectError && err != nil {
				t.Fatalf("No error was expected, but got error: %+v", err)
			}
		})
	}
}

type marshalErrorFormatter struct{}

func (mef marshalErrorFormatter) InputFormatter(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (mef marshalErrorFormatter) OutputFormatter(v interface{}) ([]byte, error) {
	return nil, fmt.Errorf("I had an error")
}

func TestGetStatsPanicsIfActionsAreUnmarshallable(t *testing.T) {
	defer func() {
		recoveryValue := recover()
		if recoveryValue == nil {
			t.Fatal("GetStats did not convert error to panic")
		}
	}()
	tracker := actiontracker.NewWithCustomActionFormatter(marshalErrorFormatter{})
	tracker.GetStats()
}

func TestMaxActionsHaveBeenAddedReturnsError(t *testing.T) {
	actionName := "jump"
	tracker := actiontracker.NewMaxedCountJSONActionTracker(actionName)
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
