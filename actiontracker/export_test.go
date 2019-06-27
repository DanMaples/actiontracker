package actiontracker

const TooManyValuesError = tooManyValuesError

var MarshalJSON = &marshalJSON

func NewMaxedCountActionTracker(actionKey string) ActionTracker {
	ati := &actionTrackerImpl{actions: make(map[string]*actionAverage)}
	ati.actions[actionKey] = &actionAverage{count: maxUint - 1}
	return ati
}
