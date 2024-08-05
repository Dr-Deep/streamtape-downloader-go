package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type progressReader struct {
	Reader io.Reader
	Size   int64
	Pos    int64
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	if err == nil {
		pr.Pos += int64(n)
		fmt.Printf("\rDownloading... %.2f%%", float64(pr.Pos)/float64(pr.Size)*100)
	}

	return n, err
}

func downloadFile(filePath string, downloadURL string) error {
	var (
		startTime            = time.Now().UnixMilli()
		tempDownloadFilePath = "." + filePath + ".temp"
	)

	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("StatusCode != 200")
	}
	defer resp.Body.Close()

	downloadFile, err := os.OpenFile(
		tempDownloadFilePath,
		os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil
	}
	defer downloadFile.Close()

	progressReader := &progressReader{
		Reader: resp.Body,
		Size:   resp.ContentLength,
	}

	if _, err := io.Copy(downloadFile, progressReader); err != nil {
		//os.Remove() downloadFile
		return err
	}

	os.Rename(tempDownloadFilePath, filePath)

	fmt.Println("=> Download completed")
	fmt.Printf("Took: %.2fs\n", float64(time.Now().UnixMilli()-startTime)/1000)

	return nil
}

func downloadFileWithYTDL(filePath string, downloadURL string) error {
	var cmd = exec.Command("yt-dlp", "-f mp4", "-c", "-o "+filePath, downloadURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	} else {
		return cmd.Wait()
	}
}
