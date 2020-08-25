package troops

// type Job func()

// type Job func(h func(args ...interface{}) error, options ...interface{}) error

type Job struct {
	exec func(args ...interface{})
	args []interface{}
}
