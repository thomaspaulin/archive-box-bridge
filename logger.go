package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path"
	"sync"
)

var once sync.Once
var logger *log.Logger

func logFile() (*os.File, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	// todo figure out permissions for logging in /var/log/archive-box-bridge
	p := path.Join(u.HomeDir, "archive-box-bridge.log")
	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	return f, err
}

// Credit to https://stackoverflow.com/questions/18361750/correct-approach-to-global-logging-in-golang
func createLogger() *log.Logger {
	f, err := logFile()
	if err != nil {
		panic(fmt.Sprintf("Error opening log file: %v", err))
	} else {
		log.Printf("Logging to %s", f.Name())
	}
	mw := io.MultiWriter(f, os.Stdout)
	return log.New(mw, "[archive-box-bridge] ", log.Ldate|log.Ltime)
}

func GetLogger() *log.Logger {
	once.Do(func() {
		logger = createLogger()
	})
	return logger
}
