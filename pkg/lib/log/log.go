package log

import (
	"github.com/go-logr/logr"
	clog "sigs.k8s.io/controller-runtime/pkg/log"
)

// Interface is a wrapper for logr.Logger
type Interface = logr.Logger

var (
	// Logger is the base logger used application.  It delegates
	// to another logr.Logger.  You *must* call SetLogger to
	// get any actual logging.
	Logger = clog.Log

	// SetLogger sets a concrete logging implementation for all deferred Loggers.
	SetLogger = clog.SetLogger
)
