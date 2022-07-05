// Package zap contains helpers for setting up a new logr.Logger instance
// using the Zap logging framework.
package zap

import (
	"io"
	"os"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns a brand new Logger configured with Opts.
func New(opts ...Opts) logr.Logger {
	return zapr.NewLogger(NewRaw(opts...))
}

// Opts allows to manipulate Options
type Opts func(*Options)

// Options contains all possible settings
type Options struct {
	// If Development is true, a Zap development config will be used
	// (no sampling), otherwise a Zap production
	// config will be used (sampling).
	Development bool
	// The encoder to use, defaults to console when Development is true
	// and JSON otherwise
	Encoder zapcore.Encoder
	// The destination to write to, defaults to os.Stdout
	DestWritter io.Writer
	// The level to use, defaults to Debug when Development is true and
	// Info otherwise
	Level *zap.AtomicLevel
	// StacktraceLevel is the level at and above which stacktraces will
	// be recorded for all messages.
	StacktraceLevel *zap.AtomicLevel
	// Raw zap.Options to configure on the underlying zap logger
	ZapOpts []zap.Option
}

// addDefaults adds defaults to the Options
func (o *Options) addDefaults() {
	if o.DestWritter == nil {
		o.DestWritter = os.Stdout
	}

	if o.Development {
		if o.Encoder == nil {
			encCfg := zap.NewDevelopmentEncoderConfig()
			o.Encoder = zapcore.NewConsoleEncoder(encCfg)
		}
		if o.Level == nil {
			lvl := zap.NewAtomicLevelAt(zap.DebugLevel)
			o.Level = &lvl
		}
		if o.StacktraceLevel == nil {
			lvl := zap.NewAtomicLevelAt(zap.PanicLevel)
			o.StacktraceLevel = &lvl
		}
		o.ZapOpts = append(o.ZapOpts, zap.Development())

	} else {
		if o.Encoder == nil {
			encCfg := zap.NewProductionEncoderConfig()
			o.Encoder = zapcore.NewJSONEncoder(encCfg)
		}
		if o.Level == nil {
			lvl := zap.NewAtomicLevelAt(zap.InfoLevel)
			o.Level = &lvl
		}
		if o.StacktraceLevel == nil {
			lvl := zap.NewAtomicLevelAt(zap.PanicLevel)
			o.StacktraceLevel = &lvl
		}
		o.ZapOpts = append(o.ZapOpts,
			zap.WrapCore(func(core zapcore.Core) zapcore.Core {
				return zapcore.NewSampler(core, time.Second, 100, 100)
			}))
	}

	o.ZapOpts = append(o.ZapOpts, zap.AddStacktrace(o.StacktraceLevel))
}

// NewRaw returns a new zap.Logger configured with the passed Opts
// or their defaults.
func NewRaw(opts ...Opts) *zap.Logger {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}
	o.addDefaults()

	// this basically mimics New<type>Config, but with a custom sink
	sink := zapcore.AddSync(o.DestWritter)

	o.ZapOpts = append(o.ZapOpts, zap.AddCallerSkip(1), zap.ErrorOutput(sink))
	log := zap.New(zapcore.NewCore(o.Encoder, sink, *o.Level))
	log = log.WithOptions(o.ZapOpts...)
	return log
}
