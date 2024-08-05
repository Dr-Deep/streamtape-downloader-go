package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
)

// returns title, link, error
func GetStreamTapeVideoTitleAndURL(streamtapeVideoURL string) (string, string, error) {

	resp, err := http.Get(streamtapeVideoURL)
	if err != nil {
		return "", "", err
	}

	if resp.StatusCode != 200 {
		return "", "", errors.New("Server Response Code != 200")
	}

	_respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	respBody := string(_respBody)

	var (
		// Video Link
		videoLink = regexp.MustCompile(
			`/get_video\?id=[a-zA-Z0-9]+&expires=[a-zA-Z0-9]+&ip=[a-zA-Z0-9]+&token=[a-zA-Z0-9]+`,
		).FindString(respBody)

		// Video Title
		videoTitle = regexp.MustCompile(
			`.*<meta name=\"og:title\" content=\"(.*?)\">`,
		).FindStringSubmatch(respBody)[1]
	)

	if videoLink == "" {
		return "", "", errors.New("couldt get videoLink")
	}
	if videoTitle == "" {
		return "", "", errors.New("couldt get videoTitle")
	}

	url, err := url.Parse(streamtapeVideoURL)
	if err != nil {
		return "", "", errors.New("couldt parse StreamTape-Video-URL")
	}

	videoLink = fmt.Sprintf(
		"%s://%s%s",
		url.Scheme,
		url.Host,
		videoLink,
	)

	return videoTitle, videoLink, nil
}
