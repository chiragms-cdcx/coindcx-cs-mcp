package logger

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
)

func PrettyJSON(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	var v any
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return s
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return s
	}
	return string(b)
}

const (
	LevelError = iota
	LevelInfo
	LevelDebug
)

var (
	level      = LevelInfo
	levelOnce  sync.Once
	outputOnce sync.Once
)

func initLevel() {
	levelOnce.Do(func() {
		switch strings.ToLower(strings.TrimSpace(os.Getenv("LOG_LEVEL"))) {
		case "debug":
			level = LevelDebug
		case "info", "":
			level = LevelInfo
		case "error":
			level = LevelError
		default:
			level = LevelInfo
		}
	})
}

func initOutput() {
	outputOnce.Do(func() {
		if path := strings.TrimSpace(os.Getenv("LOG_FILE")); path != "" {
			f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				log.Printf("[ERROR] LOG_FILE=%q: %v; using stderr", path, err)
				return
			}
			log.SetOutput(f)
		}
	})
}

func Debug(format string, v ...any) {
	initLevel()
	initOutput()
	if level >= LevelDebug {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func Info(format string, v ...any) {
	initLevel()
	initOutput()
	if level >= LevelInfo {
		log.Printf("[INFO] "+format, v...)
	}
}

func Error(format string, v ...any) {
	initLevel()
	initOutput()
	if level >= LevelError {
		log.Printf("[ERROR] "+format, v...)
	}
}
