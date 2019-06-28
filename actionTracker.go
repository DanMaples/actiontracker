package actiontracker

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"
)

const maxUint = ^uint(0)
const tooManyValuesError = "can't continue to track action, too many values have been added to track"

//ActionTracker is the interface for an actionTracker
type ActionTracker interface {
	AddAction(string) error
	GetStats() string
}

//NewWithCustomActionFormatter will create and return
//a new ActionTracker with the supplied ActionFormatter
func NewWithCustomActionFormatter(af ActionFormatter) ActionTracker {
	return &actionTrackerImpl{
		ActionFormatter: af,
		actions:         make(map[string]*actionAverage),
	}
}

//NewWithJSONActionFormatter will create and return
//a new ActionTracker that takes JSON input
func NewWithJSONActionFormatter() ActionTracker {
	return NewWithCustomActionFormatter(NewJSONFormatter())
}

//actionAverage is a stuct used to keep track of the average of an action
type actionAverage struct {
	value float64
	count uint
}

//actionTrackerImpl is the concrete implementation of the ActionTracker interface
type actionTrackerImpl struct {
	sync.RWMutex
	actions map[string]*actionAverage
	ActionFormatter
}

//AddAction will parse the input and add the action to the tracker
func (ati *actionTrackerImpl) AddAction(rawInput string) error {
	var parsedInput StucturedActionInput
	if err := ati.InputFormatter([]byte(rawInput), &parsedInput); err != nil {
		return err
	}
	ati.Lock()
	defer ati.Unlock()
	if _, exists := ati.actions[parsedInput.Action]; !exists {
		ati.actions[parsedInput.Action] = &actionAverage{}
	} else if ati.actions[parsedInput.Action].count == maxUint {
		return errors.New(tooManyValuesError)
	}
	ati.actions[parsedInput.Action].count++
	ati.actions[parsedInput.Action].value = ati.actions[parsedInput.Action].value + (parsedInput.Time-ati.actions[parsedInput.Action].value)/float64(ati.actions[parsedInput.Action].count)
	return nil
}

//GetStats will return the stats about the actions from the tracker
//Output will be rounded to the nearest 3 decimal places
func (ati *actionTrackerImpl) GetStats() string {
	output := make([]*StructuredStatsOutput, 0)
	ati.RLock()
	sortedActions := ati.getSortedActions()
	for _, action := range sortedActions {
		output = append(output, &StructuredStatsOutput{
			Action: action,
			Avg:    math.Round(ati.actions[action].value*1000) / 1000, //round to the nearest 3 decimal places
		})
	}
	ati.RUnlock()
	statsBytes, err := ati.OutputFormatter(output, "", "    ")
	if err != nil {
		panic(fmt.Sprintf("programming error detected: %+v", err))
	}
	return string(statsBytes)
}

//getSortedActions will return a slice of actions, sorted alphabetically
func (ati *actionTrackerImpl) getSortedActions() []string {
	actions := make([]string, len(ati.actions))
	index := 0
	for action := range ati.actions {
		actions[index] = action
		index++
	}
	sort.Strings(actions)
	return actions
}
