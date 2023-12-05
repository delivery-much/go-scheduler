package scheduler

import "errors"

// emptyDB represents an empty job database that always returns an error
type emptyDB struct{}

func (edb *emptyDB) InitJobDB() error {
	return errors.New("Tried to access the scheduler DB to manage jobs, but the go-scheduler library was not instantiated")
}

func (edb *emptyDB) List(f Finder) ([]*Job, error) {
	return []*Job{}, errors.New("Tried to access the scheduler DB to manage jobs, but the go-scheduler library was not instantiated")
}

func (edb *emptyDB) ListExpiredSchedules() ([]*Job, error) {
	return []*Job{}, errors.New("Tried to access the scheduler DB to manage jobs, but the go-scheduler library was not instantiated")
}

func (edb *emptyDB) SaveJob(j Job) error {
	return errors.New("Tried to access the scheduler DB to manage jobs, but the go-scheduler library was not instantiated")
}

func (edb *emptyDB) DeleteJob(j Job) error {
	return errors.New("Tried to access the scheduler DB to manage jobs, but the go-scheduler library was not instantiated")
}
