package log

var logger *SimpleLogger

func Debug(a ...any) {
	if logger == nil {
		logger = NewSimpleLogger(DEFAULT_LEVEL)
	}

	logger.Debug(a)
}

func Info(a ...any) {
	if logger == nil {
		logger = NewSimpleLogger(DEFAULT_LEVEL)
	}

	logger.Info(a)
}

func Error(a ...any) {
	if logger == nil {
		logger = NewSimpleLogger(DEFAULT_LEVEL)
	}

	logger.Error(a)
}

func Warn(a ...any) {
	if logger == nil {
		logger = NewSimpleLogger(DEFAULT_LEVEL)
	}

	logger.Warn(a)
}

func Fatal(a ...any) {
	if logger == nil {
		logger = NewSimpleLogger(DEFAULT_LEVEL)
	}

	logger.Fatal(a)
}
