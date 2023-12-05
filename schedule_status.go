package scheduler

type ScheduleStatus string

const (
	DONE     = ScheduleStatus("DONE")
	FAILED   = ScheduleStatus("FAILED")
	PENDING  = ScheduleStatus("PENDING")
	CANCELED = ScheduleStatus("CANCELED")
)

// String returns the schedule status in string notation
func (t ScheduleStatus) String() string {
	return string(t)
}
