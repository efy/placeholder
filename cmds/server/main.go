package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	addr = flag.String("addr", ":9000", "listen address")
)

func main() {
	flag.Parse()

	// Handle setting the port from an environment variable
	// when running in a heroku like environment
	port := os.Getenv("PORT")
	if port != "" {
		port = ":" + port
		addr = &port
	}

	log.Println("starting image placeholder server on port", *addr)

	http.HandleFunc("/", placeholderHandler)

	log.Fatal(http.ListenAndServe(*addr, nil))
}

func placeholderHandler(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	width := 600
	height := 400

	if len(parts) >= 2 {
		dimensions := strings.Split(parts[1], "x")
		le := len(dimensions)

		if le == 2 {
			w, err := strconv.ParseInt(dimensions[0], 10, 32)
			if err != nil {
				log.Println("error converting width to integer")
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			h, err := strconv.ParseInt(dimensions[1], 10, 32)
			if err != nil {
				log.Println("error converting height to integer")
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			width = int(w)
			height = int(h)
		}

		if le == 1 {
			w, err := strconv.ParseInt(dimensions[0], 10, 32)
			if err != nil {
				log.Println("error converting width to integer")
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			width = int(w)
			height = int(w)
		}
	}

	phOpts := &PlaceholderOptions{
		Width:  width,
		Height: height,
	}
	ph, err := GeneratePlaceholder(phOpts)
	if err != nil {
		log.Println(err)
		return
	}

	rw.Header().Set("Content-Type", "image/png")
	rw.Header().Set("Content-Length", strconv.Itoa(len(*ph)))
	rw.Write(*ph)
}

type PlaceholderOptions struct {
	Width  int
	Height int
}

// GeneratePlaceholder creates and encodes a png image with the given placeholder options
func GeneratePlaceholder(opts *PlaceholderOptions) (*[]byte, error) {
	// Create the blank graphic
	m := image.NewRGBA(image.Rect(0, 0, opts.Width, opts.Height))
	color := color.RGBA{0, 0, 0, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)

	// Encode as a png
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, m); err != nil {
		return nil, fmt.Errorf("error encoding image")
	}
	byt := buf.Bytes()
	return &byt, nil
}
