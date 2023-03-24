package logger

import (
	"fmt"
	"gateway/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"sync"
	"time"
)

type Level = zapcore.Level

// 定义 Level 类型，表示日志级别
const (
	InfoLevel   Level = zap.InfoLevel   // 0, default level
	WarnLevel   Level = zap.WarnLevel   // 1
	ErrorLevel  Level = zap.ErrorLevel  // 2
	DPanicLevel Level = zap.DPanicLevel // 3, used in development logger
	// PanicLevel logs a message, then panics
	PanicLevel Level = zap.PanicLevel // 4
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel Level = zap.FatalLevel // 5
	DebugLevel Level = zap.DebugLevel // -1
)

type Field = zap.Field

func (l *Logger) Debug(msg string, fields ...Field) {
	l.l.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.l.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.l.Error(msg, fields...)
}

func (l *Logger) DPanic(msg string, fields ...Field) {
	l.l.DPanic(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...Field) {
	l.l.Panic(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.l.Fatal(msg, fields...)
}

// Debugf logs a message with sprintf syntax at DebugLevel
// Debugf 使用 sprintf 语法在 DebugLevel 记录一条消息
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.l.Debug(fmt.Sprintf(template, args...))
}

// Infof logs a message with sprintf syntax at InfoLevel
// Infof 使用 sprintf 语法在 InfoLevel 记录一条消息
func (l *Logger) Infof(template string, args ...interface{}) {
	l.l.Info(fmt.Sprintf(template, args...))
}

// Warnf logs a message with sprintf syntax at WarnLevel
// Warnf 使用 sprintf 语法在 WarnLevel 记录一条消息
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.l.Warn(fmt.Sprintf(template, args...))
}

// Errorf logs a message with sprintf syntax at ErrorLevel
// Errorf 使用 sprintf 语法在 ErrorLevel 记录一条消息
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.l.Error(fmt.Sprintf(template, args...))
}

// DPanicf logs a message with sprintf syntax at DPanicLevel
// DPanicf 使用 sprintf 语法在 DPanicLevel 记录一条消息
func (l *Logger) DPanicf(template string, args ...interface{}) {
	l.l.DPanic(fmt.Sprintf(template, args...))
}

// Panicf logs a message with sprintf syntax at PanicLevel
// Panicf 使用 sprintf 语法在 PanicLevel 记录一条消息
func (l *Logger) Panicf(template string, args ...interface{}) {
	l.l.Panic(fmt.Sprintf(template, args...))
}

// Fatalf logs a message with sprintf syntax at FatalLevel
// Fatalf 使用 sprintf 语法在 FatalLevel 记录一条消息
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.l.Fatal(fmt.Sprintf(template, args...))
}

// function variables for all field types
// in github.com/uber-go/zap/field.go
// 定义一些函数，用于包装不同类型的日志记录器
var (
	Skip        = zap.Skip
	Binary      = zap.Binary
	Bool        = zap.Bool
	Boolp       = zap.Boolp
	ByteString  = zap.ByteString
	Complex128  = zap.Complex128
	Complex128p = zap.Complex128p
	Complex64   = zap.Complex64
	Complex64p  = zap.Complex64p
	Float64     = zap.Float64
	Float64p    = zap.Float64p
	Float32     = zap.Float32
	Float32p    = zap.Float32p
	Int         = zap.Int
	Intp        = zap.Intp
	Int64       = zap.Int64
	Int64p      = zap.Int64p
	Int32       = zap.Int32
	Int32p      = zap.Int32p
	Int16       = zap.Int16
	Int16p      = zap.Int16p
	Int8        = zap.Int8
	Int8p       = zap.Int8p
	String      = zap.String
	Stringp     = zap.Stringp
	Uint        = zap.Uint
	Uintp       = zap.Uintp
	Uint64      = zap.Uint64
	Uint64p     = zap.Uint64p
	Uint32      = zap.Uint32
	Uint32p     = zap.Uint32p
	Uint16      = zap.Uint16
	Uint16p     = zap.Uint16p
	Uint8       = zap.Uint8
	Uint8p      = zap.Uint8p
	Uintptr     = zap.Uintptr
	Uintptrp    = zap.Uintptrp
	Reflect     = zap.Reflect
	Namespace   = zap.Namespace
	Stringer    = zap.Stringer
	Time        = zap.Time
	Timep       = zap.Timep
	Stack       = zap.Stack
	StackSkip   = zap.StackSkip
	Duration    = zap.Duration
	Durationp   = zap.Durationp
	Any         = zap.Any

	Info    = std.Info    // Info 级别日志
	Infof   = std.Infof   // 格式化输出 Info 级别日志
	Warn    = std.Warn    // Warn 级别日志
	Warnf   = std.Warnf   // 格式化输出 Warn 级别日志
	Error   = std.Error   // Error 级别日志
	Errorf  = std.Errorf  // 格式化输出 Error 级别日志
	DPanic  = std.DPanic  // DPanic 级别日志
	DPanicf = std.DPanicf // 格式化输出 DPanic 级别日志
	Panic   = std.Panic   // Panic 级别日志
	Panicf  = std.Panicf  // 格式化输出 Panic 级别日志
	Fatal   = std.Fatal   // Fatal 级别日志
	Fatalf  = std.Fatalf  // 格式化输出 Fatal 级别日志
	Debug   = std.Debug   // Debug 级别日志
	Debugf  = std.Debugf  // 格式化输出 Debug 级别日志
)

// not safe for concurrent use
func ResetDefault(l *Logger) {
	std = l
	Info = std.Info
	Infof = std.Infof
	Warn = std.Warn
	Warnf = std.Warnf
	Error = std.Error
	Errorf = std.Errorf
	DPanic = std.DPanic
	DPanicf = std.DPanicf
	Panic = std.Panic
	Panicf = std.Panicf
	Fatal = std.Fatal
	Fatalf = std.Fatalf
	Debug = std.Debug
	Debugf = std.Debugf
}

type Logger struct {
	l     *zap.Logger // zap ensure that zap.Logger is safe for concurrent use
	level Level
}

// 定义全局变量，存储默认日志记录器
var std = New(os.Stderr, InfoLevel, WithCaller(true))

// Default 返回默认 Logger 实例
func Default() *Logger {
	return std
}

// Option 是 zap.Logger 的配置选项，定义在 zap 包中
type Option = zap.Option

// WithCaller 可添加是否输出日志调用堆栈信息的选项，定义在 zap 包中
var (
	WithCaller    = zap.WithCaller    // 添加日志调用堆栈信息选项
	AddStacktrace = zap.AddStacktrace // 添加堆栈信息选项
)

// RotateOptions 保存日志轮换配置
type RotateOptions struct {
	MaxSize    int  // 单个日志文件最大尺寸，单位为 MB
	MaxAge     int  // 日志文件保存时间，单位为天
	MaxBackups int  // 最多保留的日志文件数
	Compress   bool // 是否压缩历史日志文件
}

// LevelEnablerFunc 定义一个能够判断日志级别是否被允许的函数类型
type LevelEnablerFunc func(lvl Level) bool

// TeeOption 定义一个 TeeLogger 需要的选项
type TeeOption struct {
	Filename string
	Ropt     RotateOptions
	Left     LevelEnablerFunc
}

// NewTeeWithRotate 基于 zap 实现的支持日志轮换的输出日志到多个文件的 Logger
// tops：输出到文件的选项列表
// opts：Logger 的配置选项
func NewTeeWithRotate(tops []TeeOption, opts ...Option) *Logger {
	var cores []zapcore.Core
	cfg := zap.NewProductionConfig()
	// 定义时间格式
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
	}

	for _, top := range tops {
		top := top

		lv := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return top.Left(Level(lvl))
		})

		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   top.Filename,
			MaxSize:    top.Ropt.MaxSize,
			MaxBackups: top.Ropt.MaxBackups,
			MaxAge:     top.Ropt.MaxAge,
			Compress:   top.Ropt.Compress,
		})

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.AddSync(w),
			lv,
		)
		cores = append(cores, core)
	}

	logger := &Logger{
		l: zap.New(zapcore.NewTee(cores...), opts...),
	}
	return logger
}

// New create a new logger (not support logger rotating).
// New 创建一个新的记录器（不支持文件轮换）
func New(writer io.Writer, level Level, opts ...Option) *Logger {
	if writer == nil {
		panic("the writer is nil")
	}
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
	}
	// 使用 zapcore.NewCore 创建一个新的 zapcore.Core，将输出到的写入器和日志级别等信息传递进去。
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.AddSync(writer),
		zapcore.Level(level),
	)
	logger := &Logger{
		l:     zap.New(core, opts...),
		level: level,
	}
	return logger
}

// Sync flushes any buffered log entries. Applications should take care to call Sync before exiting.
// 同步刷新所有缓存的日志条目。应用程序在退出之前应该注意调用Sync。
func (l *Logger) Sync() error {
	return l.l.Sync()
}

// Sync flushes any buffered log entries for the global logger. Applications should take care to call Sync before exiting.
// 同步刷新全局日志记录器的所有缓存的日志条目。应用程序在退出之前应该注意调用Sync。
func Sync() error {
	if std != nil {
		return std.Sync()
	}
	return nil
}

// loggerOnce用于确保在并发环境下只执行一次代码块，从而实现初始化操作的线程安全
var loggerOnce sync.Once

// InitLogger initializes the default logger with the loaded logger configuration
// 使用已加载的日志配置初始化默认日志记录器
func init() {
	loggerOnce.Do(func() {
		logConfig := configs.GetLogConfig()

		ropt := RotateOptions{
			MaxSize:    logConfig.MaxSize,
			MaxAge:     logConfig.MaxAge,
			MaxBackups: logConfig.MaxBackups,
			Compress:   false,
		}

		// Initialize the logger with logger rotation
		// 使用日志轮换初始化日志记录器
		level := zapLevelFromString(logConfig.Level)
		teeOption := TeeOption{
			Filename: logConfig.Filename,
			Ropt:     ropt,
			Left:     func(lvl Level) bool { return lvl >= level },
		}

		logger := NewTeeWithRotate([]TeeOption{teeOption}, WithCaller(true))
		ResetDefault(logger)
	})
}

// Convert logger level string to zapcore.Level
// 将日志级别的字符串转换为zapcore.Level类型
func zapLevelFromString(level string) Level {
	switch level {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "dpanic":
		return DPanicLevel
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}
