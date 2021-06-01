package request

import (
	"log"
	"os"
)

type Logger interface {
	Print(string)
}

type ConsoleLog struct{}

func (l *ConsoleLog) Print(logstr string) {
	log.Print(logstr)
}

type FileLog struct {
	Path string
	log  *log.Logger
}

func NewFileLogger(pathstr string) *FileLog {
	file, err := os.OpenFile(pathstr, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil
	}
	logger := log.New(file, "request", log.Ldate|log.Ltime|log.Lshortfile)
	return &FileLog{
		Path: pathstr,
		log:  logger,
	}
}

func (fl *FileLog) Print(logstr string) {
	fl.log.Print(logstr)
}
