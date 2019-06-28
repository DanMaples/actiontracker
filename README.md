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


# To Use
In your imports section of your go file that will be using this library, include:

```
import "github.com/DanMaples/actiontracker"

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
    statsString := tracker.GetStats()
}
```