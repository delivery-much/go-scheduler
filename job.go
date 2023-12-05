package scheduler

import "time"

var (
	deleteOnDone   = false
	deleteOnCancel = false
)

// Job represents a schedule job.
type Job struct {
	// ID its the job ID in the database
	ID string

	// ScheduleType represents the job schedule type, if its a job that runs only once (SIMPLE) or if its a job that runs recurrently (RECURRENT)
	ScheduleType

	// Status represents the job schedule current status, if its pending, failed, done, or canceled
	Status ScheduleStatus

	// NextRunAt defines when the job should run
	NextRunAt time.Time

	// LastRunAt defines when the job was last ran
	LastRunAt *time.Time

	// ScheduleString its the schedule string that was defined when the job schedule was configured,
	// when its a RECURRENT schedule
	ScheduleString string

	// ScheduleLimitDate defines the limit date that the RECURRENT job will run.
	//
	// When the limit date arrives, the job is set as DONE.
	// If no limit date is set, the job will run forever until its manually canceled on deleted.
	ScheduleLimitDate *time.Time

	// Name represents the job definition name
	Name string

	// Data represents the extra data that the user can provide when defining a job
	Data map[string]any
}

// Done sets the job schedule status as DONE and saves it on the database
func (j *Job) Done() error {
	j.Status = DONE

	if deleteOnDone {
		return db.DeleteJob(*j)
	}

	return db.SaveJob(*j)
}

// Fail sets the job schedule status as FAILED and saves it on the database
func (j *Job) Fail() error {
	j.Status = FAILED

	return db.SaveJob(*j)
}

// Cancel sets the job schedule status as CANCELED and saves it on the database
func (j *Job) Cancel() error {
	j.Status = CANCELED

	if deleteOnCancel {
		return db.DeleteJob(*j)
	}

	return db.SaveJob(*j)
}

// Delete deletes the job from the database
func (j *Job) Delete() error {
	return db.DeleteJob(*j)
}

// IsDone returns true if the job status is DONE, and false otherwise
func (j *Job) IsDone() bool {
	return j.Status == DONE
}

// HasFailed returns true if the job status is FAILED, and false otherwise
func (j *Job) HasFailed() bool {
	return j.Status == FAILED
}

// IsCanceled returns true if the job status is CANCELED, and false otherwise
func (j *Job) IsCanceled() bool {
	return j.Status == CANCELED
}

// IsPending returns true if the job status is PENDING, and false otherwise
func (j *Job) IsPending() bool {
	return j.Status == PENDING
}

// IsSimple return true if the job schedule type is SIMPLE, and false otherwise
func (j *Job) IsSimple() bool {
	return j.ScheduleType == SIMPLE
}

// IsRecurrent return true if the job schedule type is RECURRENT, and false otherwise
func (j *Job) IsRecurrent() bool {
	return j.ScheduleType == RECURRENT
}
