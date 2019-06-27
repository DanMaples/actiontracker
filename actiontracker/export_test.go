package actiontracker

const MaxUint = maxUint

type ActionTrackerImpl = actionTrackerImpl
type ActionAverage = actionAverage

func NewMaxedCountActionTracker(actionKey string) ActionTracker {
	ati := &actionTrackerImpl{actions: make(map[string]*actionAverage)}
	ati.actions[actionKey] = &actionAverage{count: maxUint - 1}
	return ati
}
