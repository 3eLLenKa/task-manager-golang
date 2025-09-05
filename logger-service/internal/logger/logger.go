package logger

import (
	"log"
	"os"
)

type Logger struct {
	file *os.File
}

func New(path string) (*Logger, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, err
	}

	return &Logger{
		file: file,
	}, nil
}

func (l *Logger) Write(msg string) {
	log.SetOutput(l.file)
	log.Println(msg)
}

func (l *Logger) Close() error {
	return l.file.Close()
}
