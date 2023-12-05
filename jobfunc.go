package scheduler

// JobFunc represents a function that can handle jobs
type JobFunc func(*Job) (err error)
