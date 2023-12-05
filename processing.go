package scheduler

import (
	"time"
)

func processJobs(rate time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("A panic occurred while processing jobs: %v", r)
		}

		processJobs(rate)
	}()

	for {
		time.Sleep(rate)

		process()
	}
}

func process() {
	jobs, err := db.ListExpiredSchedules()
	if err != nil {
		logger.Errorf("Failed to list expired schedules, %v", err)
		return
	}

	for _, j := range jobs {
		jobFunc := jobDefinitions[j.Name]
		if jobFunc == nil {
			logger.Errorf("Job %s was scheduled but no job definition with the name %s was found", j.ID, j.Name)
			failJob(j)
			continue
		}

		err := (*jobFunc)(j)
		if err != nil {
			logger.Errorf("Job %s failed, %v", j.ID, err)
			failJob(j)
			continue
		}

		now := now()
		j.LastRunAt = &now
		if j.IsSimple() ||
			(j.ScheduleLimitDate != nil && j.ScheduleLimitDate.Before(now)) {
			err = j.Done()
			if err != nil {
				logger.Errorf("Failed to save job %s after it was done processing, %v", j.ID, err)
			}
			continue
		}

		if j.ScheduleString == "" {
			logger.Errorf("Tried to re-schedule recurrent job %s, but it had no ScheduleString", j.ID)
			failJob(j)
			continue
		}

		nra, err := getNextScheduleDate(j.ScheduleString)
		if err != nil {
			logger.Errorf("Failed to get next schedule date for job %s, %v", j.ID, err)
			failJob(j)
			continue
		}

		j.Status = PENDING
		j.NextRunAt = nra
		err = db.SaveJob(*j)
		if err != nil {
			logger.Errorf("Failed to save job %s on the database to be re-scheduled, %v", j.ID, err)
			failJob(j)
			continue
		}
	}
}

func failJob(j *Job) {
	err := j.Fail()
	if err != nil {
		logger.Errorf("Failed to save job %s after it failed, %v", j.ID, err)
	}
}
