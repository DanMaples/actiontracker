package actiontracker

import "encoding/json"

//JSONFormatter is a struct to be used as an ActionFormatter
type JSONFormatter struct{}

//NewJSONFormatter returns an ActionFormatter that parses JSON
func NewJSONFormatter() ActionFormatter {
	return &JSONFormatter{}
}

//InputFormatter parses JSON into structs
func (jf JSONFormatter) InputFormatter(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

//OutputFormatter parses structs into JSON
func (jf JSONFormatter) OutputFormatter(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}
