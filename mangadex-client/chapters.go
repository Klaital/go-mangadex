package mangadex_client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

type ChapterFeedResponse struct {
	Result   string            `json:"result"`
	Response string            `json:"response"`
	Data     []ChapterFeedData `json:"data"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
	Total    int               `json:"total"`
}

type ChapterFeedData struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Title              string `json:"title"`
		Volume             string `json:"volume"`
		Chapter            string `json:"chapter"`
		Pages              int    `json:"pages"`
		TranslatedLanguage string `json:"translatedLanguage"`
		Uploader           string `json:"uploader"`
		ExternalUrl        string `json:"externalUrl"`
		Version            int    `json:"version"`
		CreatedAt          string `json:"createdAt"`
		UpdatedAt          string `json:"updatedAt"`
		PublishAt          string `json:"publishAt"`
		ReadableAt         string `json:"readableAt"`
	} `json:"attributes"`
	Relationships []struct {
		Id         string `json:"id"`
		Type       string `json:"type"`
		Related    string `json:"related"`
		Attributes struct {
		} `json:"attributes"`
	} `json:"relationships"`
}

// ListChapters returns data about the chapters available for a given manga
// Fetched data is cached locally for future requests.
func (c Client) ListChapters(id string, force bool) ([]ChapterFeedData, error) {
	cachePath := filepath.Join(c.CacheDir, id, "chapters.json")
	var chapterData []ChapterFeedData

	// Fetch from the cache
	if data, err := os.ReadFile(cachePath); err != nil {
		// handle cache miss
		slog.Debug("Cache miss for chapters", "id", id)
	} else {
		if err := json.Unmarshal(data, &chapterData); err != nil {
			// Bad data in the cache. Treat it as a miss, as the data format may have changed.
			slog.Warn("Bad chapter data in cache", "id", id)
		}
	}

	// If the data was loaded from the cache, just return that
	if !force && len(chapterData) > 0 {
		return chapterData, nil
	}

	// Query the Mangadex API for fresh data
	// TODO: add proper pagination support
	url := fmt.Sprintf("%s/manga/%s/feed?limit=300&translatedLanguage[]=en&contentRating[]=safe&contentRating[]=suggestive&contentRating[]=erotica&contentRating[]=pornographic&includeFutureUpdates=1&order[createdAt]=asc&order[updatedAt]=asc&order[publishAt]=asc&order[readableAt]=asc&order[volume]=asc&order[chapter]=asc",
		c.BaseURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching chapters feed for %s: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching chapters feed for %s returned status %s", id, resp.Status)
	}
	var feedData ChapterFeedResponse
	if err = json.NewDecoder(resp.Body).Decode(&feedData); err != nil {
		return nil, fmt.Errorf("decoding chapters feed for %s: %w", id, err)
	}

	// Cache the result
	if err := os.MkdirAll(filepath.Join(c.CacheDir, id), 0755); err != nil {
		return nil, fmt.Errorf("make cache dir '%s': %w", cachePath, err)
	}
	if data, err := json.MarshalIndent(feedData.Data, "", "  "); err != nil {
		return nil, fmt.Errorf("marshal chapters feed for %s: %w", id, err)
	} else {
		if err = os.WriteFile(cachePath, data, 644); err != nil {
			return nil, fmt.Errorf("cache chapters feed for %s: %w", id, err)
		}
	}

	return feedData.Data, nil
}
