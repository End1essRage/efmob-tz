package logger

import (
	"context"
	"sync"

	"github.com/go-chi/chi/v5/middleware"
	l "github.com/sirupsen/logrus"
)

var (
	instance *Wrapper
	once     sync.Once
)

type Wrapper struct {
	logger *l.Logger
}
type LogOptions struct {
	Pkg  string
	Func string
	Ctx  context.Context
}

// получить инстанс логгера
func Logger() *Wrapper {
	if instance == nil {
		return New("-", true, true)
	}
	return instance
}

// конфигурация нового логгера
func New(serviceName string, jsonFormat bool, debug bool) *Wrapper {
	once.Do(func() {
		instance = &Wrapper{
			logger: newLogger(serviceName, jsonFormat, debug),
		}
	})
	return instance
}

// вывод сообщения
// Универсальный метод для логирования
func (log *Wrapper) WithFields(opts LogOptions) *l.Entry {
	entry := log.logger.WithField("@caller", opts.Pkg).WithField("@func", opts.Func)

	if opts.Ctx != nil {
		reqID := middleware.GetReqID(opts.Ctx)
		if reqID != "" {
			entry = entry.WithField("request_id", reqID)
		}
	}

	return entry
}

// Для обратной совместимости
func (log *Wrapper) Log(pkg, fn string) *l.Entry {
	return log.WithFields(LogOptions{Pkg: pkg, Func: fn})
}

// вывод сообщения без параметров -> Log(pkg, fn string)
func (log *Wrapper) Logger() *l.Logger {
	return log.logger
}

func newLogger(serviceName string, jsonFormat bool, debug bool) *l.Logger {
	//переименовываем ключевы поля для elk
	fieldMap := l.FieldMap{
		l.FieldKeyTime:  "@timestamp",
		l.FieldKeyLevel: "@level",
		l.FieldKeyMsg:   "@message",
		l.FieldKeyFunc:  "@caller",
	}

	//создаем логгер
	logger := l.New()
	if jsonFormat {
		logger.SetFormatter(&l.JSONFormatter{
			//задаем формат вывода для elk
			TimestampFormat: "01/02/2006 15:04:05",
			//передаем мапу с неймингом полей
			FieldMap: fieldMap,
		})
	} else {
		logger.SetFormatter(&l.TextFormatter{})
	}

	//уровень логов
	if debug {
		logger.Level = l.DebugLevel
	}

	//записываем название сервиса
	logger.AddHook(&MetaHook{Name: serviceName})

	return logger
}
