package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

var LevelNameMapping = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	PanicLevel: "PANIC",
	FatalLevel: "FATAL",
}

const (
	flagLevel             = "log.level"
	flagDisableCaller     = "log.disable-caller"
	flagDisableStacktrace = "log.disable-stacktrace"
	flagFormat            = "log.format"
	flagEnableColor       = "log.enable-color"
	flagOutputPaths       = "log.output-paths"
	flagErrorOutputPaths  = "log.error-output-paths"
	flagDevelopment       = "log.development"
	flagName              = "log.name"

	consoleFormat = "console"
	jsonFormat    = "json"
)

// newDefaultOptions 返回一个带有默认配置的 options 实例
func newDefaultOptions() *options {
	return &options{
		Level:             InfoLevel,                      // 默认日志等级为 InfoLevel
		Format:            consoleFormat,                  // 默认输出格式为控制台格式
		OutputPaths:       []OutputPath{StdoutOutputPath}, // 默认输出到 stdout
		ErrorOutputPaths:  []string{"stderr"},             // 默认错误输出到 stderr
		DisableCaller:     false,                          // 启用调用者信息（文件:行号）
		DisableStacktrace: false,                          // 启用堆栈跟踪
		EnableColor:       true,                           // 启用颜色（适用于终端）
		Development:       false,                          // 非开发模式
		MaxSize:           100,                            // 默认每个日志文件最大 100MB
		MaxBackups:        5,                              // 最多保留 5 个备份文件
		MaxAge:            7,                              // 日志文件最多保留 7 天
		Compress:          false,                          // 不压缩旧文件
		Name:              "app-logger",                   // 默认 logger 名称
		FilePath:          "app.log",                      // 默认日志文件路径
	}
}

type OutputPath string

const (
	FileOutputPath   OutputPath = "file"
	StdoutOutputPath OutputPath = "stdout"
)

type options struct {
	OutputPaths       []OutputPath `mapstructure:"output-paths"`       // stdout / file
	FilePath          string       `mapstructure:"file"`               // 文件路径
	ErrorOutputPaths  []string     `mapstructure:"error-output-paths"` // zap 内部错误日志输出路径
	Level             Level        `mapstructure:"level"`              // 日志等级
	Format            string       `mapstructure:"format"`             // 输出格式化(JSON / TEXT)
	DisableCaller     bool         `mapstructure:"disable-caller"`     //
	DisableStacktrace bool         `mapstructure:"disable-stacktrace"` //

	EnableColor bool `mapstructure:"enable-color"` // 开启颜色
	Development bool `mapstructure:"development"`  // 是否为开发模式

	// 新增的日志轮转配置项
	MaxSize    int  `mapstructure:"max-size"`    // MB
	MaxBackups int  `mapstructure:"max-backups"` // 最多保留多少个备份文件
	MaxAge     int  `mapstructure:"max-age"`     // 文件最多保留几天
	Compress   bool `mapstructure:"compress"`    // 是否压缩旧文件

	Name string `mapstructure:"name"`
}

// Build 构建一个zap.logger 日志
func (o *options) build() (*zap.Logger, error) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(LevelNameMapping[o.Level])); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	// 构建 encoder config
	encodeLevel := zapcore.CapitalLevelEncoder
	if o.Format == consoleFormat && o.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     timeEncoder,
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	var encoder zapcore.Encoder
	if o.Format == jsonFormat {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 收集所有 Core
	var cores []zapcore.Core

	for _, path := range o.OutputPaths {
		if path == StdoutOutputPath {
			core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(zapLevel))
			cores = append(cores, core)
		} else if path == FileOutputPath {
			writer := &lumberjack.Logger{
				Filename:   o.FilePath,
				MaxSize:    o.MaxSize,
				MaxBackups: o.MaxBackups,
				MaxAge:     o.MaxAge,
				Compress:   o.Compress,
			}
			core := zapcore.NewCore(encoder, zapcore.AddSync(writer), zap.NewAtomicLevelAt(zapLevel))
			cores = append(cores, core)
		}
	}

	if len(cores) == 0 {
		return nil, fmt.Errorf("no valid output paths configured")
	}

	// 合并多个 Core
	core := zapcore.NewTee(cores...)

	// 构造 Logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.PanicLevel))

	if o.Name != "" {
		logger = logger.Named(o.Name)
	}

	return logger, nil
}

func (o *options) Validate() []error {
	var errs []error

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("not a valid log format: %q", o.Format))
	}

	return errs
}

// AddFlags adds flags for log to the specified FlagSet object.
/*
func (o *options) AddFlags(fs *pflag.FlagSet) {
	level := LevelNameMapping[o.Level]
	fs.StringVar(&level, flagLevel, level, "Minimum log output `LEVEL`.")
	fs.BoolVar(&o.DisableCaller, flagDisableCaller, o.DisableCaller, "Disable output of caller information in the log.")
	fs.BoolVar(&o.DisableStacktrace, flagDisableStacktrace,
		o.DisableStacktrace, "Disable the log to record a stack trace for all messages at or above panic level.")
	fs.StringVar(&o.Format, flagFormat, o.Format, "Log output `FORMAT`, support plain or json format.")
	fs.BoolVar(&o.EnableColor, flagEnableColor, o.EnableColor, "Enable output ansi colors in plain format logs.")
	fs.StringSliceVar(&o.OutputPaths, flagOutputPaths, o.OutputPaths, "Output paths of log.")
	fs.StringSliceVar(&o.ErrorOutputPaths, flagErrorOutputPaths, o.ErrorOutputPaths, "Error output paths of log.")
	fs.BoolVar(
		&o.Development,
		flagDevelopment,
		o.Development,
		"Development puts the logger in development mode, which changes "+
			"the behavior of DPanicLevel and takes stacktraces more liberally.",
	)
	fs.StringVar(&o.Name, flagName, o.Name, "The name of the logger.")
}
*/

type Option func(*options)

func WithLevel(level Level) Option {
	return func(o *options) { o.Level = level }
}

func WithOutputPaths(outputPaths ...OutputPath) Option {
	return func(o *options) { o.OutputPaths = outputPaths }
}

func WithErrorOutputPaths(errorOutputPaths []string) Option {
	return func(o *options) { o.ErrorOutputPaths = errorOutputPaths }
}

func WithDisableCaller(disableCaller bool) Option {
	return func(o *options) { o.DisableCaller = disableCaller }
}

func WithDisableStacktrace(disableStacktrace bool) Option {
	return func(o *options) { o.DisableStacktrace = disableStacktrace }
}

func WithEnableColor(enableColor bool) Option {
	return func(o *options) { o.EnableColor = enableColor }
}

func WithName(n string) Option {
	return func(o *options) { o.Name = n }
}
func WithDevelopment(development bool) Option {
	return func(o *options) { o.Development = development }
}
func WithMaxSize(maxSize int) Option {
	return func(o *options) { o.MaxSize = maxSize }
}
func WithMaxBackups(maxBackups int) Option {
	return func(o *options) { o.MaxBackups = maxBackups }
}
func WithMaxAge(maxAge int) Option {
	return func(o *options) { o.MaxAge = maxAge }
}
func WithCompress(compress bool) Option {
	return func(o *options) { o.Compress = compress }
}
func WithFilePath(filePath string) Option {
	return func(o *options) { o.FilePath = filePath }
}
func (l *Logger) SetOptions(opts ...Option) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, opt := range opts {
		opt(l.opt)
	}
}
