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
	count uint
}

//actionTrackerImpl is the implementation of the interface
type actionTrackerImpl struct {
	sync.RWMutex
	actions map[string]*actionAverage
}

//AddAction will add an action
func (ati *actionTrackerImpl) AddAction(key string, value float64) error {
	ati.Lock()
	defer ati.Unlock()

	if _, exists := ati.actions[key]; !exists {
		ati.actions[key] = &actionAverage{}
	} else if ati.actions[key].count == maxUint {
		return errors.New("can't continue to track action, too many values have been added to track")
	}

	ati.actions[key].count++
	ati.actions[key].value = ati.actions[key].value + (value-ati.actions[key].value)/float64(ati.actions[key].count)
	return nil
}

//GetStats will get the stats
func (ati *actionTrackerImpl) GetStats() string {
	retString := ""
	ati.RLock()
	defer ati.RUnlock()

	for actionAverageName, theActionAverage := range ati.actions {
		retString += fmt.Sprintf("key:%s value:%v\n", actionAverageName, theActionAverage.value)
	}
	return retString
}

//New will create a new ActionTracker
func New() ActionTracker {
	return &actionTrackerImpl{actions: make(map[string]*actionAverage)}
}
