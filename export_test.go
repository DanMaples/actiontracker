package actiontracker

//TooManyValuesError is the error message for when
//AddAction() returns an error.
const TooManyValuesError = tooManyValuesError

//NewMaxedCountJSONActionTracker will return a JSON action tracker
//where the supplied key's count is maxUint -1.  Used for testing only.
func NewMaxedCountJSONActionTracker(actionKey string) ActionTracker {
	ati := &actionTrackerImpl{
		ActionFormatter: NewJSONFormatter(),
		actions:         make(map[string]*actionAverage),
	}
	ati.actions[actionKey] = &actionAverage{count: maxUint - 1}
	return ati
}
