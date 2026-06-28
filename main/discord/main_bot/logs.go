package mainbot

import (
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"
)

type logLine struct {
	at      time.Time
	level   string
	message string
}

type logBuffer struct {
	mu    sync.Mutex
	lines []logLine
}

func (b *logBuffer) add(level, message string) {
	b.mu.Lock()
	b.lines = append(b.lines, logLine{
		at:      time.Now(),
		level:   level,
		message: strings.TrimSpace(message),
	})
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
	out.WriteString("  MAIN BOT — SESSION LOG\n")
	out.WriteString("═══════════════════════════════════════\n\n")
	for _, line := range b.lines {
		out.WriteString(fmt.Sprintf("[%s] %-5s %s\n",
			line.at.Format("15:04:05"), line.level, line.message))
	}
	out.WriteString(fmt.Sprintf("\n───────────────────────────────────────\n  %d entries\n", len(b.lines)))
	return out.String()
}

func (b *logBuffer) reset() {
	b.mu.Lock()
	b.lines = nil
	b.mu.Unlock()
}

var (
	sessionLogs logBuffer
	logActive   bool
	logMu       sync.Mutex
)

func startLogCapture() {
	logMu.Lock()
	defer logMu.Unlock()
	sessionLogs.reset()
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

func botLog(format string, args ...interface{}) {
	sessionLogs.add("INFO", fmt.Sprintf(format, args...))
}

func botLogError(format string, args ...interface{}) {
	sessionLogs.add("ERROR", fmt.Sprintf(format, args...))
}
