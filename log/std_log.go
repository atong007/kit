package log

import (
	"go.uber.org/zap"
)

type StdLog struct {
	log *zap.SugaredLogger
}

func New() *StdLog {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()
	zap.AddCallerSkip(1)
	return &StdLog{log: sugar}
}

func (l *StdLog) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *StdLog) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *StdLog) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *StdLog) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *StdLog) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *StdLog) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}
