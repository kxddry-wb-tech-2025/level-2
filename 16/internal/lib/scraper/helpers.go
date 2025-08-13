package scraper

import (
	"io"
	"net/url"
	"os"
	"path/filepath"
	"wget/internal/lib/parser"
)

func saveFile(url *url.URL, content io.ReadCloser, contentType, outputDir string) (string, error) {
	defer content.Close()

	full := parser.GenerateLocalPath(url, contentType, outputDir)

	dir := filepath.Dir(full)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	file, err := os.Create(full)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, content)
	if err != nil {
		return "", err
	}
	return full, nil
}
