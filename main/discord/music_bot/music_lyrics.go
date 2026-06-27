package musicbot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type geniusSearchResponse struct {
	Response struct {
		Hits []struct {
			Result struct {
				URL string `json:"url"`
			} `json:"result"`
		} `json:"hits"`
	} `json:"response"`
}

func fetchLyrics(query, token string) ([]string, string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.genius.com/search", nil)
	if err != nil {
		return nil, "", err
	}
	q := req.URL.Query()
	q.Set("q", query)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("genius api error: %d", resp.StatusCode)
	}

	var payload geniusSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, "", err
	}
	if len(payload.Response.Hits) == 0 {
		return nil, "", fmt.Errorf("lyrics not found on genius")
	}

	songURL := payload.Response.Hits[0].Result.URL
	page, err := http.Get(songURL)
	if err != nil {
		return nil, songURL, err
	}
	defer page.Body.Close()

	doc, err := goquery.NewDocumentFromReader(page.Body)
	if err != nil {
		return nil, songURL, err
	}

	var parts []string
	doc.Find("div").Each(func(_ int, sel *goquery.Selection) {
		class, _ := sel.Attr("class")
		if strings.Contains(class, "Lyrics__Container") {
			text := strings.TrimSpace(sel.Text())
			if text != "" {
				parts = append(parts, text)
			}
		}
	})
	if len(parts) == 0 {
		return nil, songURL, fmt.Errorf("lyrics page empty")
	}

	return chunkText(strings.Join(parts, "\n"), 1990), songURL, nil
}

func chunkText(text string, size int) []string {
	if len(text) <= size {
		return []string{text}
	}
	var chunks []string
	lines := strings.Split(text, "\n")
	var current strings.Builder
	for _, line := range lines {
		if current.Len()+len(line)+1 > size {
			chunks = append(chunks, current.String())
			current.Reset()
		}
		if current.Len() > 0 {
			current.WriteByte('\n')
		}
		current.WriteString(line)
	}
	if current.Len() > 0 {
		chunks = append(chunks, current.String())
	}
	return chunks
}
