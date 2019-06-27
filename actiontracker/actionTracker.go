package actiontracker

import (
	"fmt"
	"sync"
)

//ActionTracker is the interface for an actionTracker
type ActionTracker interface {
	AddAction(key string, value float64) error
	GetStats() string
}

type item struct {
	value float64
	count uint
}

//actionTrackerImpl is the implementation of the interface
type actionTrackerImpl struct {
	sync.RWMutex
	dataStore map[string]item
}

//AddAction will add an action
func (ati *actionTrackerImpl) AddAction(key string, value float64) error {
	ati.Lock()
	defer ati.Unlock()

	currentItem := ati.dataStore[key]
	currentItem.value = ((float64(currentItem.count) * currentItem.value) + value) / float64(currentItem.count+1)
	currentItem.count++
	ati.dataStore[key] = currentItem
	return nil
}

//GetStats will get the stats
func (ati *actionTrackerImpl) GetStats() string {
	retString := ""
	ati.RLock()
	defer ati.RUnlock()

	for itemName, theItem := range ati.dataStore {
		retString += fmt.Sprintf("key:%s value:%v\n", itemName, theItem.value)
	}
	return retString
}

//New will create a new ActionTracker
func New() ActionTracker {
	return &actionTrackerImpl{dataStore: make(map[string]item)}
}
