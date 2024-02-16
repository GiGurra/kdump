package util_log

import (
	"log/slog"
	"os"
)

type LogCfg struct {
	MinLevel        slog.Level
	IncludeSrcLines bool
	LoggerName      string
}

func (l LogCfg) Default() LogCfg {
	l.MinLevel = slog.LevelInfo
	l.IncludeSrcLines = false
	return l
}

func ConfigureGcpCompatibleJsonDefaultSlog(
	cfg LogCfg,
) {
	if cfg.LoggerName == "" {
		cfg.LoggerName = "root"
	}
	slog.SetDefault(NewGcpCompatibleJsonSlog(cfg))
}

func NewGcpCompatibleJsonSlog(
	cfg LogCfg,
) *slog.Logger {
	result := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: cfg.IncludeSrcLines,
		Level:     cfg.MinLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) == 0 && a.Key == "level" {
				return slog.Attr{Key: "severity", Value: a.Value}
			}
			if len(groups) == 0 && a.Key == "msg" {
				return slog.Attr{Key: "message", Value: a.Value}
			}
			return a
		},
	}))
	if cfg.LoggerName != "" {
		result = result.With(slog.String("logger", cfg.LoggerName))
	}
	return result
}
