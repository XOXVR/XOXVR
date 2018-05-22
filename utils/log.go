package utils

import (
	"log"
	"os"
)

// SetupLogFile sets out output for log
func SetupLogFile(logPath string) (err error) {
	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	log.SetOutput(f)
	log.Println("----------------------------------")
	log.Println("SESSION START")
	log.Println("----------------------------------")

	return
}
