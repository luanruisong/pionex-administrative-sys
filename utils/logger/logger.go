package logger

import (
	"context"
	"errors"
	"io"
	"os"
	"pionex-administrative-sys/utils/app"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	gormlogger "gorm.io/gorm/logger"
)

var (
	log   *zap.Logger
	sugar *zap.SugaredLogger
)

func init() {
	// 获取日志目录

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochMillisTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var writer io.Writer = os.Stdout
	if app.FileLog() {
		// 使用 lumberjack 进行日志轮转
		writer = &lumberjack.Logger{
			Filename:   app.LogPath(time.Now().Format("2006-01-02.log")),
			MaxSize:    100, // MB
			MaxBackups: 30,
			MaxAge:     30, // 天
			Compress:   true,
		}

	}
	log = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(writer),
		zapcore.InfoLevel,
	), zap.AddCaller(), zap.AddCallerSkip(1))
	sugar = log.Sugar()
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

func With(fields ...zap.Field) *zap.Logger {
	return log.With(fields...)
}

func Sync() error {
	return log.Sync()
}

// Sugar Logger 封装

func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	sugar.Infow(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	sugar.Errorw(msg, keysAndValues...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	sugar.Debugw(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	sugar.Warnw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	sugar.Fatalw(msg, keysAndValues...)
}

// GetLogger 获取原始 zap.Logger
func GetLogger() *zap.Logger {
	return log
}

// GetSugar 获取 SugaredLogger
func GetSugar() *zap.SugaredLogger {
	return sugar
}

// Writer 实现 io.Writer 接口，用于接管第三方库的日志输出
type Writer struct {
	level zapcore.Level
}

func (w *Writer) Write(p []byte) (n int, err error) {
	msg := string(p)
	// 去除末尾换行符
	if len(msg) > 0 && msg[len(msg)-1] == '\n' {
		msg = msg[:len(msg)-1]
	}
	switch w.level {
	case zapcore.DebugLevel:
		log.Debug(msg)
	case zapcore.InfoLevel:
		log.Info(msg)
	case zapcore.WarnLevel:
		log.Warn(msg)
	case zapcore.ErrorLevel:
		log.Error(msg)
	default:
		log.Info(msg)
	}
	return len(p), nil
}

// NewWriter 创建指定级别的 Writer
func NewWriter(level zapcore.Level) *Writer {
	return &Writer{level: level}
}

// InfoWriter 返回 Info 级别的 Writer
func InfoWriter() *Writer {
	return &Writer{level: zapcore.InfoLevel}
}

// ErrorWriter 返回 Error 级别的 Writer
func ErrorWriter() *Writer {
	return &Writer{level: zapcore.ErrorLevel}
}

// GormLogger GORM 日志适配器
type GormLogger struct {
	SlowThreshold time.Duration
	LogLevel      gormlogger.LogLevel
}

// NewGormLogger 创建 GORM 日志适配器
func NewGormLogger() *GormLogger {
	return &GormLogger{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      gormlogger.Info,
	}
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		sugar.Infof(msg, data...)
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		sugar.Warnf(msg, data...)
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		sugar.Errorf(msg, data...)
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && !errors.Is(err, gormlogger.ErrRecordNotFound):
		log.Error("gorm",
			zap.Error(err),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		log.Warn("gorm slow query",
			zap.Duration("elapsed", elapsed),
			zap.Duration("threshold", l.SlowThreshold),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case l.LogLevel >= gormlogger.Info:
		log.Info("gorm",
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	}
}

func (l *GormLogger) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	return sql, params
}
