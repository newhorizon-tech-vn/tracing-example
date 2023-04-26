package log

const (
	//DebugLevel has verbose message
	DebugLevel = "debug"
	//InfoLevel is default log level
	InfoLevel = "info"
	//WarnLevel is for log messages about possible issues
	WarnLevel = "warn"
	//ErrorLevel is for log errors
	ErrorLevel = "error"
	//FatalLevel is for log fatal messages. The sytem shutsdown after log the message.
	FatalLevel = "fatal"
	//PanicLevel logs a message, then panics.
	PanicLevel = "panic"
)
