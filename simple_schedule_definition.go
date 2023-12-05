package scheduler

import (
	"fmt"
	"time"
)

type simpleScheduleDefinition struct {
	nextRunAt time.Time
}

// Do effectivelly schedules the job on the database to run in the configured time, given the job name.
//
// You can also provide extra data that will be saved with the job.
func (ssd *simpleScheduleDefinition) Do(jobName string, data ...map[string]any) (err error) {
	jobFunc := jobDefinitions[jobName]
	if jobFunc == nil {
		err = fmt.Errorf("No job definition with the name %s was found", jobName)
		return
	}

	d := make(map[string]any)
	if len(data) > 0 {
		d = data[0]
	}

	job := Job{
		Status:       PENDING,
		ScheduleType: SIMPLE,
		NextRunAt:    ssd.nextRunAt,
		Name:         jobName,
		Data:         d,
	}

	return db.SaveJob(job)
}
