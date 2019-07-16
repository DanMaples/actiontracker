package actiontracker

import "testing"

func TestActionGetRoundedAverage(t *testing.T) {
	testAction := newAction()
	testAction.add(333.3333333333)

	var testCases = []struct {
		digits   uint
		expected float64
	}{
		{1, 333.3},
		{5, 333.33333},
		{10, 333.3333333333},
	}
	for _, testCase := range testCases {
		actual := testAction.getRoundedAvg(testCase.digits)
		if testCase.expected != actual {
			t.Errorf("Expected rounded average of %v does not match actual rounded average of %v", testCase.expected, actual)
		}
	}
}
