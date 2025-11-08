package kit

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/ncraft-io/ncraft/go/pkg/ncraft/logs"
)

type kitLogger struct {
	*logs.SugaredLogger
}

func Logger() log.Logger {
	return NewKitLogger(logs.Logger())
}

func NewKitLogger(zap *logs.SugaredLogger) log.Logger {
	return &kitLogger{zap}
}

// Log err & level should be the first key
// msg may be first or only after the level
func (l kitLogger) Log(keyvals ...interface{}) error {
	length := len(keyvals)
	secondMsgField := func(l int, key interface{}) bool { return l >= 4 && fmt.Sprint(key) == "msg" }

	if length >= 2 {
		key := fmt.Sprint(keyvals[0])
		switch key {
		case "err":
			l.SugaredLogger.Errorw("", keyvals...)
		case "msg":
			l.SugaredLogger.Infow(fmt.Sprint(keyvals[1]), keyvals[2:]...)
		case "level":
			level := fmt.Sprint(keyvals[1])
			switch level {
			case "debug":
				if secondMsgField(length, keyvals[2]) {
					l.SugaredLogger.Debugw(fmt.Sprint(keyvals[3]), keyvals[4:]...)
				} else {
					l.SugaredLogger.Debugw("", keyvals[2:]...)
				}
			default:
				fallthrough
			case "info":
				if secondMsgField(length, keyvals[2]) {
					l.SugaredLogger.Infow(fmt.Sprint(keyvals[3]), keyvals[4:]...)
				} else {
					l.SugaredLogger.Infow("", keyvals[2:]...)
				}
			case "warn":
				if secondMsgField(length, keyvals[2]) {
					l.SugaredLogger.Warnw(fmt.Sprint(keyvals[3]), keyvals[4:]...)
				} else {
					l.SugaredLogger.Warnw("", keyvals[2:]...)
				}
			case "error":
				if secondMsgField(length, keyvals[2]) {
					l.SugaredLogger.Errorw(fmt.Sprint(keyvals[3]), keyvals[4:]...)
				} else {
					l.SugaredLogger.Errorw("", keyvals[2:]...)
				}
			}
		default:
			l.SugaredLogger.Infow("", keyvals...)
		}
	}

	return nil
}
