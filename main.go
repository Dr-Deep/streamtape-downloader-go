package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

/*
   downloads fortsetzen können / bytes zählen
   golang download lib? statt yt-dlp?
   groöße am ende vergleichen ob download wirklich ganz ist
   wenn yt-dlp installiert ist das benutzen wenn net build-in dl
*/

func main() {

	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "append StreamTape-Links to command line")
		panic("no StreamTape-Video-URLS found")
	}

	fmt.Printf("=> Downloading %v StreamTape-Videos\n", len(os.Args[1:]))
	for i, streamTapeVideoURL := range os.Args[1:] {
		var startTime = time.Now().UnixMilli()

		// Fetching URL
		fmt.Printf("=>  (%v/%v) Fetching StreamTape-Video-URL\n", i+1, len(os.Args[1:]))
		title, url, err := GetStreamTapeVideoTitleAndURL(streamTapeVideoURL)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"=> (%v/%v) Skipping Video because of an Error: %s\n",
				i+1,
				len(os.Args[1:]),
				err.Error(),
			)
			continue
		}

		// Downloading
		fmt.Printf("=> (%v/%v) Downloading %s\n", i+1, len(os.Args[1:]), title)
		if err := downloadFileWithYTDL(
			filepath.Clean(title),
			url,
		); err != nil {
			fmt.Fprintf(
				os.Stderr,
				"=> (%v/%v) Skipping Video because of an Error: %s\n",
				i+1,
				len(os.Args[1:]),
				err.Error(),
			)
			continue
		}

		// Finish
		fmt.Printf(
			"=> (%v/%v) Download finishied in %.2fs\n",
			i+1,
			len(os.Args[1:]),
			float64(time.Now().UnixMilli()-startTime)/1000,
		)
	}

	fmt.Printf("=> (%v/%v) Done!", len(os.Args[1:]), len(os.Args[1:]))
}
