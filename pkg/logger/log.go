package logger

import (
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
	l.Debug(msg, fields...)
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

// function variables for all field types
// in github.com/uber-go/zap/field.go

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

	Info   = std.Info
	Warn   = std.Warn
	Error  = std.Error
	DPanic = std.DPanic
	Panic  = std.Panic
	Fatal  = std.Fatal
	Debug  = std.Debug
)

// not safe for concurrent use
func ResetDefault(l *Logger) {
	std = l
	Info = std.Info
	Warn = std.Warn
	Error = std.Error
	DPanic = std.DPanic
	Panic = std.Panic
	Fatal = std.Fatal
	Debug = std.Debug
}

type Logger struct {
	l     *zap.Logger // zap ensure that zap.Logger is safe for concurrent use
	level Level
}

var std = New(os.Stderr, InfoLevel, WithCaller(true))

func Default() *Logger {
	return std
}

type Option = zap.Option

var (
	WithCaller    = zap.WithCaller
	AddStacktrace = zap.AddStacktrace
)

type RotateOptions struct {
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
}

type LevelEnablerFunc func(lvl Level) bool

type TeeOption struct {
	Filename string
	Ropt     RotateOptions
	Left     LevelEnablerFunc
}

func NewTeeWithRotate(tops []TeeOption, opts ...Option) *Logger {
	var cores []zapcore.Core
	cfg := zap.NewProductionConfig()
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
func New(writer io.Writer, level Level, opts ...Option) *Logger {
	if writer == nil {
		panic("the writer is nil")
	}
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
	}

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

func (l *Logger) Sync() error {
	return l.l.Sync()
}

func Sync() error {
	if std != nil {
		return std.Sync()
	}
	return nil
}

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
