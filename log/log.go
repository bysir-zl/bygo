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
	lErr      *log.Logger
	callDepth int
	maxLevel  int
	tag       string
}

var defLogger *Logger

func NewLogger() *Logger {
	return &Logger{
		lErr:      log.New(os.Stderr, "", log.Ltime|log.Ldate|log.Lshortfile),
		l:         log.New(os.Stdout, "", log.Ltime|log.Ldate|log.Lshortfile),
		callDepth: 3,
		maxLevel:  0,
	}
}

func (p *Logger) Verbose(v ...interface{}) {
	p.OutPutStd(LVerbose, p.tag, v...)
}
func (p *Logger) Debug(v ...interface{}) {
	p.OutPutStd(LDebug, p.tag, v...)
}
func (p *Logger) Info(v ...interface{}) {
	p.OutPutStd(LInfo, p.tag, v...)
}
func (p *Logger) Warn(v ...interface{}) {
	p.OutPutStd(LWarn, p.tag, v...)
}
func (p *Logger) Error(v ...interface{}) {
	p.OutPutStd(LError, p.tag, v...)
}
func (p *Logger) Assert(v ...interface{}) {
	p.OutPutStd(LAssert, p.tag, v...)
}

func (p *Logger) VerboseT(tag string, v ...interface{}) {
	p.OutPutStd(LVerbose, tag, v...)
}
func (p *Logger) DebugT(tag string, v ...interface{}) {
	p.OutPutStd(LDebug, tag, v...)
}
func (p *Logger) InfoT(tag string, v ...interface{}) {
	p.OutPutStd(LInfo, tag, v...)
}
func (p *Logger) WarnT(tag string, v ...interface{}) {
	p.OutPutStd(LWarn, tag, v...)
}
func (p *Logger) ErrorT(tag string, v ...interface{}) {
	p.OutPutStd(LError, tag, v...)
}
func (p *Logger) AssertT(tag string, v ...interface{}) {
	p.OutPutStd(LAssert, tag, v...)
}

func (p *Logger) SetTag(tag string) {
	p.tag = tag
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
		l += " "

		if tag != "" {
			tag = "[" + tag + "]"
		}

		var arg []interface{}
		var format string

		// 支持format
		if len(v) >= 2 {
			if _format, ok := v[0].(string); ok {
				if strings.Contains(_format, "%") {
					format = "%s " + _format
					if tag != "" {
						arg = append([]interface{}{tag}, v[1:]...)
					}else{
						arg = v[1:]
					}
				}
			}
		}

		if format == "" {
			if tag != "" {
				arg = append([]interface{}{tag}, v...)
			}else{
				arg = v
			}
			format = strings.Repeat("%+v ", len(arg))
		}

		if level == LError {
			p.lErr.SetPrefix(l)
			p.lErr.Output(p.callDepth, fmt.Sprintf(format, arg...))
		} else {
			p.l.SetPrefix(l)
			p.l.Output(p.callDepth, fmt.Sprintf(format, arg...))
		}
	}
}

func (p *Logger) SetLogLevel(level int) {
	p.maxLevel = level
}

func (p *Logger) SetCallDepth(d int) {
	p.callDepth = d
}

func Verbose(v ...interface{}) {
	defLogger.OutPutStd(LVerbose, defLogger.tag, v...)
}
func Debug(v ...interface{}) {
	defLogger.OutPutStd(LDebug, defLogger.tag, v...)
}
func Info(v ...interface{}) {
	defLogger.OutPutStd(LInfo, defLogger.tag, v...)
}
func Warn(v ...interface{}) {
	defLogger.OutPutStd(LWarn, defLogger.tag, v...)
}
func Error(v ...interface{}) {
	defLogger.OutPutStd(LError, defLogger.tag, v...)
}
func Assert(v ...interface{}) {
	defLogger.OutPutStd(LAssert, defLogger.tag, v...)
}

func VerboseT(tag string, v ...interface{}) {
	defLogger.OutPutStd(LVerbose, tag, v...)

}
func DebugT(tag string, v ...interface{}) {
	defLogger.OutPutStd(LDebug, tag, v...)
}
func InfoT(tag string, v ...interface{}) {
	defLogger.OutPutStd(LInfo, tag, v...)
}
func WarnT(tag string, v ...interface{}) {
	defLogger.OutPutStd(LWarn, tag, v...)
}
func ErrorT(tag string, v ...interface{}) {
	defLogger.OutPutStd(LError, tag, v...)
}
func AssertT(tag string, v ...interface{}) {
	defLogger.OutPutStd(LAssert, tag, v...)
}

func SetLogLevel(level int) {
	defLogger.SetLogLevel(level)
}

func init() {
	defLogger = NewLogger()
}
