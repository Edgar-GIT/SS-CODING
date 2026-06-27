package musicbot

import (
	"encoding/json"
	"fmt"
	"os"

	"ss-coding/discord/deps"
)

type UserPlaylist struct {
	Name  string   `json:"name"`
	Songs []string `json:"songs"`
}

func loadUserPlaylist(userID string) (*UserPlaylist, error) {
	path := deps.PlaylistPath(userID)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var playlist UserPlaylist
	if err := json.Unmarshal(data, &playlist); err != nil {
		return nil, err
	}
	return &playlist, nil
}

func saveUserPlaylist(userID string, playlist *UserPlaylist) error {
	if err := deps.EnsureDirs(); err != nil {
		return err
	}
	data, err := json.MarshalIndent(playlist, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(deps.PlaylistPath(userID), data, 0o644)
}

func createUserPlaylist(userID, name string) error {
	existing, err := loadUserPlaylist(userID)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("you already have a playlist")
	}
	return saveUserPlaylist(userID, &UserPlaylist{Name: name, Songs: []string{}})
}

func addSongToPlaylist(userID, song string) error {
	playlist, err := loadUserPlaylist(userID)
	if err != nil {
		return err
	}
	if playlist == nil {
		return fmt.Errorf("create a playlist first with %screateplaylist", commandPrefix)
	}
	playlist.Songs = append(playlist.Songs, song)
	return saveUserPlaylist(userID, playlist)
}
