package scheduler

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// Config represents the configuration for the go-scheduler lib
type Config struct {
	// DB represents a user created struct that implements the JobDatabase interface.
	//
	// Users should provide this value when they want to have full controll over the
	// actions that the library will execute in its database.
	//
	// If no user created DB is specified, the MongoDB value should be provided
	DB JobDatabase

	// Logger represents a user created struct that implements the Logger interface.
	//
	// This logger will be used to log information when managing jobs.
	// If no logger is specified, the library will log nothing.
	Logger Logger

	// MongoDB represents the configuration values that the library need to start a job DB on a mongoDB connection.
	// This configuration uses the original mongoDB driver to do so.
	//
	// When providing this configuration, the library will have access to the user's mongoDB connection,
	// since it will be responsible to manage jobs.
	//
	// If this value is not specified, the DB value should be provided.
	MongoDB *MongoJobDBConfig

	// ProcessingRate represents the rate that the library will process jobs.
	//
	// Defaut: 1 minute
	ProcessingRate time.Duration

	// Location represents the location that the library should use when generating time values.
	//
	// Default: UTC
	Location string

	// DeleteOnDone defines if, when a job is done, the job should be deleted from the database.
	//
	// Default: false
	DeleteOnDone bool

	// DeleteOnCancel defines if, when a job is canceled, the job should be deleted from the database.
	//
	// Default: false
	DeleteOnCancel bool
}

type MongoJobDBConfig struct {
	// Conn its the mongo db connection that the user can provide.
	Conn *mongo.Client

	// DbName its the database name that the library should use to save jobs
	//
	// Default: go-scheduler
	DbName string

	// CollName its the collection name that the library should use to save jobs
	//
	// Default: scheduler-jobs
	CollName string
}
