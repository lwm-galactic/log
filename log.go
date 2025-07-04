package log

import (
	"fmt"
	"go.uber.org/zap"
	"sync"
)

// ILogger Logger represents the ability to log messages, both errors and not.
type ILogger interface {
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)

	Infof(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Panicf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})

	Infow(msg string, keysAndValues ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	WithValues(keysAndValues ...interface{}) ILogger
	SetOptions(opts ...Option)
	Build() error
	//WithName(name string) ILogger
	//WithContext(ctx context.Context) context.Context
	//Write(p []byte) (n int, err error)
	Flush()
}
type Logger struct {
	opt       *Options
	logger    *zap.Logger
	mu        sync.Mutex
	entryPool *sync.Pool
}

func NewLoggerByOption(options *Options) ILogger {
	return &Logger{
		opt: options,
	}
}
func NewLogger() ILogger {
	return &Logger{
		opt: NewDefaultOptions(),
	}
}

func (l *Logger) Build() error {
	logger, err := l.opt.build()
	if err != nil {
		return err
	}
	l.logger = logger
	return nil
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger.Sugar().Infof(format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logger.Sugar().Debugf(format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logger.Sugar().Warnf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Sugar().Errorf(format, v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Sugar().Fatalf(format, v...)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.logger.Sugar().Panicf(format, v...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Infow(msg, keysAndValues...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Debugw(msg, keysAndValues...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Warnw(msg, keysAndValues...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Errorw(msg, keysAndValues...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Fatalw(msg, keysAndValues...)
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.logger.Sugar().Panicw(msg, keysAndValues...)
}

func (l *Logger) WithValues(keysAndValues ...interface{}) ILogger {
	if len(keysAndValues) == 0 {
		return l
	}

	fields := make([]zap.Field, 0, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		key := keysAndValues[i]
		var val interface{} = "MISSING"
		if i+1 < len(keysAndValues) {
			val = keysAndValues[i+1]
		}
		fields = append(fields, zap.Any(fmt.Sprintf("%v", key), val))
	}

	newLogger := *l
	newLogger.logger = l.logger.With(fields...)
	return &newLogger
}
func (l *Logger) Flush() {
	_ = l.logger.Sync()
}
