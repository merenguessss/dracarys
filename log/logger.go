package log

import (
	"fmt"
	"log"
	"os"
)

type Level int

// 日志等级
const (
	UNKNOWN Level = iota
	TRACE
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

// 日志颜色.
const (
	green    = "\033[1;97;42m"
	white    = "\033[1;90;47m"
	dracarys = "\033[1;97;75m"
	yellow   = "\033[1;90;43m"
	red      = "\033[1;97;41m"
	blue     = "\033[1;97;44m"
	cyan     = "\033[1;97;46m"
	reset    = "\033[0m"
)

// word颜色.
const (
	wGreen  = "\033[1;32;75m"
	wWhite  = "\033[1;37;75m"
	wYellow = "\033[1;33;75m"
	wRed    = "\033[1;31;75m"
	wBlue   = "\033[1;34;75m"
	wCyan   = "\033[1;36;75m"
)

// 默认calldepth。
const defaultCalldepth = 4

// 日志等级转化为字符串.
func (level Level) String() string {
	switch level {
	case TRACE:
		return "[TRACE]"
	case DEBUG:
		return "[DEBUG]"
	case INFO:
		return "[INFO]"
	case WARNING:
		return "[WARNING]"
	case ERROR:
		return "[ERROR]"
	case FATAL:
		return "[FATAL]"
	default:
		return "[UNKNOWN]"
	}
}

// 日志等级转化为字符串颜色.
func (level Level) WordColor() string {
	switch level {
	case TRACE:
		return wGreen
	case DEBUG:
		return wBlue
	case INFO:
		return wWhite
	case WARNING:
		return wYellow
	case ERROR:
		return wRed
	case FATAL:
		return wRed
	default:
		return wRed
	}
}

// 日志等级转换为颜色.
func (level Level) Color() string {
	switch level {
	case TRACE:
		return green
	case DEBUG:
		return blue
	case INFO:
		return white
	case WARNING:
		return yellow
	case ERROR:
		return red
	case FATAL:
		return red
	default:
		return cyan
	}
}

// Log 日志接口.
// 可以自行实现,实现后注入DefaultLogger即可使用自定义.
type Log interface {
	Trace(v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
	Tracef(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warningf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}

func Trace(v ...interface{}) {
	DefaultLogger.Trace(v...)
}

func Tracef(format string, v ...interface{}) {
	DefaultLogger.Tracef(format, v...)
}

func Debug(v ...interface{}) {
	DefaultLogger.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	DefaultLogger.Debugf(format, v...)
}

func Info(v ...interface{}) {
	DefaultLogger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	DefaultLogger.Infof(format, v...)
}

func Warning(v ...interface{}) {
	DefaultLogger.Warning(v...)
}

func Warningf(format string, v ...interface{}) {
	DefaultLogger.Warningf(format, v...)
}

func Error(v ...interface{}) {
	DefaultLogger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	DefaultLogger.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	DefaultLogger.Fatal(v...)
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	DefaultLogger.Fatalf(format, v...)
	os.Exit(1)
}

var DefaultLogger = New(DefaultOptions)

var DefaultOptions = &Options{
	Level: DEBUG,
}

type logger struct {
	*log.Logger
	opts *Options
}

type Options struct {
	Path  string `yaml:"path"`
	Level Level  `yaml:"level"`
}

var New = func(opts *Options) *logger {
	out := os.Stdout
	if opts.Path != "" {
		logFile, err := os.OpenFile(opts.Path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("open file, ERROR : ", err)
		}
		out = logFile
	}
	fmt.Println(opts.Path)

	return &logger{
		log.New(out, fmt.Sprintf("%s[DRACARYS]%s", dracarys, reset), log.LstdFlags|log.Lshortfile),
		opts,
	}
}

func (log *logger) Trace(v ...interface{}) {
	if log.opts.Level > TRACE {
		return
	}
	log.Logger.Output(defaultCalldepth, log.formatLog(TRACE, fmt.Sprint(v...)))
}

func (log *logger) Debug(v ...interface{}) {
	if log.opts.Level > DEBUG {
		return
	}
	log.Logger.Output(defaultCalldepth, log.formatLog(DEBUG, fmt.Sprint(v...)))
}

func (log *logger) Info(v ...interface{}) {
	if log.opts.Level > INFO {
		return
	}
	log.Logger.Output(defaultCalldepth, log.formatLog(INFO, fmt.Sprint(v...)))
}

func (log *logger) Warning(v ...interface{}) {
	if log.opts.Level > WARNING {
		return
	}
	log.Logger.Output(defaultCalldepth, log.formatLog(WARNING, fmt.Sprint(v...)))

}

func (log *logger) Error(v ...interface{}) {
	if log.opts.Level > ERROR {
		return
	}
	log.Logger.Output(defaultCalldepth, log.formatLog(ERROR, fmt.Sprint(v...)))

}

func (log *logger) Fatal(v ...interface{}) {
	if log.opts.Level > FATAL {
		return
	}
	log.Logger.Output(defaultCalldepth, log.formatLog(FATAL, fmt.Sprint(v...)))

}

func (log *logger) Tracef(format string, v ...interface{}) {
	if log.opts.Level > TRACE {
		return
	}
	log.Logger.Output(defaultCalldepth, log.formatLog(TRACE, fmt.Sprintf(format, v...)))
}

func (log *logger) Debugf(format string, v ...interface{}) {
	if log.opts.Level > DEBUG {
		return
	}
	log.Logger.Output(defaultCalldepth, log.formatLog(DEBUG, fmt.Sprintf(format, v...)))
}

func (log *logger) Infof(format string, v ...interface{}) {
	if log.opts.Level > INFO {
		return
	}
	log.Logger.Output(defaultCalldepth, log.formatLog(INFO, fmt.Sprintf(format, v...)))
}

func (log *logger) Warningf(format string, v ...interface{}) {
	if log.opts.Level > WARNING {
		return
	}
	log.Logger.Output(defaultCalldepth, log.formatLog(WARNING, fmt.Sprintf(format, v...)))
}

func (log *logger) Errorf(format string, v ...interface{}) {
	if log.opts.Level > ERROR {
		return
	}
	log.Logger.Output(defaultCalldepth, log.formatLog(ERROR, fmt.Sprintf(format, v...)))
}

func (log *logger) Fatalf(format string, v ...interface{}) {
	if log.opts.Level > FATAL {
		return
	}
	log.Logger.Output(defaultCalldepth, log.formatLog(FATAL, fmt.Sprintf(format, v...)))
}

// formatLog 格式化log.
func (log *logger) formatLog(level Level, data string) string {
	return fmt.Sprintf("%s%s%s%s %s%s%s", reset,
		level.Color(), level.String(), reset,
		level.WordColor(), data, reset)
}
