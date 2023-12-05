package scheduler

import "github.com/delivery-much/mock-helper/mock"

type loggerMock struct {
	mock.Mock
}

func newLoggerMock() *loggerMock {
	return &loggerMock{
		mock.NewMock(),
	}
}

func (lm *loggerMock) Error(message string) {
	lm.RegisterMethodCall("Error", message)
}

func (lm *loggerMock) Errorf(format string, a ...any) {
	lm.RegisterMethodCall("Errorf", format, a)
}
