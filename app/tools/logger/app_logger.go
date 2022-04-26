package logger

import (
	"github.com/sanctumlabs/curtz/app/pkg"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type appLogger struct {
	// name is the name of the logger that is published to log as a scope
	name string

	// log is the instance of the logrus logger
	log *logrus.Entry
}

func newLogger(name string) *appLogger {
	newLogger := logrus.New()

	// Log as JSON instead of the default ASCII formatter.
	newLogger.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	newLogger.SetOutput(os.Stdout)

	l := &appLogger{
		name: name,
		log: newLogger.WithFields(logrus.Fields{
			logFieldScope: name,
			logFieldType:  LogTypeLog,
		}),
	}

	l.EnableJSONOutput(defaultJSONOutput)
	return l
}

// EnableJSONOutput enables JSON formatted output log
func (l *appLogger) EnableJSONOutput(enabled bool) {
	var formatter logrus.Formatter

	fieldMap := logrus.FieldMap{
		// If time field name is conflicted, logrus adds "fields." prefix.
		// So rename to unused field @time to avoid the confliction.
		logrus.FieldKeyTime:  logFieldTimeStamp,
		logrus.FieldKeyLevel: logFieldLevel,
		logrus.FieldKeyMsg:   logFieldMessage,
	}

	hostname, _ := os.Hostname()
	l.log.Data = logrus.Fields{
		logFieldScope:    l.log.Data[logFieldScope],
		logFieldType:     LogTypeLog,
		logFieldInstance: hostname,
		logFieldVer:      pkg.Version,
	}

	if enabled {
		formatter = &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap:        fieldMap,
		}
	} else {
		formatter = &logrus.TextFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap:        fieldMap,
		}
	}

	l.log.Logger.SetFormatter(formatter)
}

// SetAppID sets app_id field in the log. Default value is empty string
func (l *appLogger) SetAppID(id string) {
	//nolint SA4005
	l.log = l.log.WithField(logFieldAppID, id)
}

func toLogrusLevel(lvl LogLevel) logrus.Level {
	// ignore error because it will never happen
	l, _ := logrus.ParseLevel(string(lvl))
	return l
}

// SetOutputLevel sets log output level
func (l *appLogger) SetOutputLevel(outputLevel LogLevel) {
	l.log.Logger.SetLevel(toLogrusLevel(outputLevel))
}

// WithLogType specify the log_type field in log. Default value is LogTypeLog
func (l *appLogger) WithLogType(logType string) Logger {
	return &appLogger{
		name: l.name,
		log:  l.log.WithField(logFieldType, logType),
	}
}

// Info logs a message at level Info.
func (l *appLogger) Info(args ...interface{}) {
	l.log.Log(logrus.InfoLevel, args...)
}

// Infof logs a message at level Info.
func (l *appLogger) Infof(format string, args ...interface{}) {
	l.log.Logf(logrus.InfoLevel, format, args...)
}

// Debug logs a message at level Debug.
func (l *appLogger) Debug(args ...interface{}) {
	l.log.Log(logrus.DebugLevel, args...)
}

// Debugf logs a message at level Debug.
func (l *appLogger) Debugf(format string, args ...interface{}) {
	l.log.Logf(logrus.DebugLevel, format, args...)
}

// Warn logs a message at level Warn.
func (l *appLogger) Warn(args ...interface{}) {
	l.log.Log(logrus.WarnLevel, args...)
}

// Warnf logs a message at level Warn.
func (l *appLogger) Warnf(format string, args ...interface{}) {
	l.log.Logf(logrus.WarnLevel, format, args...)
}

// Error logs a message at level Error.
func (l *appLogger) Error(args ...interface{}) {
	l.log.Log(logrus.ErrorLevel, args...)
}

// Errorf logs a message at level Error.
func (l *appLogger) Errorf(format string, args ...interface{}) {
	l.log.Logf(logrus.ErrorLevel, format, args...)
}

// Fatal logs a message at level Fatal then the process will exit with status set to 1.
func (l *appLogger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

// Fatalf logs a message at level Fatal then the process will exit with status set to 1.
func (l *appLogger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}
