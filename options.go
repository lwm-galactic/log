package log

import "github.com/spf13/pflag"

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

// AddFlags adds flags for log to the specified FlagSet object.
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

func (l *logger) SetOptions(opts ...Option) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, opt := range opts {
		opt(l.opt)
	}
}
