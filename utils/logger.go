package utils

import (
	"fmt"
	"os"
	"time"
)

const MSG_INFO byte = 0x0
const MSG_WARNING byte = 0x1
const MSG_FATAL byte = 0x2
const MSG_VERBOSE byte = 0x3
const MSG_DEBUG byte = 0x4

var MSG_STRING = map[byte]string{
	MSG_INFO:    "INFO    ",
	MSG_WARNING: "WARNING ",
	MSG_FATAL:   "FATAL   ",
	MSG_VERBOSE: "VERBOSE ",
	MSG_DEBUG:   "DEBUG   ",
}

type Log struct {
	UseDebug        bool
	UseVerbose      bool
	UseTimestamp    bool
	TimestampFormat string
}

func (l Log) Message(log_level byte, message ...interface{}) {
	var msg string
	if l.UseTimestamp {
		if len(l.TimestampFormat) == 0 {
			l.TimestampFormat = time.RFC3339
		}
		timestamp := time.Now().Format(time.RFC3339)
		msg = timestamp + " " + MSG_STRING[log_level] + ":"
	} else {
		msg = MSG_STRING[log_level] + ":"
	}

	all := append([]interface{}{msg}, message...)

	fmt.Println(all...)
}

func (l Log) Info(message ...interface{}) {
	l.Message(MSG_INFO, message...)
}

func (l Log) Warning(message ...interface{}) {
	l.Message(MSG_WARNING, message...)
}

func (l Log) Fatal(message ...interface{}) {
	l.Message(MSG_FATAL, message...)
	os.Exit(1)
}

func (l Log) Verbose(message ...interface{}) {
	if l.UseDebug || l.UseVerbose {
		l.Message(MSG_VERBOSE, message...)
	}
}

func (l Log) Debug(message ...interface{}) {
	if l.UseDebug {
		l.Message(MSG_DEBUG, message...)
	}
}
