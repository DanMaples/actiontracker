package actiontracker

import (
	"fmt"
	"sort"
	"sync"
)

const defaultDecimalPlaces = 3

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
		actions:         make(map[string]action),
	}
}

//NewWithJSONActionFormatter will create and return
//a new ActionTracker that takes JSON input
func NewWithJSONActionFormatter() ActionTracker {
	return NewWithCustomActionFormatter(NewJSONFormatter())
}

//actionTrackerImpl is the concrete implementation of the ActionTracker interface
type actionTrackerImpl struct {
	sync.RWMutex
	actions map[string]action
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
		ati.actions[parsedInput.Action] = newAction()
	}
	return ati.actions[parsedInput.Action].add(parsedInput.Time)
}

//GetStats will return the stats about the actions from the tracker
func (ati *actionTrackerImpl) GetStats() string {
	ati.RLock()
	output := make([]*StructuredStatsOutput, len(ati.actions))
	sortedActions := ati.getSortedActions()
	for index, action := range sortedActions {
		output[index] = &StructuredStatsOutput{
			Action: action,
			Avg:    ati.actions[action].getRoundedAvg(defaultDecimalPlaces),
		}
	}
	ati.RUnlock()
	statsBytes, err := ati.OutputFormatter(output)
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
