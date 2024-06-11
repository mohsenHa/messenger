package zaplogger

import (
	"fmt"
	"github.com/mohsenHa/messenger/logger/loggerentity"
	"github.com/mohsenHa/messenger/pkg/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
)

type ZapLogger struct {
	config Config
	logger *zap.Logger
}

const (
	runtimeCallerSkip = 3
)

var (
	DriverName         = "zap"
	once               = sync.Once{}
	zapLogLevelMapping = map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"fatal": zapcore.FatalLevel,
	}
)

type Config struct {
	Level      string
	Filename   string
	Filepath   string
	LocalTime  bool
	MaxBackups int
	MaxSize    int
	MaxAge     int
}

func NewZapLogger(cfg Config) *ZapLogger {
	logger := &ZapLogger{config: cfg}
	logger.Init()
	return logger
}

func (zl *ZapLogger) getLogLevel() zapcore.Level {
	level, exists := zapLogLevelMapping[zl.config.Level]
	if !exists {
		level = zapcore.DebugLevel
	}
	return level
}

func (zl *ZapLogger) Init() {
	once.Do(func() {
		logger, _ := zap.NewProduction()

		config := zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		defaultEncoder := zapcore.NewJSONEncoder(config)
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   zl.config.Filename,
			LocalTime:  zl.config.LocalTime,
			MaxSize:    zl.config.MaxSize,    // megabytes
			MaxBackups: zl.config.MaxBackups, // megabytes
			MaxAge:     zl.config.MaxAge,     // days
		})

		stdOutWriter := zapcore.AddSync(os.Stdout)
		core := zapcore.NewTee(
			zapcore.NewCore(defaultEncoder, writer, zl.getLogLevel()),
			zapcore.NewCore(defaultEncoder, stdOutWriter, zl.getLogLevel()),
		)
		logger = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
		zl.logger = logger
	})
}

func (zl *ZapLogger) Debug(cat loggerentity.Category, sub loggerentity.SubCategory, msg string,
	extra map[loggerentity.ExtraKey]interface{}) {

	params := prepareLogInfo(cat, sub, extra)
	zl.logger.Debug(msg, params...)
}

func (zl *ZapLogger) Debugf(template string, args ...interface{}) {
	zl.logger.Debug(fmt.Sprintf(template, args...), zap.Any("trace", trace.Parse(runtimeCallerSkip)))
}

func (zl *ZapLogger) Info(cat loggerentity.Category, sub loggerentity.SubCategory, msg string,
	extra map[loggerentity.ExtraKey]interface{}) {

	params := prepareLogInfo(cat, sub, extra)
	zl.logger.Info(msg, params...)
}

func (zl *ZapLogger) Infof(template string, args ...interface{}) {
	zl.logger.Info(fmt.Sprintf(template, args...), zap.Any("trace", trace.Parse(runtimeCallerSkip)))
}

func (zl *ZapLogger) Warn(cat loggerentity.Category, sub loggerentity.SubCategory, msg string,
	extra map[loggerentity.ExtraKey]interface{}) {

	params := prepareLogInfo(cat, sub, extra)
	zl.logger.Warn(msg, params...)
}

func (zl *ZapLogger) Warnf(template string, args ...interface{}) {
	zl.logger.Warn(fmt.Sprintf(template, args...), zap.Any("trace", trace.Parse(runtimeCallerSkip)))
}

func (zl *ZapLogger) Error(cat loggerentity.Category, sub loggerentity.SubCategory, msg string,
	extra map[loggerentity.ExtraKey]interface{}) {

	params := prepareLogInfo(cat, sub, extra)
	zl.logger.Error(msg, params...)
}

func (zl *ZapLogger) Errorf(template string, args ...interface{}) {
	zl.logger.Error(fmt.Sprintf(template, args...))
}

func (zl *ZapLogger) Fatal(cat loggerentity.Category, sub loggerentity.SubCategory, msg string,
	extra map[loggerentity.ExtraKey]interface{}) {

	params := prepareLogInfo(cat, sub, extra)
	zl.logger.Fatal(msg, params...)
}

func (zl *ZapLogger) Fatalf(template string, args ...interface{}) {
	zl.logger.Fatal(fmt.Sprintf(template, args...))
}

func prepareLogInfo(cat loggerentity.Category, sub loggerentity.SubCategory,
	extra map[loggerentity.ExtraKey]interface{}) []zap.Field {

	if extra == nil {
		extra = make(map[loggerentity.ExtraKey]interface{})
	}
	extra["Category"] = cat
	extra["SubCategory"] = sub
	extra["trace"] = trace.Parse(0)

	return logParamsToZapParams(extra)
}
