package placeholder

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

type ImageOptions struct {
	Width  int
	Height int
	Color  color.RGBA
}

var (
	DefaultImageOptions = &ImageOptions{
		Width:  600,
		Height: 400,
		Color:  color.RGBA{0, 0, 0, 255},
	}
)

// baseImage creates an RGBA image of the specified size and color
func baseImage(w int, h int, c color.RGBA) image.Image {
	m := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(m, m.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)

	return m
}

// GenerateImage creates and encodes a png image with the given placeholder options
func GenerateImage(opts *ImageOptions) (*[]byte, error) {
	// Create the blank graphic
	m := baseImage(opts.Width, opts.Height, opts.Color)

	// Encode as a png
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, m); err != nil {
		return nil, fmt.Errorf("error encoding image")
	}
	byt := buf.Bytes()
	return &byt, nil
}
