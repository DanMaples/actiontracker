package actiontracker

//ActionFormatter defines an interface to format action inputs and outputs
type ActionFormatter interface {
	InputFormatter(data []byte, v interface{}) error
	OutputFormatter(v interface{}) ([]byte, error)
}

//StucturedActionInput is a struct used to parse raw action input into.
//Exported so custom ActionFormatters can use it.
type StucturedActionInput struct {
	Action string  `json:"action"`
	Time   float64 `json:"time"`
}

//StructuredStatsOutput is a struct used to parse action output data into.
//Exported so custom ActionFormatters can use it.
type StructuredStatsOutput struct {
	Action string  `json:"action"`
	Avg    float64 `json:"avg"`
}
