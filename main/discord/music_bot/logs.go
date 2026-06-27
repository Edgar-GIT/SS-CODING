package musicbot

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type logBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *logBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	b.buf.Write(p)
	b.mu.Unlock()
	return len(p), nil
}

func (b *logBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func (b *logBuffer) Reset() {
	b.mu.Lock()
	b.buf.Reset()
	b.mu.Unlock()
}

var (
	sessionLogs logBuffer
	logMirror   io.Writer = os.Stderr
	logActive   bool
	logMu       sync.Mutex
)

func startLogCapture() {
	logMu.Lock()
	defer logMu.Unlock()
	sessionLogs.Reset()
	logActive = true
	log.SetOutput(io.MultiWriter(&sessionLogs, logMirror))
}

func stopLogCapture() string {
	logMu.Lock()
	defer logMu.Unlock()
	if logActive {
		log.SetOutput(logMirror)
		logActive = false
	}
	return sessionLogs.String()
}

func SessionLogs() string {
	return sessionLogs.String()
}

func botLog(format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}
	sessionLogs.Write([]byte(line))
	logMirror.Write([]byte(line))
}
