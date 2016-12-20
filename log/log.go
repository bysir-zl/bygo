package log

import (
	"log"
	"fmt"
	"os"
)

const (
	LVerbose = iota + 1
	LDebug
	LInfo
	LWarn
	LError
	LAssert
)

var maxLevel = 0
var callDepth = 3

var Logger *log.Logger

func Verbose(tag string, v ...interface{}) {
	OutPutStd(callDepth,LVerbose, tag, v...)
}

func Debug(tag string, v ...interface{}) {
	OutPutStd(callDepth,LDebug, tag, v...)
}
func Info(tag string, v ...interface{}) {
	OutPutStd(callDepth,LInfo, tag, v...)
}

func Warn(tag string, v ...interface{}) {
	OutPutStd(callDepth,LWarn, tag, v...)
}
func Error(tag string, v ...interface{}) {
	OutPutStd(callDepth,LError, tag, v...)
}

func Assert(tag string, v ...interface{}) {
	OutPutStd(callDepth,LAssert, tag, v...)
}

func OutPutStd(callDepth,level int, tag string, v ...interface{}) {
	if maxLevel < level {
		l := ""
		switch level {
		case LVerbose:
			l = "[V]"
		case LInfo:
			l = "[I]"
		case LDebug:
			l = "[D]"
		case LWarn:
			l = "[W]"
		case LError:
			l = "[E]"
		case LAssert:
			l = "[A]"
		}
		tag = "[" + tag + "] "
		Logger.SetPrefix(l+" ")
		Logger.Output(callDepth, fmt.Sprintln(append([]interface{}{tag}, v...)...))
	}
}

func SetLogLevel(level int) {
	maxLevel = level
}

func SetCallDepth(d int) {
	callDepth = d
}

func init() {
	Logger = log.New(os.Stderr, "", log.Ltime | log.Ldate | log.Lshortfile)
}