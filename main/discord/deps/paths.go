package deps

import (
	"os"
	"path/filepath"
	"sync"
)

// Runtime layout (gitignored, created by InstallAll):
//
//	temp_music/
//	├── downloads/   — temp audio files
//	├── playlists/   — user playlist JSON
//	└── tools/       — yt-dlp + ffmpeg (auto-downloaded)

var (
	once         sync.Once
	musicBotDir  string
	tempMusicDir string
	downloadsDir string
	playlistsDir string
	toolsDir     string
)

func initPaths() {
	musicBotDir = resolveMusicBotDir()
	tempMusicDir = filepath.Join(musicBotDir, "temp_music")
	downloadsDir = filepath.Join(tempMusicDir, "downloads")
	playlistsDir = filepath.Join(tempMusicDir, "playlists")
	toolsDir = filepath.Join(tempMusicDir, "tools")
}

func resolveMusicBotDir() string {
	if dir, err := os.Getwd(); err == nil {
		for {
			candidate := filepath.Join(dir, "discord", "music_bot")
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				break
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	if exe, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(exe), "music_bot")
	}
	return filepath.Join("discord", "music_bot")
}

func EnsureDirs() error {
	once.Do(initPaths)
	for _, dir := range []string{tempMusicDir, downloadsDir, playlistsDir, toolsDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return nil
}

func DownloadsDir() string {
	once.Do(initPaths)
	return downloadsDir
}

func PlaylistPath(userID string) string {
	once.Do(initPaths)
	return filepath.Join(playlistsDir, userID+".json")
}

func toolPath(name string) string {
	once.Do(initPaths)
	return filepath.Join(toolsDir, name)
}
