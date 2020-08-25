package troops

type ILogger interface {
	Infof(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}
