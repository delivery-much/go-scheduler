package scheduler

import (
	"github.com/delivery-much/mock-helper/mock"
)

type databaseMock struct {
	mock.Mock
}

func newDatabaseMock() *databaseMock {
	return &databaseMock{
		mock.NewMock(),
	}
}

func (dm *databaseMock) InitJobDB() (err error) {
	dm.RegisterMethodCall("InitJobDB")

	res := dm.GetMethodResponse("InitJobDB")
	if len(res) == 0 {
		return
	}

	return res.GetError(0)
}

func (dm *databaseMock) List(f Finder) (js []*Job, err error) {
	dm.RegisterMethodCall("List", f)

	res := dm.GetMethodResponse("List")
	if len(res) == 0 {
		return
	}

	return res.Get(0).([]*Job), res.GetError(1)
}

func (dm *databaseMock) ListExpiredSchedules() (js []*Job, err error) {
	dm.RegisterMethodCall("ListExpiredSchedules")

	res := dm.GetMethodResponse("ListExpiredSchedules")
	if len(res) == 0 {
		return
	}

	return res.Get(0).([]*Job), res.GetError(1)
}

func (dm *databaseMock) SaveJob(j Job) (err error) {
	dm.RegisterMethodCall("SaveJob", j)

	res := dm.GetMethodResponse("SaveJob")
	if len(res) == 0 {
		return
	}

	return res.GetError(0)
}

func (dm *databaseMock) DeleteJob(j Job) (err error) {
	dm.RegisterMethodCall("DeleteJob")

	res := dm.GetMethodResponse("DeleteJob")
	if len(res) == 0 {
		return
	}

	return res.GetError(0)
}
