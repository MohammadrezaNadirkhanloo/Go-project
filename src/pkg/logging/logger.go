package logging

import "github.com/MohammadrezaNadirkhanloo/Go-project/config"

type Logger interface {
	Init()

	Debug(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{})
	Debugf(template string, args ...interface{})

	Info(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{})
	Infof(template string, args ...interface{})

	Warn(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{})
	Warnf(template string, args ...interface{})

	Error(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{})
	Errorf(template string, args ...interface{})

	Fatal(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{})
	Fatalf(template string, args ...interface{})
	// Error(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{})
	// Errorf(template string, args ...interface{})

	// Fatal(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{})
	// Fatalf(template string, args ...interface{})
}

// file <- filebeat -> elasticsearch -> kibana
func NewLogger(cfg *config.Config) Logger {
	switch cfg.Logger.Logger {
	case "zap":
		return newZapLogger(cfg)
	case "zerolog":
		return newZeroLogger(cfg)
	}
	panic("logger not supported")
}
