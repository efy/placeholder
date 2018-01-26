package main

import (
	"flag"
	"log"
	"net/http"
	"os"
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
