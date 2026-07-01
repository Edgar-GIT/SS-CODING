package transporter

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RequestLogEntry struct {
	Time      time.Time
	IP        string
	Method    string
	Path      string
	UserAgent string
}

type IPLogSummary struct {
	IP        string
	Hits      int
	FirstSeen time.Time
	LastSeen  time.Time
	LastPath  string
	UserAgent string
}

var (
	logMu          sync.Mutex
	logEntries     []RequestLogEntry
	logProxyServer *http.Server
	logProxyPort   int
)

func StartLoggingProxy(targetPort int) (int, error) {
	logMu.Lock()
	if logProxyServer != nil {
		port := logProxyPort
		logMu.Unlock()
		return port, nil
	}
	logMu.Unlock()

	target, err := url.Parse("http://localhost:" + strconv.Itoa(targetPort))
	if err != nil {
		return 0, err
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host
	}

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recordRequest(r)
			proxy.ServeHTTP(w, r)
		}),
	}

	port := listener.Addr().(*net.TCPAddr).Port

	logMu.Lock()
	logProxyServer = server
	logProxyPort = port
	logMu.Unlock()

	go func() {
		_ = server.Serve(listener)
	}()

	return port, nil
}

func StopLoggingProxy() {
	logMu.Lock()
	server := logProxyServer
	logProxyServer = nil
	logProxyPort = 0
	logMu.Unlock()

	if server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}
}

func LoggingProxyRunning() bool {
	logMu.Lock()
	defer logMu.Unlock()
	return logProxyServer != nil
}

func RequestLogs() []RequestLogEntry {
	logMu.Lock()
	defer logMu.Unlock()

	entries := make([]RequestLogEntry, len(logEntries))
	copy(entries, logEntries)
	return entries
}

func IPLogSummaries() []IPLogSummary {
	logMu.Lock()
	defer logMu.Unlock()

	byIP := make(map[string]*IPLogSummary)
	for _, entry := range logEntries {
		summary := byIP[entry.IP]
		if summary == nil {
			summary = &IPLogSummary{
				IP:        entry.IP,
				FirstSeen: entry.Time,
			}
			byIP[entry.IP] = summary
		}
		summary.Hits++
		summary.LastSeen = entry.Time
		summary.LastPath = entry.Path
		summary.UserAgent = entry.UserAgent
	}

	summaries := make([]IPLogSummary, 0, len(byIP))
	for _, summary := range byIP {
		summaries = append(summaries, *summary)
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].LastSeen.After(summaries[j].LastSeen)
	})

	return summaries
}

func ClearRequestLogs() {
	logMu.Lock()
	logEntries = nil
	logMu.Unlock()
}

func recordRequest(r *http.Request) {
	entry := RequestLogEntry{
		Time:      time.Now(),
		IP:        requestIP(r),
		Method:    r.Method,
		Path:      r.URL.RequestURI(),
		UserAgent: r.UserAgent(),
	}

	logMu.Lock()
	logEntries = append(logEntries, entry)
	if len(logEntries) > 500 {
		logEntries = logEntries[len(logEntries)-500:]
	}
	logMu.Unlock()
}

func requestIP(r *http.Request) string {
	for _, header := range []string{"CF-Connecting-IP", "X-Real-IP", "X-Forwarded-For"} {
		value := strings.TrimSpace(r.Header.Get(header))
		if value == "" {
			continue
		}
		if header == "X-Forwarded-For" {
			value = strings.TrimSpace(strings.Split(value, ",")[0])
		}
		if value != "" {
			return value
		}
	}

	forwarded := strings.TrimSpace(r.Header.Get("Forwarded"))
	if forwarded != "" {
		for _, part := range strings.Split(forwarded, ";") {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(strings.ToLower(part), "for=") {
				return strings.Trim(strings.TrimPrefix(part, "for="), "\"[]")
			}
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
