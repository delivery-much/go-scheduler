package scheduler

import "time"

var (
	location = time.UTC
)

// now returns the current time in the given location that is configured in the library
func now() time.Time {
	return time.Now().In(location)
}
