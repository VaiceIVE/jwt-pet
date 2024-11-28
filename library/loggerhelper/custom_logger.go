package loggerhelper

import (
	"go.uber.org/zap"
)

type CustomLogger struct {
	zapLogger *zap.Logger
}

func NewCustomLogger(options ...zap.Option) *CustomLogger {
	zapLogger, err := zap.NewProduction(options...)
	defer zapLogger.Sync()
	zap.ReplaceGlobals(zapLogger)

	if err != nil {
		zap.S().Fatalf("configure logger error %v", err)
	}

	return &CustomLogger{zapLogger: zapLogger}
}
func (l *CustomLogger) WithTracing() *zap.Logger {
	return l.zapLogger.With()
}

func (l *CustomLogger) NoTracing() *zap.Logger {
	return l.zapLogger.With()
}

func (l *CustomLogger) SugarWithTracing() *zap.SugaredLogger {
	return l.zapLogger.With().Sugar()
}

func (l *CustomLogger) SugarNoTracing() *zap.SugaredLogger {
	return l.zapLogger.With().Sugar()
}