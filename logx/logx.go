package logx

import (
	"fmt"
	"io"
	"time"
)

type logFrame struct {
	Time  string `json:"time"`
	Level Level  `json:"level"`
}

type Logger struct {
	writers      []io.Writer
	defaultLevel Level
	timeformat   string
}

func NewLogger(writers ...io.Writer) (*Logger, error) {
	if len(writers) < 1 {
		return nil, fmt.Errorf("logx.NewLogger: %w", errNoWriter)
	}
	lg := &Logger{
		writers:      writers,
		defaultLevel: Info,
		timeformat:   time.RFC3339,
	}
	return lg, nil
}

func (l *Logger) SetDefaultLevel(level Level) {
	l.defaultLevel = level
}

func (l *Logger) SetTimeFormat(format string) {
	l.timeformat = format
}
