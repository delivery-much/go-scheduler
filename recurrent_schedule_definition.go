package scheduler

import (
	"fmt"
	"time"
)

type recurrentScheduleDefinition struct {
	schedule  string
	limitDate *time.Time
}

// Do effectivelly schedules the job on the database to run in the configured time, given the job name.
//
// You can also provide extra data that will be saved with the job.
func (rsd *recurrentScheduleDefinition) Do(jobName string, data ...map[string]any) (err error) {
	jobFunc := jobDefinitions[jobName]
	if jobFunc == nil {
		err = fmt.Errorf("No job definition with the name %s was found", jobName)
		return
	}

	t, err := getNextScheduleDate(rsd.schedule)
	if err != nil {
		return
	}

	d := make(map[string]any)
	if len(data) > 0 {
		d = data[0]
	}

	job := Job{
		Status:            PENDING,
		ScheduleType:      RECURRENT,
		NextRunAt:         t,
		ScheduleString:    rsd.schedule,
		ScheduleLimitDate: rsd.limitDate,
		Name:              jobName,
		Data:              d,
	}

	return db.SaveJob(job)
}

// Until sets a limit date for the RECURRENT job to run.
//
// When the limit date arrives, the job is set as DONE.
//
// IMPORTANT: Please note that, so that the library flow works properly,
// the provided time value (t) should be in the same timezone as configured in the library instantiation.
func (rsd *recurrentScheduleDefinition) Until(t time.Time) *recurrentScheduleDefinition {
	rsd.limitDate = &t
	return rsd
}
