package filestorage

import (
	"bytes"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type ImageProcessor struct {
    MaxWidth int
    Quality  int
}

func NewImageProcessor(maxWidth, quality int) *ImageProcessor {
    return &ImageProcessor{
        MaxWidth: maxWidth,
        Quality:  quality,
    }
}

// ProcessImage decodes, resizes, and encodes the image.
// It returns a Buffer containing the processed image data.
func (p *ImageProcessor) ProcessImage(r io.Reader, filename string) (*bytes.Buffer, string, error) {
    // Decode
    img, err := imaging.Decode(r)
    if err != nil {
        return nil, "", err
    }

    // Resize if needed
    if img.Bounds().Dx() > p.MaxWidth {
        img = imaging.Resize(img, p.MaxWidth, 0, imaging.Lanczos)
    }

    // Enocde to JPEG or keep original?
    // Let's force convert to JPEG for consistency and compression unless it's transparent?
    // For now, let's just stick to JPEG for compression for everything unless user asks otherwise
    // Check extension
    ext := strings.ToLower(filepath.Ext(filename))
    
    buf := new(bytes.Buffer)
    
    // If it's PNG or WEBP, we might want to keep transparency? 
    // Simplify: Convert everything to JPEG for now to ensure compression, 
    // unless explicitly requested to handle transparency in the future.
    // Actually, let's try to preserve PNG if input is PNG, but compress it.
    // disintegration/imaging supports Encode
    
    var contentType string
    if ext == ".png" {
        // Encode as PNG
         err = imaging.Encode(buf, img, imaging.PNG, imaging.PNGCompressionLevel(png.BestCompression))
         contentType = "image/png"
    } else {
        // Default to JPEG
        err = jpeg.Encode(buf, img, &jpeg.Options{Quality: p.Quality})
        contentType = "image/jpeg"
    }

    if err != nil {
        return nil, "", err
    }

    return buf, contentType, nil
}
