package logrus

import (
	"bytes"
	"context"
	"sync"
	"time"
)

var (
	bufferPool *sync.Pool
)

func init() {
	bufferPool = &sync.Pool{New: func() any { return new(bytes.Buffer) }}

}

// Fields Custom field value map
type Fields map[string]any

// Level type hierarchy
type Level uint32

// Entry type
type Entry struct {
	Logger *Logger
	Data   Fields
	Time   time.Time
	Level  Level
	// Message passed to Trace, Debug, Info, Warn, Error, Fatal or Panic
	Message string
	Context context.Context
	// err may contain a field formatting error
	err string
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
		Data:   make(Fields, 6),
	}
}

// Dup Deep copy once
func (entry *Entry) Dup() *Entry {
	data := make(Fields, len(entry.Data))
	for k, v := range entry.Data {
		data[k] = v
	}
	return &Entry{Logger: entry.Logger, Data: data, Time: entry.Time, Context: entry.Context, err: entry.err}
}
