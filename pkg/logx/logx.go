// Package logx provides ...
package logx

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sync"
)

var std = NewLogging()

// NewLogging log printer
func NewLogging(opts ...Option) *Logging {
	options := Options{
		level:  LevelInfo,
		writer: os.Stderr,
		flag:   stdFlags,
	}
	// apply opts
	for _, o := range opts {
		o(&options)
	}

	// new logging
	logging := &Logging{options: options}
	// object pool
	logging.entryPool.New = logging.newEntry
	logging.bufferPool.New = logging.newBuffer
	return logging
}

// Logging logging setup.
type Logging struct {
	mu         sync.Mutex
	entryPool  sync.Pool
	bufferPool sync.Pool

	options Options
}

// newEntry new entry
func (log *Logging) newEntry() interface{} {
	return &Entry{
		logging: log,
		Data:    make(map[string]interface{}, 4),
	}
}

// releaseEntry release entry
func (log *Logging) releaseEntry(e *Entry) {
	e.Data = map[string]interface{}{}
	log.entryPool.Put(e)
}

// newBuffer new buffer
func (log *Logging) newBuffer() interface{} {
	return new(bytes.Buffer)
}

// releaseBuffer release buffer
func (log *Logging) releaseBuffer(buf *bytes.Buffer) {
	buf.Reset()
	log.bufferPool.Put(buf)
}

// output output log
func (log *Logging) output(l Level, msg string) {
	if !log.V(l) {
		return
	}
	// get & put
	entry, _ := log.entryPool.Get().(*Entry)
	entry.Level = l
	entry.Message = msg

	calldepth := 4
	entry.Output(calldepth)
}

// Info logs to INFO log.
func (log *Logging) Info(args ...interface{}) {
	log.output(LevelInfo, fmt.Sprintln(args...))
}

// Info logs to INFO log.
func (log *Logging) Infof(format string, args ...interface{}) {
	log.output(LevelInfo, fmt.Sprintf(format, args...))
}

// Warning logs to WARNING log.
func (log *Logging) Warning(args ...interface{}) {
	log.output(LevelWarning, fmt.Sprintln(args...))
}

// Warning logs to WARNING log.
func (log *Logging) Warningf(format string, args ...interface{}) {
	log.output(LevelWarning, fmt.Sprintf(format, args...))
}

// Error logs to ERROR log.
func (log *Logging) Error(args ...interface{}) {
	log.output(LevelError, fmt.Sprintln(args...))
}

// Error logs to ERROR log.
func (log *Logging) Errorf(format string, args ...interface{}) {
	log.output(LevelError, fmt.Sprintf(format, args...))
}

// Fatal logs to ERROR log. with os.Exit(1).
func (log *Logging) Fatal(args ...interface{}) {
	log.output(LevelFatal, fmt.Sprintln(args...))
	os.Exit(1)
}

// Fatal logs to ERROR log. with os.Exit(1).
func (log *Logging) Fatalf(format string, args ...interface{}) {
	log.output(LevelFatal, fmt.Sprintf(format, args...))
	os.Exit(1)
}

// V reports whether verbosity level log is at least the requested verbose level.
func (log *Logging) V(l Level) bool {
	return log.options.level <= l
}

// Info logs to INFO log.
func Info(args ...interface{}) {
	std.Info(args...)
}

// Info logs to INFO log.
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

// Warning logs to WARNING log.
func Warning(args ...interface{}) {
	std.Warning(args...)
}

// Warning logs to WARNING log.
func Warningf(format string, args ...interface{}) {
	std.Warningf(format, args...)
}

// Error logs to ERROR log.
func Error(args ...interface{}) {
	std.Error(args...)
}

// Error logs to ERROR log.
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

// Fatal logs to ERROR log. with os.Exit(1).
func Fatal(args ...interface{}) {
	std.Fatal(args...)
}

// Fatal logs to ERROR log. with os.Exit(1).
func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}

// WithData custom data
func WithData(data map[string]interface{}) *Entry {
	entry, _ := std.entryPool.Get().(*Entry)
	entry.Data = data
	return entry
}

// WithContext context
func WithContext(ctx context.Context) *Entry {
	entry, _ := std.entryPool.Get().(*Entry)
	entry.context = ctx
	return entry
}
