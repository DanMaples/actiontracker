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
This library uses float64 data types to keep track of the averages.

Output is rounded to the nearest 3 decimal places.