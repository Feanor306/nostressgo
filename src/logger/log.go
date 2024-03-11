package logger

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
)

type Log struct {
	zerolog.Logger
}

var instantiated *Log
var once sync.Once

func New() *Log {
	once.Do(func() {
		instantiated = &Log{
			Logger: zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger(),
		}
	})
	return instantiated
}
