package goxash3d_fwgs

import (
	"fmt"
	"os"
	"sync"
	"time"
	"path/filepath"
)

var (
	logFile     *os.File
	logMutex    sync.Mutex
	logFilePath = "/cs16dedicated/goxash3d.log"
)

func InitFileLogger() {
	logMutex.Lock()
	defer logMutex.Unlock()

	// Автоматически создаем папку, если ее нет
	dir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("[LOG ERROR] Failed to create dir %s: %v\n", dir, err)
	}

	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("[LOG ERROR] Failed to open log file %s: %v\n", logFilePath, err)
		return
	}
	fmt.Printf("[LOG SUCCESS] Opened log file at %s\n", logFilePath)
}

func FileLog(level, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	timestamp := time.Now().Format("2006-01-02T15:04:05.000-07:00")
	entry := fmt.Sprintf("%s [%s] %s\n", timestamp, level, msg)

	logMutex.Lock()
	defer logMutex.Unlock()

	if logFile != nil {
		logFile.WriteString(entry)
		logFile.Sync() // Гарантируем мгновенную запись на диск
	}
}

func init() {
	InitFileLogger()
}

