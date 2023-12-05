package scheduler

type ScheduleType string

const (
	SIMPLE    = ScheduleType("SIMPLE")
	RECURRENT = ScheduleType("RECURRENT")
)

// String returns the schedule type in string notation
func (t ScheduleType) String() string {
	return string(t)
}
