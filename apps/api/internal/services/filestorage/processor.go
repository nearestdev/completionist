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

// ProcessImage resizes the image to MaxWidth if needed, then encodes it.
// PNG inputs are preserved as PNG; everything else is encoded as JPEG.
func (p *ImageProcessor) ProcessImage(r io.Reader, filename string) (*bytes.Buffer, string, error) {
	img, err := imaging.Decode(r)
	if err != nil {
		return nil, "", err
	}

	if img.Bounds().Dx() > p.MaxWidth {
		img = imaging.Resize(img, p.MaxWidth, 0, imaging.Lanczos)
	}

	ext := strings.ToLower(filepath.Ext(filename))
	buf := new(bytes.Buffer)
	var contentType string
	if ext == ".png" {
		err = imaging.Encode(buf, img, imaging.PNG, imaging.PNGCompressionLevel(png.BestCompression))
		contentType = "image/png"
	} else {
		err = jpeg.Encode(buf, img, &jpeg.Options{Quality: p.Quality})
		contentType = "image/jpeg"
	}

	if err != nil {
		return nil, "", err
	}

	return buf, contentType, nil
}
