package scheduler

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockDependencies mocks the database and logger library dependencies and returns a pointer to the mocks
func mockDependencies() (
	dbMock *databaseMock,
	logMock *loggerMock,
) {
	dbMock = newDatabaseMock()
	logMock = newLoggerMock()

	db = dbMock
	logger = logMock

	return
}

func TestProcessJobs(t *testing.T) {
	t.Run("When the database fails", func(t *testing.T) {
		t.Run("Should log an error if the database fails to list jobs", func(t *testing.T) {
			dbMock, loggerMock := mockDependencies()

			mockErr := errors.New("mock!!")
			dbMock.SetMethodResponse("ListExpiredSchedules", []*Job{}, mockErr)

			process()

			assert.True(t, dbMock.CalledOnce())
			assert.True(t, loggerMock.Method("Errorf").CalledWith("Failed to list expired schedules, %v"))
		})
		t.Run("Should do nothing if there are no expired jobs", func(t *testing.T) {
			dbMock, loggerMock := mockDependencies()

			dbMock.SetMethodResponse("ListExpiredSchedules", []*Job{}, nil)

			process()

			assert.True(t, dbMock.CalledOnce())
			assert.False(t, loggerMock.Called())
		})
	})
	t.Run("When the job fails", func(t *testing.T) {
		t.Run("Should log an error and fail the job if the job has no definition", func(t *testing.T) {
			dbMock, loggerMock := mockDependencies()

			mockJob := Job{
				Name: "a name that was not defined",
			}

			dbMock.SetMethodResponse("ListExpiredSchedules", []*Job{&mockJob}, nil)

			process()

			assert.True(t, dbMock.Method("ListExpiredSchedules").CalledOnce())
			assert.True(t, dbMock.Method("SaveJob").CalledOnce())
			assert.True(t, dbMock.Method("SaveJob").CalledWith(mockJob))
			assert.True(t, mockJob.HasFailed())
			assert.True(t, loggerMock.Method("Errorf").CalledWith("Job %s was scheduled but no job definition with the name %s was found"))
		})
		t.Run("Should log an error and fail the job if the job function returns an error", func(t *testing.T) {
			dbMock, loggerMock := mockDependencies()

			mockJobName := "MYMOCKJOB!"
			mockErr := errors.New("MOCK ERROR")
			mockJobFunc := func(j *Job) error {
				return mockErr
			}

			mockJob := Job{
				Name: mockJobName,
			}

			Define(mockJobName, mockJobFunc)

			dbMock.SetMethodResponse("ListExpiredSchedules", []*Job{&mockJob}, nil)

			process()

			assert.True(t, dbMock.Method("ListExpiredSchedules").CalledOnce())
			assert.True(t, dbMock.Method("SaveJob").CalledOnce())
			assert.True(t, dbMock.Method("SaveJob").CalledWith(mockJob))
			assert.True(t, mockJob.HasFailed())
			assert.True(t, loggerMock.Method("Errorf").CalledWith("Job %s failed, %v"))
		})
	})
	t.Run("When the job succeeds", func(t *testing.T) {
		t.Run("Should set the job as done if the job schedule is SIMPLE", func(t *testing.T) {
			dbMock, loggerMock := mockDependencies()

			mockJobName := "MYMOCKJOB!"
			mockJobFunc := func(j *Job) error {
				return nil
			}

			mockJob := Job{
				Name:         mockJobName,
				ScheduleType: SIMPLE,
			}

			Define(mockJobName, mockJobFunc)

			dbMock.SetMethodResponse("ListExpiredSchedules", []*Job{&mockJob}, nil)

			process()

			assert.True(t, dbMock.Method("ListExpiredSchedules").CalledOnce())
			assert.True(t, dbMock.Method("SaveJob").CalledOnce())
			assert.True(t, dbMock.Method("SaveJob").CalledWith(mockJob))
			assert.True(t, mockJob.IsDone())
			assert.False(t, loggerMock.Called())
		})
		t.Run("Should set the job as done if the job schedule is RECURRENT, but the limit date is lesser than now", func(t *testing.T) {
			dbMock, loggerMock := mockDependencies()

			mockJobName := "MYMOCKJOB!"
			mockJobFunc := func(j *Job) error {
				return nil
			}

			ld := time.Now().Add(-time.Hour)
			mockJob := Job{
				Name:              mockJobName,
				ScheduleType:      RECURRENT,
				ScheduleLimitDate: &ld,
			}

			Define(mockJobName, mockJobFunc)

			dbMock.SetMethodResponse("ListExpiredSchedules", []*Job{&mockJob}, nil)

			process()

			assert.True(t, dbMock.Method("ListExpiredSchedules").CalledOnce())
			assert.True(t, dbMock.Method("SaveJob").CalledOnce())
			assert.True(t, dbMock.Method("SaveJob").CalledWith(mockJob))
			assert.True(t, mockJob.IsDone())
			assert.False(t, loggerMock.Called())
		})
		t.Run("Should log an error and fail the job if the job schedule is RECURRENT, but the job has no scheduleString", func(t *testing.T) {
			dbMock, loggerMock := mockDependencies()

			mockJobName := "MYMOCKJOB!"
			mockJobFunc := func(j *Job) error {
				return nil
			}

			mockJob := Job{
				Name:           mockJobName,
				ScheduleType:   RECURRENT,
				ScheduleString: "",
			}

			Define(mockJobName, mockJobFunc)

			dbMock.SetMethodResponse("ListExpiredSchedules", []*Job{&mockJob}, nil)

			process()

			assert.True(t, dbMock.Method("ListExpiredSchedules").CalledOnce())
			assert.True(t, dbMock.Method("SaveJob").CalledOnce())
			assert.True(t, dbMock.Method("SaveJob").CalledWith(mockJob))
			assert.True(t, mockJob.HasFailed())
			assert.True(t, loggerMock.CalledWith("Tried to re-schedule recurrent job %s, but it had no ScheduleString"))
		})
		t.Run("Should log an error and fail the job if the job schedule is RECURRENT, but the job scheduleString is invalid", func(t *testing.T) {
			dbMock, loggerMock := mockDependencies()

			mockJobName := "MYMOCKJOB!"
			mockJobFunc := func(j *Job) error {
				return nil
			}

			mockJob := Job{
				Name:           mockJobName,
				ScheduleType:   RECURRENT,
				ScheduleString: "an invalid schedule string",
			}

			Define(mockJobName, mockJobFunc)

			dbMock.SetMethodResponse("ListExpiredSchedules", []*Job{&mockJob}, nil)

			process()

			assert.True(t, dbMock.Method("ListExpiredSchedules").CalledOnce())
			assert.True(t, dbMock.Method("SaveJob").CalledOnce())
			assert.True(t, dbMock.Method("SaveJob").CalledWith(mockJob))
			assert.True(t, mockJob.HasFailed())
			assert.True(t, loggerMock.CalledWith("Failed to get next schedule date for job %s, %v"))
		})
		t.Run("Should re-schedule job if the job is RECURRENT and its schedule string is valid", func(t *testing.T) {
			dbMock, loggerMock := mockDependencies()

			mockJobName := "MYMOCKJOB!"
			mockJobFunc := func(j *Job) error {
				return nil
			}

			mockJob := Job{
				Name:           mockJobName,
				ScheduleType:   RECURRENT,
				ScheduleString: "monday at 12:45",
			}

			Define(mockJobName, mockJobFunc)

			dbMock.SetMethodResponse("ListExpiredSchedules", []*Job{&mockJob}, nil)

			process()

			assert.True(t, dbMock.Method("ListExpiredSchedules").CalledOnce())
			assert.True(t, dbMock.Method("SaveJob").CalledOnce())
			assert.True(t, dbMock.Method("SaveJob").CalledWith(mockJob))
			assert.True(t, mockJob.IsPending())
			assert.False(t, loggerMock.Called())
		})
	})
}
