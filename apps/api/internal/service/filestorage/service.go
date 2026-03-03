package filestorage

import (
	"context"
	"io"
)

type Service interface {
    // Upload uploads a file to the storage provider
    Upload(ctx context.Context, file io.Reader, filename string, contentType string) (string, error)
    
    // Delete deletes a file from the storage provider (optional)
    Delete(ctx context.Context, key string) error
}
