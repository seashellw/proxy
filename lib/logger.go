package lib

import (
	"encoding/csv"
	"os"
	"sync"
	"time"
)

func NewLogger(fileName string, length int) *Logger {
	return &Logger{
		FileName: fileName,
		length:   length,
		lock:     &sync.Mutex{},
		list:     [][]string{},
	}
}

type Logger struct {
	FileName string
	length   int
	list     [][]string
	lock     *sync.Mutex
}

func (logger *Logger) Write(msg []string) {
	logger.lock.Lock()
	defer logger.lock.Unlock()
	newLine := append([]string{time.Now().Format("2006-01-02 15:04:05")}, msg...)
	logger.list = append(logger.list, newLine)
	file, err := os.OpenFile(logger.FileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Write(newLine)
	writer.Flush()

	if len(logger.list) >= 2*logger.length {
		logger.list = logger.list[logger.length:]
		file, err := os.OpenFile(logger.FileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			return
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		writer.WriteAll(logger.list)
	}
}

func (logger *Logger) Read(start int, end int) [][]string {
	logger.lock.Lock()
	defer logger.lock.Unlock()
	return logger.list[start:end]
}

func (logger *Logger) Error(msg []string) {
	go logger.Write(append([]string{"error"}, msg...))
}

func (logger *Logger) Info(msg []string) {
	go logger.Write(append([]string{"info"}, msg...))
}
