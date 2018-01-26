package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/efy/placeholder"
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

	phOpts := &placeholder.ImageOptions{
		Width:  width,
		Height: height,
	}
	ph, err := placeholder.GenerateImage(phOpts)
	if err != nil {
		log.Println(err)
		return
	}

	rw.Header().Set("Content-Type", "image/png")
	rw.Header().Set("Content-Length", strconv.Itoa(len(*ph)))
	rw.Write(*ph)
}

// Extracts dimensions from the request path
func extractDimensions(path string) (int, int) {
	var w, h int
	re := regexp.MustCompile(`\/(\d+)x?(\d+)?`)
	m := re.FindAllStringSubmatch(path, 2)

	if len(m) == 0 {
		return w, h
	}

	if len(m[0]) == 2 {
		ws, _ := strconv.ParseInt(m[0][1], 10, 32)
		w = int(ws)
		h = w
	}

	if len(m[0]) == 3 {
		ws, _ := strconv.ParseInt(m[0][1], 10, 32)
		hs, _ := strconv.ParseInt(m[0][2], 10, 32)
		w = int(ws)
		h = int(hs)
	}

	return w, h
}

// Extracts color from the request path
func extractColor(path string) color.RGBA {
	c := color.RGBA{}

	re := regexp.MustCompile(`\/.+\/(.{6})`)
	m := re.FindAllStringSubmatch(path, 1)

	if len(m) == 0 {
		return c
	}

	if len(m[0]) <= 2 {
		col, err := hexToRGBA(m[0][1])
		if err != nil {
			return c
		}
		c = col
	}

	return c
}

func hexToRGBA(hex string) (color.RGBA, error) {
	var r, g, b uint8
	l := len(hex)
	c := color.RGBA{}

	if l != 3 && l != 6 {
		return c, fmt.Errorf("invalid hex string")
	}

	// Handle shorthand hex
	if l < 6 {
		hex += hex
	}

	if l == 6 {
		rs := hex[0:2]
		gs := hex[2:4]
		bs := hex[4:6]

		r64, err := strconv.ParseUint(rs, 16, 8)
		if err != nil {
			return c, fmt.Errorf("error converting %s", hex)
		}
		g64, err := strconv.ParseUint(gs, 16, 8)
		if err != nil {
			return c, fmt.Errorf("error converting %s", hex)
		}
		b64, err := strconv.ParseUint(bs, 16, 8)
		if err != nil {
			return c, fmt.Errorf("error converting %s", hex)
		}

		r = uint8(r64)
		g = uint8(g64)
		b = uint8(b64)
	}

	return color.RGBA{r, g, b, 255}, nil
}
