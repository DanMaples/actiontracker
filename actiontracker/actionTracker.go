package actiontracker

import (
	"errors"
	"fmt"
	"sync"
)

const maxUint = ^uint(0)

//ActionTracker is the interface for an actionTracker
type ActionTracker interface {
	AddAction(key string, value float64) error
	GetStats() string
}

type actionAverage struct {
	value float64
	Count uint //Count is exported soley for mocking purposes. Acceptable because actionAverage is not exported.
}

//actionTrackerImpl is the implementation of the interface
type actionTrackerImpl struct {
	sync.RWMutex
	Actions map[string]*actionAverage //Actions is exported soley for mocking purposes. Acceptable because actionTrackerImpl is not exported.
}

//AddAction will add an action
func (ati *actionTrackerImpl) AddAction(key string, value float64) error {
	ati.Lock()
	defer ati.Unlock()

	if _, exists := ati.Actions[key]; !exists {
		ati.Actions[key] = &actionAverage{}
	} else if ati.Actions[key].Count == maxUint {
		return errors.New("can't continue to track action, too many values have been added to track")
	}

	ati.Actions[key].Count++
	ati.Actions[key].value = ati.Actions[key].value + (value-ati.Actions[key].value)/float64(ati.Actions[key].Count)
	return nil
}

//GetStats will get the stats
func (ati *actionTrackerImpl) GetStats() string {
	retString := ""
	ati.RLock()
	defer ati.RUnlock()

	for actionAverageName, theactionAverage := range ati.Actions {
		retString += fmt.Sprintf("key:%s value:%v\n", actionAverageName, theactionAverage.value)
	}
	return retString
}

//New will create a new ActionTracker
func New() ActionTracker {
	return &actionTrackerImpl{Actions: make(map[string]*actionAverage)}
}
