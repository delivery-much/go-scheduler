package scheduler

import (
	"errors"
	"fmt"
	"time"
)

// start library base values
var (
	// db its the library designated database
	db JobDatabase = &emptyDB{}

	// logger its the library designated logger
	logger Logger = &emptyLogger{}

	// jobDefinitions maps the job names to its designated functions
	jobDefinitions map[string]*JobFunc = make(map[string]*JobFunc, 0)
)

// Init initis the scheduler library
func Init(c Config) (err error) {
	switch {
	case c.DB != nil:
		db = c.DB

	case c.MongoDB != nil && c.MongoDB.Conn != nil:
		if c.MongoDB.DbName == "" {
			c.MongoDB.DbName = "go-scheduler"
		}
		if c.MongoDB.CollName == "" {
			c.MongoDB.CollName = "scheduler-jobs"
		}
		db = newMongo(c.MongoDB.Conn, c.MongoDB.DbName, c.MongoDB.CollName)

	default:
		err = errors.New("No job DB or mongo conection was provided")
		return
	}

	err = db.InitJobDB()
	if err != nil {
		err = fmt.Errorf("Failed to init job db, %v", err)
		return
	}

	if c.Logger != nil {
		logger = c.Logger
	}

	if c.ProcessingRate.Seconds() == float64(0) {
		c.ProcessingRate = time.Minute
	}

	if c.Location != "" {
		l, err := time.LoadLocation(c.Location)
		if err != nil {
			return fmt.Errorf("Failed to load the provided location '%s', %v", c.Location, err)
		}

		location = l
	}

	deleteOnCancel = c.DeleteOnCancel
	deleteOnDone = c.DeleteOnDone

	go processJobs(c.ProcessingRate)
	return
}

// Define inserts a new job definition
// given the job name and the function to be called when that job is triggered.
//
// If the jobName was already previously defined, the previous function will be overridden
func Define(jobName string, fn JobFunc) {
	jobDefinitions[jobName] = &fn
}

// In creates a new definition of a SIMPLE and PENDING job to be run once in the provided duration.
//
// This function does not save the schedule in the database yet,
// the function Do must be called subsequently so that the job can be defined and saved.
func In(d time.Duration) *simpleScheduleDefinition {
	return &simpleScheduleDefinition{
		nextRunAt: now().Add(d),
	}
}

// On creates a new definition of a SIMPLE and PENDING job to be run once in the provided date time.
//
// This function does not save the schedule in the database yet,
// the function Do must be called subsequently so that the job can be validated and saved.
//
// IMPORTANT: Please note that, for the library flow to function correctly,
// the provided time value (t) should be in the same timezone as configured during the library instantiation.
func On(t time.Time) *simpleScheduleDefinition {
	return &simpleScheduleDefinition{
		nextRunAt: t,
	}
}

// Every schedules a RECURRENT job to run repeatedly, given the schedule string.
//
// The schedule string expects the following formats:
//
// - A time interval string (Ex.: "1 minute", "2 months", "6 years")
//
// - A time string in HH:MM format (Ex.: "11:27")
//
// - A weekday string (Ex.: "monday", "friday")
//
// - A weekday and time string (Ex.: "monday at 12:00", "friday at 15:08")
func Every(schedule string) *recurrentScheduleDefinition {
	return &recurrentScheduleDefinition{
		schedule: schedule,
	}
}

// List lists jobs on the database given the finder.
func List(f Finder) ([]*Job, error) {
	return db.List(f)
}
