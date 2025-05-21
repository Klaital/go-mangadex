package main

import (
	"flag"
	"fmt"
	mangadex "github.com/klaital/go-mangadex/mangadex-client"
	"log/slog"
	"os"
)

const baseURL = "https://api.mangadex.org"
const cacheDir = "cache/describe"

var mangadexClient *mangadex.Client

func describeManga(mangaID string) error {
	md, err := mangadexClient.DescribeMangaFull(mangaID, false)
	if err != nil {
		slog.Error("Failed to describe manga", "err", err, "id", mangaID)
		return err
	}

	md.PrintDetails()
	return nil
}

func downloadChapters(mangaID, language string) error {
	slog.Error("Download function not implemented yet")
	return nil
}

func main() {
	describeCmd := flag.NewFlagSet("describe", flag.ExitOnError)
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	langFlag := downloadCmd.String("lang", "en", "Language code to filter chapters")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'describe' or 'download' subcommands")
		os.Exit(1)
	}

	mangadexClient = mangadex.NewClient(baseURL)

	switch os.Args[1] {
	case "describe":
		describeCmd.Parse(os.Args[2:])
		if describeCmd.NArg() < 1 {
			fmt.Println("Usage: describe <manga_id>")
			os.Exit(1)
		}
		mangaID := describeCmd.Arg(0)
		if err := describeManga(mangaID); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	case "download":
		downloadCmd.Parse(os.Args[2:])
		if downloadCmd.NArg() < 1 {
			fmt.Println("Usage: download [--lang=en] <manga_id>")
			os.Exit(1)
		}
		mangaID := downloadCmd.Arg(0)
		if err := downloadChapters(mangaID, *langFlag); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

	default:
		fmt.Println("Expected 'describe' or 'download' subcommands")
		os.Exit(1)
	}
}
