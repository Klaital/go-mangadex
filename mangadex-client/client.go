package mangadex_client

import (
	"fmt"
	"log/slog"
	"strconv"
)

type Client struct {
	CacheDir    string
	BaseURL     string
	DownloadDir string
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:     baseURL,
		CacheDir:    "cache",
		DownloadDir: "download",
	}
}

type MangaDetails struct {
	Manga
	Chapters []ChapterFeedData
}

func MustAtof(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		slog.Error("Failed to parse chapter number", "chapter", s, "err", err)
	}
	return f
}

func (md MangaDetails) GetLatestChapter() *ChapterFeedData {
	var newestChapterNum float64 = 0.0
	var newestChapter *ChapterFeedData
	for _, chapter := range md.Chapters {
		chapterOrdinal := MustAtof(chapter.Attributes.Chapter)
		if chapterOrdinal > newestChapterNum {
			newestChapterNum = chapterOrdinal
			newestChapter = &chapter
		}
	}
	// if nil, then we failed to find that last chapter
	return newestChapter
}

func (c Client) DescribeMangaFull(id string, force bool) (*MangaDetails, error) {
	var md MangaDetails
	mr, err := c.DescribeManga(id, force)
	if err != nil {
		return nil, fmt.Errorf("describe manga: %w", err)
	}

	md.Manga = mr.Data

	// Add the chapter data
	md.Chapters, err = c.ListChapters(id, force)
	if err != nil {
		return nil, fmt.Errorf("list chapters: %w", err)
	}

	return &md, nil
}

func (md MangaDetails) PrintDetails() {
	fmt.Printf("%s    (%s)\n", md.Manga.Attributes.Title["en"], md.Id)
	fmt.Printf("%d chapters\n", len(md.Chapters))
	if md.Manga.Attributes.LatestUploadedChapter == "" {
		fmt.Println("No latest chapter specified")
	} else {
		latestChapter := md.GetLatestChapter()
		if latestChapter == nil {
			fmt.Printf("Latest chapter %s not found\n", md.Manga.Attributes.LatestUploadedChapter)
		} else {
			fmt.Printf("Latest: Vol %s, Chapter %s - %s (uploaded %s)\n", latestChapter.Attributes.Volume, latestChapter.Attributes.Chapter, latestChapter.Attributes.Title, latestChapter.Attributes.ReadableAt)
		}
	}
}
