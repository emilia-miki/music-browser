package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

const logFileName string = "integration_test.log"

var (
	Error       *log.Logger
	Warning     *log.Logger
	Info        *log.Logger
	Stdout      *log.Logger
	StdoutRed   *log.Logger
	StdoutGreen *log.Logger
)

type coloredWriter struct {
	writer    io.Writer
	ansiColor []byte
}

func (cw *coloredWriter) Write(p []byte) (int, error) {
	cw.writer.Write(cw.ansiColor)
	return cw.writer.Write(p)
}

func newColoredWriter(writer io.Writer, ansiColor string) *coloredWriter {
	return &coloredWriter{
		writer:    writer,
		ansiColor: []byte(ansiColor),
	}
}

func Initialize(verbose bool) {
	logFile, err := os.Create(logFileName)
	if err != nil {
		panic(fmt.Sprintf(
			"Unable to create log file %s: %s",
			logFileName, err,
		))
	}

	stdoutWriter := newColoredWriter(os.Stdout, "\033[0m")
	stdoutRedWriter := newColoredWriter(os.Stdout, "\033[31m")
	stdoutGreenWriter := newColoredWriter(os.Stdout, "\033[32m")
	stdoutYellowWriter := newColoredWriter(os.Stdout, "\033[33m")

	stdoutLogWriter := io.MultiWriter(logFile, stdoutWriter)
	stdoutRedLogWriter := io.MultiWriter(logFile, stdoutRedWriter)
	stdoutGreenLogWriter := io.MultiWriter(logFile, stdoutGreenWriter)

	errLogWriter := io.MultiWriter(logFile, stdoutRedWriter)
	var warningLogWriter io.Writer = logFile
	var infoLogWriter io.Writer = logFile
	if verbose {
		warningLogWriter = io.MultiWriter(logFile, stdoutYellowWriter)
		infoLogWriter = io.MultiWriter(logFile, stdoutWriter)
	}

	Stdout = log.New(stdoutLogWriter, "", 0)
	StdoutRed = log.New(stdoutRedLogWriter, "", 0)
	StdoutGreen = log.New(stdoutGreenLogWriter, "", 0)

	Error = log.New(errLogWriter, "ERROR: ", log.Flags())
	Warning = log.New(warningLogWriter, "WARNING: ", log.Flags())
	Info = log.New(infoLogWriter, "INFO: ", log.Flags())
}

func Inform() {
	fmt.Println("\033[0mMore detailed logs written to file " + logFileName)
}
