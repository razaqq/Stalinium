package downloader

import (
	"Stalinium/Installer/bridge"
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

func DownloadMod(modDir string, b *bridge.AppBridge) {
	url := "https://github.com/razaqq/Stalinium/releases/latest/download/StaliniumMod.zip"
	filePath := path.Join(os.TempDir(), path.Base(url))

	// create file
	out, err := os.Create(filePath)
	if err != nil {
		b.Error("Failed to create file path.")
		return
	}
	defer func() {
		if err := out.Close(); err != nil {
		}
	}()

	// download the file
	resp, err := http.Get(url)
	if err != nil {
		b.Error("Failed to download file.")
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
		}
	}()

	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		b.Error("Failed to parse file header.")
		return
	}

	// connect counter
	b.StartTime = time.Now()
	b.Total = uint64(size)
	if _, err = io.Copy(out, io.TeeReader(resp.Body, b)); err != nil {
		b.Error("Failed to download file.")
		return
	}

	if err := unzip(filePath, modDir); err != nil {
		b.Error("Failed to unpack mod.")
	}
	b.Success()
}

func unzip(src, dest string) error {
	zipReader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	for _, file := range zipReader.Reader.File {
		zippedFile, err := file.Open()
		if err != nil {
			return err
		}
		defer zippedFile.Close()

		finalPath := path.Join(dest, file.Name)

		if file.FileInfo().IsDir() {
			err := os.MkdirAll(finalPath, file.Mode())
			if err != nil {
				return err
			}
		} else {
			outputFile, err := os.OpenFile(finalPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer outputFile.Close()
			if _, err = io.Copy(outputFile, zippedFile); err != nil {
				return err
			}
		}
	}
	return nil
}
