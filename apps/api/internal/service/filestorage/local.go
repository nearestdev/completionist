package filestorage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type LocalStorage struct {
	UploadDir string
	BaseURL   string
}

func NewLocalStorage(uploadDir, baseURL string) (*LocalStorage, error) {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{
		UploadDir: uploadDir,
		BaseURL:   baseURL,
	}, nil
}

func (s *LocalStorage) Upload(ctx context.Context, file io.Reader, filename string, contentType string) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(filename)
	newFilename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.New().String(), ext)
	dstPath := filepath.Join(s.UploadDir, newFilename)

	// Create destination file
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy content
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	// Return public URL
	return fmt.Sprintf("%s/%s", s.BaseURL, newFilename), nil
}

func (s *LocalStorage) Delete(ctx context.Context, key string) error {
    // Basic implementation: extract filename from key/URL and delete
    // For now, assuming key is the filename or full URL
    // TODO: Parse key correctly
    return nil
}
