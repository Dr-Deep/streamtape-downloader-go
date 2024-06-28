package main

// #cgo pkg-config: python-3.9-embed
// #include <Python.h>
import "C"

import (
	"encoding/json"
	"io"
	"net"
	"os"
	"path/filepath"
	"unsafe"
)

func DownloadStreamtapeLink(downloadPath string, streamtapeVideoURL string) error {

	title, downloadURL, err := GetStreamTapeVideoTitleAndURL(streamtapeVideoURL)
	if err != nil {
		return err
	}

	if downloadPath == "" {
		downloadPath = filepath.Clean(title) // || $HOME/Downloads/title
	}

	return downloadFile(downloadPath, downloadURL)
}

// könnte panicen
// returns title, downloadURL
func GetStreamTapeVideoTitleAndURL(streamtapeVideoURL string) (string, string, error) {

	// wir kommunizieren über unix sockets && env vars
	// zu faul für C-Python-API docs; Sorry

	var unixSocketPath = "/tmp/serien-downloader.sock"

	socket, err := net.Listen("unix", unixSocketPath)
	if err != nil {
		return "", "", err
	}
	defer os.Remove(unixSocketPath)

	var rawJson []byte
	go func() {
		conn, err := socket.Accept()
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		// lesen
		buf, err := io.ReadAll(conn)
		if err != nil {
			panic(err)
		}

		rawJson = buf
	}()

	py_scraper(unixSocketPath, streamtapeVideoURL)

	//
	type sus struct {
		StreamTapeVideoTitle       string
		StreamTapeVideoDownloadURL string
	}
	var susValue sus
	if err := json.Unmarshal(rawJson, &susValue); err != nil {
		return "", "", err
	}

	return susValue.StreamTapeVideoTitle, susValue.StreamTapeVideoDownloadURL, nil
}

func py_scraper(unixSocketPath string, streamtapeVideoURL string) {
	os.Setenv("StreamTapeVideoURL", streamtapeVideoURL)
	os.Setenv("unixSocketPath", unixSocketPath)
	defer os.Unsetenv("StreamTapeVideoURL")
	defer os.Unsetenv("unixSocketPath")

	var pyScript = C.CString(`
import re
import os
import sys
import json
import socket
import requests

streamtapeVideoURL = os.environ["StreamTapeVideoURL"]
unix_sock_path = os.environ["unixSocketPath"]

# stolen from  https://github.com/fluffysatoshi/streamtape2curl/blob/master/streamtape2curl.py
html = requests.get(streamtapeVideoURL).content.decode()
token = re.match(r".*document.getElementById.*\('norobotlink'\).innerHTML =.*?token=(.*?)'.*?;", html, re.M|re.S).group(1)
infix=re.match(r'.*<div id="ideoooolink" style="display:none;">(.*?token=).*?<[/]div>', html, re.M|re.S).group(1)

downloadURL=f'http:/{infix}{token}'
title=re.match(r'.*<meta name="og:title" content="(.*?)">', html, re.M|re.S).group(1)

# hier fängts an
client = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
client.connect(unix_sock_path)
client.sendall(json.dumps(
    {
        "StreamTapeVideoTitle": title,
        "StreamTapeVideoDownloadURL": downloadURL,
}
).encode())
client.close()
`)

	C.Py_Initialize()
	C.PyRun_SimpleString(pyScript)
	C.free(unsafe.Pointer(pyScript))
	C.Py_Finalize()
}
