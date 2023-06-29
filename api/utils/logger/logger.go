package logger

import (
	"fmt"
	"log"
	"os"
)

const logFileName string = "music_browser.log"

var (
	Error *log.Logger
	Info  *log.Logger
)

func Initialize() {
	logFile, err := os.Create(logFileName)
	if err != nil {
		panic(fmt.Sprintf(
			"Unable to create log file %s: %s", logFileName, err))
	}

	Error = log.New(logFile, "\033[32mERROR: ", log.Flags())
	Info = log.New(logFile, "\033[0mINFO: ", log.Flags())
}

func Inform() {
	fmt.Println("Logs have been written to file " + logFileName)
}
