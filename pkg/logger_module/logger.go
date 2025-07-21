package logger_module

import (
	"io"
	"log"
	"os"
	"sync"
)

type Logger struct {
	*log.Logger
	mutex sync.Mutex
}

var (
	instance *Logger
	once     sync.Once
)

// Создаем логгер с указанными настройками (один раз из за once.DO)
func New(out io.Writer, prefix string, flag int) *Logger {
	once.Do(func() {
		instance = &Logger{
			Logger: log.New(out, prefix, flag),
		}
	})
	return instance
}

func Get() *Logger {
	if instance == nil {
		return New(os.Stdout, "[APP]", log.LstdFlags|log.Lshortfile)
	}
	return instance
}

func (logger *Logger) Info(v ...interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	logger.SetPrefix("[INFO]")
	logger.Println(v...)
}

func (logger *Logger) Debug(v ...interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	logger.SetPrefix("[DEBUG]")
	logger.Println(v...)
}

func (logger *Logger) Error(v ...interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	logger.SetPrefix("[Error]")
	logger.Println(v...)
}

func (logger *Logger) Fatal(v ...interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()
	logger.SetPrefix("[FATAL]")
	logger.Println(v...)
	os.Exit(1)
}
