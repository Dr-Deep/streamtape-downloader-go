package main

import (
	"flag"
)

/*
   jeden block zählen, hash?
   downloads fortsetzen können
*/

var (
	streamtapeVideoURL = flag.String("streamtape-video-url", "", "valid StreamTape Video URL")
	downloadPath       = flag.String("ouput-file", "", "video output path")
)

func init() {
	flag.Parse()

	if *streamtapeVideoURL == "" {
		panic("not a valid StreamTape Video URL")
	}
}

func main() {
	if err := DownloadStreamtapeLink(*downloadPath, *streamtapeVideoURL); err != nil {
		panic(err)
	}
}
