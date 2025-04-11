package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogConfig ...
type LogConfig struct {
	Level           string `json:"level" mapstructure:"level"`
	Mode            string `json:"mode" mapstructure:"mode"`
	Encoding        string `json:"encoding" mapstructure:"encoding"`
	StacktraceLevel string `json:"stacktrace_level" mapstructure:"stacktrace_level"`
}

const (
	production      = "production"
	development     = "development"
	jsonEncoding    = "json"
	consoleEncoding = "console"
)

func LogDefaultConfig() LogConfig {
	return LogConfig{
		Level:           "debug",
		Mode:            development,
		Encoding:        jsonEncoding,
		StacktraceLevel: "",
	}
}

var levelMap = map[string]zapcore.Level{
	"debug":  zap.DebugLevel,
	"info":   zap.InfoLevel,
	"warn":   zap.WarnLevel,
	"error":  zap.ErrorLevel,
	"dpanic": zap.DPanicLevel,
	"panic":  zap.PanicLevel,
	"fatal":  zap.FatalLevel,
}

func validateLogConfig(conf LogConfig) {
	_, existed := levelMap[conf.Level]
	if !existed {
		panic("log level must be one of: debug, info, warn, error, dpanic, panic, fatal")
	}
	if conf.Mode != development && conf.Mode != production {
		panic("log mode must be one of: development, production")
	}
	if conf.Encoding != consoleEncoding && conf.Encoding != jsonEncoding {
		panic("log encoding must be one of: console, json")
	}

	if conf.StacktraceLevel != "" {
		_, existed := levelMap[conf.StacktraceLevel]
		if !existed {
			panic("invalid stacktrace_level")
		}
	}
}

// NewLogger creates a logger
func NewLogger(conf LogConfig) *zap.Logger {
	validateLogConfig(conf)

	zapConf := zap.NewProductionConfig()
	if conf.Mode == development {
		zapConf = zap.NewDevelopmentConfig()
	}

	level := zap.NewAtomicLevelAt(levelMap[conf.Level])
	zapConf.Level = level
	zapConf.Encoding = conf.Encoding

	var options []zap.Option
	if conf.StacktraceLevel != "" {
		stackLevel := zap.NewAtomicLevelAt(levelMap[conf.StacktraceLevel])
		options = append(options, zap.AddStacktrace(stackLevel))
	}
	options = append(options, zap.AddCaller())

	logger, err := zapConf.Build(options...)
	if err != nil {
		panic(err)
	}
	return logger
}
