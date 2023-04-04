package log

import (
	"fmt"
	"os"
)

const (
	DEBUG = iota // 0
	INFO
	WARN
	ERROR
	FATAL
)

const DEFAULT_LEVEL = ERROR

type SimpleLogger struct {
	level int
}

func NewSimpleLogger(level int) *SimpleLogger {
	return &SimpleLogger{
		level: level,
	}
}

func (s *SimpleLogger) Debug(a ...any) {
	if s.level <= DEBUG {
		return
	}

	fmt.Println("DEBUG", a)
}

func (s *SimpleLogger) Info(a ...any) {
	if s.level <= INFO {
		return
	}

	fmt.Println("INFO", a)
}

func (s *SimpleLogger) Error(a ...any) {
	if s.level <= ERROR {
		return
	}

	fmt.Println("ERROR", a)
}

func (s *SimpleLogger) Warn(a ...any) {
	if s.level <= WARN {
		return
	}

	fmt.Println("WARN", a)
}

func (s *SimpleLogger) Fatal(a ...any) {
	if s.level <= FATAL {
		return
	}

	fmt.Println("FATAL", a)
	os.Exit(-1)
}
