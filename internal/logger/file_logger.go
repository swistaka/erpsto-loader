package logger

import (
	"fmt"
	"os"
	"time"
	//"path/filepath"
)

type FileLogger struct {
	filePath string
}

func NewFileLogger(path string) *FileLogger {
	return &FileLogger{filePath: path}
}

func (l *FileLogger) LogError(err error) {
	logEntry := fmt.Sprintf("[%s] ERROR: %v\n", time.Now().Format(time.RFC3339), err)
	l.writeLog(logEntry)
}

func (l *FileLogger) writeLog(message string) {
	file, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	file.WriteString(message)
}
