package musicbot

import (
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"
)

type logLevel int

const (
	logInfo logLevel = iota
	logWarn
	logError
)

type logLine struct {
	at      time.Time
	level   logLevel
	message string
}

type logBuffer struct {
	mu    sync.Mutex
	lines []logLine
}

func (b *logBuffer) add(level logLevel, message string) {
	b.mu.Lock()
	b.lines = append(b.lines, logLine{at: time.Now(), level: level, message: strings.TrimSpace(message)})
	b.mu.Unlock()
}

func (b *logBuffer) formatted() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.lines) == 0 {
		return ""
	}

	var out strings.Builder
	out.WriteString("═══════════════════════════════════════\n")
	out.WriteString("  MUSIC BOT — SESSION LOG\n")
	out.WriteString("═══════════════════════════════════════\n\n")

	for _, line := range b.lines {
		out.WriteString(fmt.Sprintf("[%s] %s %s\n",
			line.at.Format("15:04:05"),
			levelTag(line.level),
			line.message,
		))
	}

	out.WriteString("\n───────────────────────────────────────\n")
	out.WriteString(fmt.Sprintf("  %d entries\n", len(b.lines)))
	return out.String()
}

func (b *logBuffer) Reset() {
	b.mu.Lock()
	b.lines = nil
	b.mu.Unlock()
}

func levelTag(level logLevel) string {
	switch level {
	case logWarn:
		return "WARN "
	case logError:
		return "ERROR"
	default:
		return "INFO "
	}
}

var (
	sessionLogs logBuffer
	logActive   bool
	logMu       sync.Mutex
)

func startLogCapture() {
	logMu.Lock()
	defer logMu.Unlock()
	sessionLogs.Reset()
	logActive = true
	log.SetOutput(io.Discard)
}

func stopLogCapture() string {
	logMu.Lock()
	defer logMu.Unlock()
	if logActive {
		log.SetOutput(io.Discard)
		logActive = false
	}
	return sessionLogs.formatted()
}

func SessionLogs() string {
	return sessionLogs.formatted()
}

func writeLog(level logLevel, format string, args ...interface{}) {
	sessionLogs.add(level, fmt.Sprintf(format, args...))
}

func botLogInfo(format string, args ...interface{})  { writeLog(logInfo, format, args...) }
func botLogWarn(format string, args ...interface{})  { writeLog(logWarn, format, args...) }
func botLogError(format string, args ...interface{}) { writeLog(logError, format, args...) }

func botLog(format string, args ...interface{}) { botLogInfo(format, args...) }
