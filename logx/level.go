package logx

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
	Fatal
)
