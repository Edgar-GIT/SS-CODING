package musicbot

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"ss-coding/discord/deps"
)

type Track struct {
	Title     string `json:"title"`
	Query     string `json:"query,omitempty"` // original playlist/search text
	Uploader  string `json:"uploader,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Duration  int    `json:"duration"`
	URL       string `json:"url,omitempty"`
	FilePath  string `json:"file_path,omitempty"`
}

func searchYouTube(query string) (*Track, error) {
	return downloadYouTube(query)
}

func searchSoundCloud(query string) (*Track, error) {
	target := query
	if !strings.HasPrefix(query, "http://") && !strings.HasPrefix(query, "https://") {
		target = "scsearch:" + query
	}
	return extractStream(context.Background(), target, "scsearch")
}

func extractStream(ctx context.Context, target, mode string) (*Track, error) {
	ytdlp, err := deps.YTDlpPath()
	if err != nil {
		return nil, err
	}
	ffmpeg, err := deps.FFmpegPath()
	if err != nil {
		return nil, err
	}

	args := []string{
		"--quiet", "--no-warnings", "--no-playlist",
		"--ffmpeg-location", filepath.Dir(ffmpeg),
		"-f", "bestaudio/best",
		"--dump-single-json",
	}
	if mode == "ytsearch" || mode == "scsearch" {
		args = append(args, target)
	} else {
		args = append(args, target)
	}

	out, err := exec.CommandContext(ctx, ytdlp, args...).Output()
	if err != nil {
		return nil, err
	}

	var payload map[string]any
	if err := json.Unmarshal(out, &payload); err != nil {
		return nil, err
	}
	if entries, ok := payload["entries"].([]any); ok && len(entries) > 0 {
		if entry, ok := entries[0].(map[string]any); ok {
			payload = entry
		}
	}

	streamURL := pickAudioURL(payload)
	if streamURL == "" {
		return nil, fmt.Errorf("no stream url")
	}

	return mapTrack(payload, streamURL, ""), nil
}

func downloadYouTube(query string) (*Track, error) {
	ytdlp, err := deps.YTDlpPath()
	if err != nil {
		return nil, err
	}
	ffmpeg, err := deps.FFmpegPath()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(deps.DownloadsDir(), fmt.Sprintf("temp_%s.mp3", uuid.NewString()))
	args := []string{
		"--quiet", "--no-warnings",
		"--ffmpeg-location", filepath.Dir(ffmpeg),
		"-f", "bestaudio/best",
		"-x", "--audio-format", "mp3",
		"--audio-quality", "192K",
		"-o", filePath,
		"ytsearch:" + query,
	}
	if err := exec.Command(ytdlp, args...).Run(); err != nil {
		return nil, err
	}
	if _, err := os.Stat(filePath); err != nil {
		base := strings.TrimSuffix(filePath, ".mp3")
		matches, _ := filepath.Glob(base + ".*")
		if len(matches) == 0 {
			return nil, fmt.Errorf("download failed")
		}
		filePath = matches[0]
	}

	meta, err := extractStream(context.Background(), "ytsearch:"+query, "ytsearch")
	if err != nil {
		meta = &Track{Title: query}
	}
	meta.FilePath = filePath
	meta.URL = ""
	return meta, nil
}

// resolveQueuedTrack resolves a title-only queue entry via streaming (fast, no full download).
func resolveQueuedTrack(ctx context.Context, query string) (*Track, error) {
	track, err := extractStream(ctx, "ytsearch:"+query, "ytsearch")
	if err != nil {
		return nil, err
	}
	if track.Query == "" {
		track.Query = query
	}
	return track, nil
}

func pickAudioURL(payload map[string]any) string {
	if formats, ok := payload["formats"].([]any); ok {
		var best string
		var bestRate float64
		for _, item := range formats {
			format, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if acodec, _ := format["acodec"].(string); acodec == "none" {
				continue
			}
			url, _ := format["url"].(string)
			if url == "" {
				continue
			}
			rate, _ := format["tbr"].(float64)
			if rate == 0 {
				rate, _ = format["abr"].(float64)
			}
			if rate >= bestRate {
				bestRate = rate
				best = url
			}
		}
		if best != "" {
			return best
		}
	}
	if url, ok := payload["url"].(string); ok {
		return url
	}
	return ""
}

func mapTrack(payload map[string]any, streamURL, filePath string) *Track {
	title, _ := payload["title"].(string)
	uploader, _ := payload["uploader"].(string)
	thumb, _ := payload["thumbnail"].(string)
	duration := 0
	switch value := payload["duration"].(type) {
	case float64:
		duration = int(value)
	case int:
		duration = value
	}
	return &Track{
		Title:     title,
		Uploader:  uploader,
		Thumbnail: thumb,
		Duration:  duration,
		URL:       streamURL,
		FilePath:  filePath,
	}
}

func formatDuration(seconds int) string {
	if seconds <= 0 {
		return "Unknown"
	}
	return fmt.Sprintf("%d:%02d", seconds/60, seconds%60)
}
