// Adapted from go-algorand/logging

package logging

import (
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// Level refers to the log logging level
type Level uint32

var (
	_globalMu  sync.RWMutex
	baseLogger Logger
)

const (
	// Panic Level level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	Panic Level = iota
	// Fatal Level level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	Fatal
	// Error Level level. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	Error
	// Warn Level level. Non-critical entries that deserve eyes.
	Warn
	// Info Level level. General operational entries about what's going on inside the
	// application.
	Info
	// Debug Level level. Usually only enabled when debugging. Very verbose logging.
	Debug
)

const auxo = "auxo"
const stackPrefix = "[Stack]"

var once sync.Once

func Init() {
	once.Do(func() {
		// baseLogger makes use of the logrus implementation
		baseLogger = NewLogrusLogger()
	})
}

func init() {
	Init()
}

// Fields maps logrus fields
type Fields = logrus.Fields

// Logger is the interface for loggers.
type Logger interface {
	// Debug logs a message at level Debug.
	Debug(...interface{})
	Debugln(...interface{})
	Debugf(string, ...interface{})

	// Info logs a message at level Info.
	Info(...interface{})
	Infoln(...interface{})
	Infof(string, ...interface{})

	// Warn logs a message at level Warn.
	Warn(...interface{})
	Warnln(...interface{})
	Warnf(string, ...interface{})

	// Error logs a message at level Error.
	Error(...interface{})
	Errorln(...interface{})
	Errorf(string, ...interface{})

	// Fatal logs a message at level Fatal.
	Fatal(...interface{})
	Fatalln(...interface{})
	Fatalf(string, ...interface{})

	// Panic logs a message at level Panic.
	Panic(...interface{})
	Panicln(...interface{})
	Panicf(string, ...interface{})

	// Set the logging version (Info by default)
	SetLevel(Level)

	// source adds file, line and function fields to the event
	source() *logrus.Entry

	// Adds a hook to the logger
	AddHook(hook logrus.Hook)
}

type logger struct {
	entry *logrus.Entry
}

func (l *logger) With(key string, value interface{}) Logger {
	return &logger{l.entry.WithField(key, value)}
}

func (l *logger) Debug(args ...interface{}) {
	l.source().Debug(args...)
}

func (l *logger) Debugln(args ...interface{}) {
	l.source().Debugln(args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.source().Debugf(format, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.source().Info(args...)
}

func (l *logger) Infoln(args ...interface{}) {
	l.source().Infoln(args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.source().Infof(format, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.source().Warn(args...)
}

func (l *logger) Warnln(args ...interface{}) {
	l.source().Warnln(args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.source().Warnf(format, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().Error(args...)
}

func (l *logger) Errorln(args ...interface{}) {
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().Errorln(args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().Errorf(format, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().Fatal(args...)
}

func (l *logger) Fatalln(args ...interface{}) {
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().Fatalln(args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().Fatalf(format, args...)
}

func (l *logger) Panic(args ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().Panic(args...)
}

func (l *logger) Panicln(args ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().Panicln(args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()
	l.source().Errorln(stackPrefix, string(debug.Stack()))
	l.source().Panicf(format, args...)
}

func (l *logger) WithFields(fields Fields) Logger {
	return &logger{l.source().WithFields(fields)}
}

func (l *logger) SetLevel(lvl Level) {
	l.entry.Logger.Level = logrus.Level(lvl)
}

func (l *logger) IsLevelEnabled(level Level) bool {
	return l.entry.Logger.Level >= logrus.Level(level)
}

func (l *logger) SetOutput(w io.Writer) {
	l.setOutput(w)
}

func (l *logger) setOutput(w io.Writer) {
	l.entry.Logger.Out = w
}

func (l *logger) source() *logrus.Entry {
	event := l.entry
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "<???>"
		line = 1
		event = event.WithFields(logrus.Fields{
			"file": file,
			"line": line,
		})
	} else {
		// Add file name and number
		slash := strings.Index(file, auxo)
		file = file[slash:]
		event = event.WithFields(logrus.Fields{
			"file": file,
			"line": line,
		})
		// Add function name if possible
		if function := runtime.FuncForPC(pc); function != nil {
			slash := strings.LastIndex(function.Name(), "/")
			name := function.Name()[slash+1:]
			event = event.WithField("function", name)
		}
	}
	return event
}

func (l *logger) AddHook(hook logrus.Hook) {
	l.entry.Logger.Hooks.Add(hook)
}

func Base() Logger {
	_globalMu.RLock()
	defer _globalMu.RUnlock()
	l := baseLogger
	return l
}

func NewLogrusLogger() Logger {
	l := logrus.New()
	configureLogrusDefault(l)
	out := logger{logrus.NewEntry(l)}
	return &out
}

func configureLogrusDefault(l *logrus.Logger) {
	l.SetLevel(logrus.DebugLevel)
	l.SetOutput(os.Stdout)
	formatter := l.Formatter
	tf, ok := formatter.(*logrus.TextFormatter)
	if ok {
		tf.TimestampFormat = "2006-01-02T15:04:05.999999999Z07:00"
		tf.FullTimestamp = true
	}
}
