# actiontracker
Library to track actions

This library will keep track of the average time for actions.

The library provides an ActionTracker interface to use.

The ActionTracker interface has two methods:

    AddAction(input string) error

    GetStats() string

Two constructors are provided:

    NewWithCustomActionFormatter(af ActionFormatter)

    NewWithJSONActionFormatter()

The ActionFormatter interface has two methods:

    InputFormatter(data []byte, v interface{}) error

    OutputFormatter(v interface{}) ([]byte, error)


# Example
The following example:

```
import (
	"fmt"
	"github.com/DanMaples/actiontracker"
)

func foo() {
    tracker := actiontracker.NewWithJSONActionFormatter()
    if err := tracker.AddAction(`{"action":"jump", "time":100}`); err != nil {
        //handle err
    }
    if err := tracker.AddAction(`{"action":"run", "time":75}`); err != nil {
            //handle err
    }

    if err := tracker.AddAction(`{"action":"jump", "time":200}`); err != nil {
            //handle err
    }
    fmt.Println(tracker.GetStats())
}
```
would produce the following output:
```
[{"action":"jump", "avg":150}, {"action":"run", "avg":75}]
```

# Notes
- Output is rounded to the nearest 3 decimal places.
- This library can make use of custom input/output parsers if you decide to use something other than JSON.
- This library is safe for concurrent use.
- To run the unit tests for this library, run "go test -v -race" from the root directory.
- Output from GetStats() will be sorted alphabetically by action name.

# Developer Considerations
- When authoring this library, a decision had to be made betwen:
(1) Exact precision for averages while only accepting integer numbers, but the possibility that data overflow could happen very quickly depending on input.
(2) High precision for averages of rational numbers, while allowing for a much greater number of action inputs before returning an error.
- Choice number two was taken so the library could accept floating point numbers.  The code was written such that it would be trivial to convert to choice number one if required.