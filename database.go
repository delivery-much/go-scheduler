package scheduler

// JobDatabase represents a database that can manipulate job documents
type JobDatabase interface {
	// InitJobDB its a function that will be called at the beggining of the library instantiation.
	//
	// It should start the job database an make it ready to read and write jobs.
	InitJobDB() error

	// ListExpiredSchedules should list jobs that are ready to run
	ListExpiredSchedules() ([]*Job, error)

	// List should list jobs given the Finder
	List(f Finder) ([]*Job, error)

	// SaveJob should save a job in its current state on the job database
	//
	// It should receive a job struct, and "upsert" it in the database. (If it's a new job, should insert a new job, if its an existent job, should update the existent job).
	SaveJob(j Job) error

	// DeleteJob should delete a job completely from the database
	DeleteJob(j Job) error
}
