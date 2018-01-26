package main

import (
	"image/color"
	"testing"
)

func TestExtractDimensions(t *testing.T) {
	tt := []struct {
		in string
		w  int
		h  int
	}{
		{
			"/400x400",
			400,
			400,
		},
		{
			"/400",
			400,
			0,
		},
		{
			"/abcxabc",
			0,
			0,
		},
	}

	for _, tr := range tt {
		t.Run(tr.in, func(t *testing.T) {
			w, h := extractDimensions(tr.in)
			if w != tr.w {
				t.Error("expected:", tr.w)
				t.Error("got     ", w)
			}
			if h != tr.h {
				t.Error("expected:", tr.h)
				t.Error("got     ", h)
			}
		})
	}
}

func TestHexToRGBA(t *testing.T) {
	tt := []struct {
		in  string
		out color.RGBA
	}{
		{
			"ffffff",
			color.RGBA{255, 255, 255, 255},
		},
	}

	for _, tr := range tt {
		t.Run(tr.in, func(t *testing.T) {
			col, _ := hexToRGBA(tr.in)
			if tr.out != col {
				t.Error("expected", tr.out)
				t.Error("got     ", col)
			}
		})
	}
}

func TestExtractColor(t *testing.T) {
	tt := []struct {
		in  string
		out color.RGBA
	}{
		{
			"/400/ffffff",
			color.RGBA{255, 255, 255, 255},
		},
	}

	for _, tr := range tt {
		t.Run(tr.in, func(t *testing.T) {
			col := extractColor(tr.in)
			if tr.out != col {
				t.Error("expected", tr.out)
				t.Error("got     ", col)
			}
		})
	}
}
