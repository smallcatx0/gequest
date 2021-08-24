package request

import (
	"log"
	"os"
)

// 兼容标准库log的接口
var _ Logger = log.New(os.Stdout, "", 0)

type Logger interface {
	Print(v ...interface{})
}

// 控制台日志
var ConsoleLog Logger = log.New(os.Stderr, "[request] ", log.Ldate|log.Ltime|log.Lshortfile)

func NewFileLogger(pathstr string) *log.Logger {
	file, err := os.OpenFile(pathstr, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil
	}
	logger := log.New(file, "[request] ", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}
