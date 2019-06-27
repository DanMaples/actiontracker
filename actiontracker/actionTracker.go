package actiontracker

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"sync"
)

const maxUint = ^uint(0)
const tooManyValuesError = "can't continue to track action, too many values have been added to track"

var marshalJSON = json.Marshal

//ActionTracker is the interface for an actionTracker
type ActionTracker interface {
	AddAction(string) error
	GetStats() string
}

//addActionInput is a struct used to parse raw json into
type addActionInput struct {
	Action string  `json:"action"`
	Time   float64 `json:"time"`
}

//getStatsOutput is a struct used to parse output data into json
type getStatsOutput struct {
	Action string  `json:"action"`
	Avg    float64 `json:"avg"`
}

//actionAverage is a stuct used to keep a running average of an action
type actionAverage struct {
	value float64
	count uint
}

//actionTrackerImpl is the implementation of the interface
type actionTrackerImpl struct {
	sync.RWMutex
	actions map[string]*actionAverage
}

//AddAction will add an action
func (ati *actionTrackerImpl) AddAction(rawInput string) error {
	var parsedInput addActionInput
	if err := json.Unmarshal([]byte(rawInput), &parsedInput); err != nil {
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

//GetStats will get the stats
func (ati *actionTrackerImpl) GetStats() string {
	output := make([]*getStatsOutput, 0)
	ati.RLock()
	sortedActions := ati.getSortedActions()
	for _, action := range sortedActions {
		output = append(output, &getStatsOutput{Action: action, Avg: ati.actions[action].value})
	}
	ati.RUnlock()
	statsBytes, err := marshalJSON(output)
	if err != nil {
		panic(fmt.Sprintf("programming error detected: %+v", err))
	}
	return string(statsBytes)
}

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

//New will create a new ActionTracker
func New() ActionTracker {
	return &actionTrackerImpl{actions: make(map[string]*actionAverage)}
}

/*TODO:
read over problem statement again for details
fill out readme
figure out rounding
rename a bunch of stuff
improve comments
*/
