package log

import (
	"gateway/configs"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 全局变量
var (
	logger *zap.Logger

	logConfig *configs.LogConfig

	fileOutput      zapcore.WriteSyncer
	errorFileOutput zapcore.WriteSyncer
)

// Init 初始化日志
func Init() {
	initLogger()
	configs.RegisterReloadCallback(initLogger) // 注册回调
}

// initLogger 初始化zap记录器
// 原理:
// 1. 创建zapcore.Core
// 2. 创建zap.Logger
// 3. 重定向标准库的日志输出到zap
// 4. 替换全局logger
func initLogger() {
	// 重置
	logConfig = nil
	logger = nil
	fileOutput = nil
	errorFileOutput = nil

	// 初始化
	logConfig = configs.GetLogConfig()
	// 创建zapcore.Core
	core := getZapCore()

	// 创建zap.Logger
	// zap.AddCaller() 添加文件和行号
	// zap.AddStacktrace(zap.ErrorLevel) 添加错误堆栈信息
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	// 重定向标准库的日志输出到zap
	zap.RedirectStdLog(logger)

	// 替换全局logger
	zap.ReplaceGlobals(logger)
}

// getEncoder 获取编码器
func getEncoder() zapcore.Encoder {
	logConfig = configs.GetLogConfig() // 获取最新的配置
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
	switch logConfig.Format {
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig)
	default:
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

// getFileLogger 获取文件日志
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

// getLogLevel 获取日志级别
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

// newCoreWithLevel 创建带有日志级别的核心
func newCoreWithLevel(encoder zapcore.Encoder, output zapcore.WriteSyncer, level zapcore.Level) zapcore.Core {
	return zapcore.NewCore(encoder, output, zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	}))
}

// getZapCore 获取zap核心
func getZapCore() zapcore.Core {
	// 获取编码器
	encoder := getEncoder()

	// 获取文件日志
	if fileOutput == nil {
		fileOutput = zapcore.AddSync(getFileLogger(logConfig.Filename))
	}

	// 获取日志级别
	logLevel := getLogLevel(logConfig.Level)

	// 创建zapcore.Core
	var cores []zapcore.Core

	cores = append(cores, newCoreWithLevel(encoder, fileOutput, logLevel))

	// 如果日志级别为debug, 则将日志输出到控制台
	if logLevel == zap.DebugLevel {
		cores = append(cores, newCoreWithLevel(encoder, zapcore.AddSync(os.Stdout), logLevel))
	}

	// 如果日志级别为error, 则将日志输出到error文件
	if errorFileOutput == nil {
		errorFileOutput = zapcore.AddSync(getFileLogger(logConfig.ErrorFilename))
	}
	cores = append(cores, newCoreWithLevel(encoder, errorFileOutput, zap.ErrorLevel))

	return zapcore.NewTee(cores...)
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
