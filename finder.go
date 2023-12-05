package scheduler

// Finder its a helper struct used to pass parameters to the database List action.
//
// Currently, only the status, name or data values can be used.
type Finder struct {
	Status string
	Name   string
	Data   map[string]any
}
