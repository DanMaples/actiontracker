package actiontracker

import (
	"errors"
	"math"
)

const maxUint = ^uint(0)
const tooManyValuesError = "can't continue to track action, too many values have been added to track"

//action interface is an interface used to manipulate actions
type action interface {
	add(float64) error
	getRoundedAvg(uint) float64
}

//actionImpl is a stuct used to keep track of an action
//it is a concrete impl of the action interface
type actionImpl struct {
	value float64
	count uint
}

//newAction will return a zero-value action
func newAction() action {
	return &actionImpl{}
}

//add will add a time to an action
func (a *actionImpl) add(time float64) error {
	if a.count == maxUint {
		return errors.New(tooManyValuesError)
	}
	a.count++
	a.value = a.value + (time-a.value)/float64(a.count)
	return nil
}

//getRoundedAvg will return the average time from an action
func (a *actionImpl) getRoundedAvg(decimalPlaces uint) float64 {
	var roundingFactor uint = 1
	var i uint
	for ; i < decimalPlaces; i++ {
		roundingFactor *= 10
	}
	return math.Round(a.value*float64(roundingFactor)) / float64(roundingFactor)
}
