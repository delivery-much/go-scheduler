package scheduler

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// jobDocument represents a job document on the mongo database
type jobDocument struct {
	ID                *primitive.ObjectID `bson:"_id,omitempty"`
	Name              string              `bson:"name"`
	Data              map[string]any      `bson:"data"`
	ScheduleType      string              `bson:"schedule_type"`
	Status            string              `bson:"status"`
	NextRunAt         time.Time           `bson:"next_run_at"`
	LastRunAt         *time.Time          `bson:"last_run_at,omitempty"`
	ScheduleString    string              `bson:"schedule_string,omitempty"`
	ScheduleLimitDate *time.Time          `bson:"schedule_limit_date,omitempty"`
}

// marshalJob marshals a job struct into a job document
func marshalJob(j Job) jobDocument {
	var id *primitive.ObjectID
	if j.ID != "" {
		parsedID, err := primitive.ObjectIDFromHex(j.ID)
		if err == nil {
			id = &parsedID
		}
	}

	return jobDocument{
		ID:                id,
		Name:              j.Name,
		Data:              j.Data,
		ScheduleType:      j.ScheduleType.String(),
		Status:            j.Status.String(),
		NextRunAt:         j.NextRunAt,
		LastRunAt:         j.LastRunAt,
		ScheduleString:    j.ScheduleString,
		ScheduleLimitDate: j.ScheduleLimitDate,
	}
}

// unmarshalJob marshals a job document into a job struct
func unmarshalJob(j jobDocument) Job {
	return Job{
		ID:                j.ID.Hex(),
		Name:              j.Name,
		Data:              j.Data,
		ScheduleType:      ScheduleType(j.ScheduleType),
		Status:            ScheduleStatus(j.Status),
		NextRunAt:         j.NextRunAt,
		LastRunAt:         j.LastRunAt,
		ScheduleString:    j.ScheduleString,
		ScheduleLimitDate: j.ScheduleLimitDate,
	}
}
