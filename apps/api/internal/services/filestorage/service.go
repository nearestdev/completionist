package filestorage

import (
	"context"
	"io"
)

type Service interface {
	Upload(ctx context.Context, file io.Reader, filename string, contentType string) (string, error)
	Delete(ctx context.Context, key string) error
}
