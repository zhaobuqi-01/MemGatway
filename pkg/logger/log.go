package logger

import (
	"fmt"
	"gateway/configs"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.Logger

	logConfig *configs.LogConfig

	fileOutput      zapcore.WriteSyncer
	errorFileOutput zapcore.WriteSyncer
)

func Init() {
	logConfig = configs.GetLogConfig()
	initLogger()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC3339),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getFileLogger(Filename string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   Filename,
		MaxSize:    logConfig.MaxSize,
		MaxBackups: logConfig.MaxBackups,
		MaxAge:     logConfig.MaxAge,
		LocalTime:  true,
		Compress:   logConfig.Compress,
	}
}

func getLogLevel(logLevelStr string) zapcore.Level {
	var logLevel zapcore.Level

	switch logLevelStr {
	case "debug":
		logLevel = zap.DebugLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	case "dpanic":
		logLevel = zap.DPanicLevel
	case "panic":
		logLevel = zap.PanicLevel
	case "fatal":
		logLevel = zap.FatalLevel
	default:
		logLevel = zap.InfoLevel
	}

	return logLevel
}

func newCoreWithLevel(encoder zapcore.Encoder, output zapcore.WriteSyncer, level zapcore.Level) zapcore.Core {
	return zapcore.NewCore(encoder, output, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	}))
}

func getZapCore() zapcore.Core {
	encoder := getEncoder()

	if fileOutput == nil {
		fileOutput = zapcore.AddSync(getFileLogger(logConfig.Filename))
	}

	logLevelStr := logConfig.Level
	logLevel := getLogLevel(logLevelStr)

	var cores []zapcore.Core
	cores = append(cores, newCoreWithLevel(encoder, fileOutput, logLevel))

	if logLevel == zap.DebugLevel {
		cores = append(cores, newCoreWithLevel(encoder, zapcore.AddSync(os.Stdout), logLevel))
	}

	if errorFileOutput == nil {
		errorFileOutput = zapcore.AddSync(getFileLogger(logConfig.ErrorFilename))
	}
	cores = append(cores, newCoreWithLevel(encoder, errorFileOutput, zap.ErrorLevel))

	return zapcore.NewTee(cores...)
}

func initLogger() {
	core := getZapCore()
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	// 重定向标准库的日志输出到zap
	zap.RedirectStdLog(logger)

	// 替换全局logger
	zap.ReplaceGlobals(logger)
}

// Close 关闭zap记录器并释放资源
func Close() error {
	if fileOutput != nil {
		if err := fileOutput.Sync(); err != nil {
			return err
		}
	}
	if errorFileOutput != nil {
		if err := errorFileOutput.Sync(); err != nil {
			return err
		}
	}
	return nil
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	logger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func InfoWithTraceID(c *gin.Context, format string, a ...any) {
	traceID, _ := c.Get("TraceID")
	message := fmt.Sprintf(format, a...)
	logger.Info(message, zap.String("TraceID", traceID.(string)))
}

func WarnWithTraceID(c *gin.Context, format string, a ...any) {
	traceID, _ := c.Get("TraceID")
	message := fmt.Sprintf(format, a...)
	logger.Warn(message, zap.String("TraceID", traceID.(string)))
}

func ErrorWithTraceID(c *gin.Context, format string, a ...any) {
	traceID, _ := c.Get("TraceID")
	message := fmt.Sprintf(format, a...)
	logger.Error(message, zap.String("TraceID", traceID.(string)))
}

func PanicWithTraceID(c *gin.Context, format string, a ...any) {
	traceID, _ := c.Get("TraceID")
	message := fmt.Sprintf(format, a...)
	logger.Error(message, zap.String("TraceID", traceID.(string)))
}

func FatalWithTraceID(c *gin.Context, format string, a ...any) {
	traceID, _ := c.Get("TraceID")
	message := fmt.Sprintf(format, a...)
	logger.Error(message, zap.String("TraceID", traceID.(string)))
}
