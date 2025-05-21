package mangadex_client

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

type Manga struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Title                          map[string]string   `json:"title"`
		AltTitles                      []map[string]string `json:"altTitles"`
		Description                    map[string]string   `json:"description"`
		IsLocked                       bool                `json:"isLocked"`
		Links                          map[string]string   `json:"links"`
		OriginalLanguage               string              `json:"originalLanguage"`
		LastVolume                     string              `json:"lastVolume"`
		LastChapter                    string              `json:"lastChapter"`
		PublicationDemographic         string              `json:"publicationDemographic"`
		Status                         string              `json:"status"`
		Year                           int                 `json:"year"`
		ContentRating                  string              `json:"contentRating"`
		ChapterNumbersResetOnNewVolume bool                `json:"chapterNumbersResetOnNewVolume"`
		AvailableTranslatedLanguages   []string            `json:"availableTranslatedLanguages"`
		LatestUploadedChapter          string              `json:"latestUploadedChapter"`
		Tags                           []struct {
			Id         string `json:"id"`
			Type       string `json:"type"`
			Attributes struct {
				Name        map[string]string `json:"name"`
				Description map[string]string `json:"description"`
				Group       string            `json:"group"`
				Version     int               `json:"version"`
			} `json:"attributes"`
			Relationships []struct {
				Id         string `json:"id"`
				Type       string `json:"type"`
				Related    string `json:"related"`
				Attributes struct {
				} `json:"attributes"`
			} `json:"relationships"`
		} `json:"tags"`
		State     string `json:"state"`
		Version   int    `json:"version"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	} `json:"attributes"`
	Relationships []struct {
		Id         string `json:"id"`
		Type       string `json:"type"`
		Related    string `json:"related"`
		Attributes struct {
		} `json:"attributes"`
	} `json:"relationships"`
}

type MangaResponse struct {
	Result   string `json:"result"`
	Response string `json:"response"`
	Data     Manga  `json:"data"`
}

// DescribeManga downloads the base manga details for a given series.
func (c Client) DescribeManga(id string, force bool) (*MangaResponse, error) {
	cachePath := filepath.Join(c.CacheDir, id, "manga.json")

	var mangaResp MangaResponse
	if data, err := os.ReadFile(cachePath); err != nil {
		// handle cache miss
		slog.Debug("Cache miss", "id", id)
	} else {
		if err := json.Unmarshal(data, &mangaResp); err != nil {
			// Bad data in the cache. Treat it as a miss, as the data format may have changed.
			slog.Warn("Bad data in cache", "id", id)
		}
	}

	// If the data was loaded from the cache, just return that
	if !force && mangaResp.Data.Id != "" {
		return &mangaResp, nil
	}

	// Fetch the manga's base data from the API
	url := fmt.Sprintf("%s/manga/%s", c.BaseURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching manga details for %s: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch manga: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&mangaResp)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// Cache the result
	if err := os.MkdirAll(filepath.Join(c.CacheDir, id), 0755); err == nil {
		if data, err := json.MarshalIndent(mangaResp, "", "  "); err == nil {
			if err = os.WriteFile(cachePath, data, 0644); err != nil {
				return nil, fmt.Errorf("error creating cache file: %w", err)
			}
		}
	}

	return &mangaResp, nil
}
