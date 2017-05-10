package log

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	LVerbose = iota + 1
	LDebug
	LInfo
	LWarn
	LError
	LAssert
)

type Logger struct {
	l         *log.Logger
	callDepth int
	maxLevel  int
}

var defLogger *Logger

func NewLogger() *Logger {
	return &Logger{
		l:         log.New(os.Stderr, "", log.Ltime|log.Ldate|log.Lshortfile),
		callDepth: 3,
		maxLevel:  0,
	}
}

func (p *Logger) Verbose(tag string, v ...interface{}) {
	p.OutPutStd(LVerbose, tag, v...)
}
func (p *Logger) Debug(tag string, v ...interface{}) {
	p.OutPutStd(LDebug, tag, v...)
}
func (p *Logger) Info(tag string, v ...interface{}) {
	p.OutPutStd(LInfo, tag, v...)
}
func (p *Logger) Warn(tag string, v ...interface{}) {
	p.OutPutStd(LWarn, tag, v...)
}
func (p *Logger) Error(tag string, v ...interface{}) {
	p.OutPutStd(LError, tag, v...)
}
func (p *Logger) Assert(tag string, v ...interface{}) {
	p.OutPutStd(LAssert, tag, v...)
}

func (p *Logger) OutPutStd(level int, tag string, v ...interface{}) {
	if p.maxLevel < level {
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
		p.l.SetPrefix(l + " ")

		var arg []interface{}
		var format string

		// 支持format
		if len(v) >= 2 {
			if _format, ok := v[0].(string); ok {
				if strings.Contains(_format, "%") {
					format = _format
					arg = append([]interface{}{tag}, v[1:]...)
				}
			}
		}

		if format == "" {
			arg = append([]interface{}{tag}, v...)
			format = strings.Repeat("%+v ", len(arg))
		}

		p.l.Output(p.callDepth, fmt.Sprintf(format, arg...))
	}
}

func (p *Logger) SetLogLevel(level int) {
	p.maxLevel = level
}

func (p *Logger) SetCallDepth(d int) {
	p.callDepth = d
}

func Verbose(tag string, v ...interface{}) {
	defLogger.Verbose(tag, v...)
}
func Debug(tag string, v ...interface{}) {
	defLogger.OutPutStd(LDebug, tag, v...)
}
func Info(tag string, v ...interface{}) {
	defLogger.OutPutStd(LInfo, tag, v...)
}
func Warn(tag string, v ...interface{}) {
	defLogger.OutPutStd(LWarn, tag, v...)
}
func Error(tag string, v ...interface{}) {
	defLogger.OutPutStd(LError, tag, v...)
}
func Assert(tag string, v ...interface{}) {
	defLogger.OutPutStd(LAssert, tag, v...)
}

func SetLogLevel(level int) {
	defLogger.SetLogLevel(level)
}

func SetCallDepth(d int) {
	defLogger.SetLogLevel(d)
}

func init() {
	defLogger = NewLogger()
}
