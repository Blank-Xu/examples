package db

import (
	"fmt"

	"go.uber.org/zap"
	"xorm.io/core"
)

type SimpleLogger struct {
	logger  *zap.Logger
	level   core.LogLevel
	showSQL bool
}

func NewSimpleLogger(log *zap.Logger, database string, logLevel core.LogLevel) *SimpleLogger {
	return &SimpleLogger{
		logger: log.With(zap.Field{Key: "database", String: database}),
		level:  logLevel,
	}
}

// Error implement core.ILogger
func (s *SimpleLogger) Error(v ...interface{}) {
	if s.level <= core.LOG_ERR && len(v) > 0 {
		switch e := v[0].(type) {
		case string:
			s.logger.Error(e)
		case error:
			s.logger.Error(e.Error())
		}
	}
}

// Errorf implement core.ILogger
func (s *SimpleLogger) Errorf(format string, v ...interface{}) {
	if s.level <= core.LOG_ERR {
		if len(v) > 0 {
			s.logger.Error(fmt.Sprintf(format, v...))
			return
		}
		s.logger.Error(format)
	}
}

// Debug implement core.ILogger
func (s *SimpleLogger) Debug(v ...interface{}) {
	if s.level <= core.LOG_DEBUG && len(v) > 0 {
		switch e := v[0].(type) {
		case string:
			s.logger.Debug(e)
		case error:
			s.logger.Debug(e.Error())
		}
	}
}

// Debugf implement core.ILogger
func (s *SimpleLogger) Debugf(format string, v ...interface{}) {
	if s.level <= core.LOG_DEBUG {
		if len(v) > 0 {
			s.logger.Debug(fmt.Sprintf(format, v...))
			return
		}
		s.logger.Debug(format)
	}
}

// Info implement core.ILogger
func (s *SimpleLogger) Info(v ...interface{}) {
	if s.level <= core.LOG_INFO && len(v) > 0 {
		switch e := v[0].(type) {
		case string:
			s.logger.Info(e)
		case error:
			s.logger.Info(e.Error())
		}
	}
}

// Infof implement core.ILogger
func (s *SimpleLogger) Infof(format string, v ...interface{}) {
	if s.level <= core.LOG_INFO {
		if len(v) > 0 {
			s.logger.Info(fmt.Sprintf(format, v...))
			return
		}
		s.logger.Info(format)
	}
}

// Warn implement core.ILogger
func (s *SimpleLogger) Warn(v ...interface{}) {
	if s.level <= core.LOG_WARNING && len(v) > 0 {
		switch e := v[0].(type) {
		case string:
			s.logger.Warn(e)
		case error:
			s.logger.Warn(e.Error())
		}
	}
}

// Warnf implement core.ILogger
func (s *SimpleLogger) Warnf(format string, v ...interface{}) {
	if s.level <= core.LOG_WARNING {
		if len(v) > 0 {
			s.logger.Warn(fmt.Sprintf(format, v...))
			return
		}
		s.logger.Warn(format)
	}
}

// Level implement core.ILogger
func (s *SimpleLogger) Level() core.LogLevel {
	return s.level
}

// SetLevel implement core.ILogger
func (s *SimpleLogger) SetLevel(l core.LogLevel) {
	s.level = l
}

// ShowSQL implement core.ILogger
func (s *SimpleLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		s.showSQL = true
		return
	}
	s.showSQL = show[0]
}

// IsShowSQL implement core.ILogger
func (s *SimpleLogger) IsShowSQL() bool {
	return s.showSQL
}
