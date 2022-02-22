package logger

import (
	"log"
	"os"
)

//LOGGER
var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

//InitLoggers: Initialize a logger for errors and for information without prefix
func InitLoggers() {
	InfoLogger = log.New(os.Stdout, "", 0)
	ErrorLogger = log.New(os.Stderr, "", 0)
}

//AddFilenameAndLinePrefix: add filename and line in logger prefixS
func AddFilenameAndLinePrefix(logger *log.Logger) {
	logger.SetFlags(log.Lshortfile)
}
