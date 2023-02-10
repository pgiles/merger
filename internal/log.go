package internal

import (
	"fmt"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"strings"
	"time"
)

// Is a global var thread safe?  Apparently, it's the job of the Handler (via Enabled(Level)?) to manage concurrency.
var lvl = slog.LevelVar{}

func init() {
	initLogger()
}

// initLogger sets logging configuration options: level, format (JSON), destinations
// (stdout, or both stdout and file) using environment variables
// LOG_LEVEL, LOG_FORMAT, LOG_FILE, respectively.
//
// The presence of
// environment variable LOG_SOURCE adds a ("source", "file:line") attribute to
// the output indicating the source code position of the log statement.
func initLogger() {
	var timeFormat = time.RFC3339
	_, logSource := os.LookupEnv("LOG_SOURCE")
	var opts = slog.HandlerOptions{
		AddSource:   logSource,
		Level:       &lvl,
		ReplaceAttr: replaceAttr(timeFormat),
	}

	var writers = []io.Writer{os.Stdout}
	if f := writeToLogFile(); f != nil {
		writers = append(writers, f)
	}

	var h slog.Handler = opts.NewTextHandler(io.MultiWriter(writers...))
	if strings.ToLower(os.Getenv("LOG_FORMAT")) == "json" {
		h = opts.NewJSONHandler(io.MultiWriter(writers...))
	}

	setDefaultLogLevel()
	slog.SetDefault(slog.New(h))
}

func replaceAttr(timeFormat string) func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			a.Value = slog.StringValue(time.Now().Format(timeFormat))
		}
		return a
	}
}

func setDefaultLogLevel() {
	// set the minimum record level that will be logged
	envLogLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
	envLogLevelAlt := strings.ToLower(os.Getenv("LEVEL"))

	if envLogLevel == "debug" || envLogLevelAlt == "debug" {
		lvl.Set(slog.LevelDebug)
	} else if envLogLevel == "info" || envLogLevelAlt == "info" {
		lvl.Set(slog.LevelInfo)
	} else {
		lvl.Set(slog.LevelError)
	}
}

func writeToLogFile() *os.File {
	if filename, b := os.LookupEnv("LOG_FILE"); b {
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("Error]: %s", err))
		}
		slog.Info(fmt.Sprintf("Logging to %v", filename))
		return f
	}

	return nil
}

// LogPanic logs to ERROR (would prefer to log as FATAL, but I'm not going to
// create custom levels to do it) and panics. It is a necessary convenience
// method that is here in absence of a log.Panic in golang.org/x/exp/slog. It's
// not ideal.
func LogPanic(msg string, err error, args ...any) {
	//os.Setenv("LOG_SOURCE", "1")
	args = append(args, slog.Any("level", "FATAL"))
	slog.Error(msg, err, args...)
	panic(fmt.Sprintf("\n%v", err))
}
