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
		in        string
		out       color.RGBA
		shoulderr bool
	}{
		{
			"ffffff",
			color.RGBA{255, 255, 255, 255},
			false,
		},
		{
			"fff",
			color.RGBA{0, 0, 0, 255},
			true,
		},
	}

	for _, tr := range tt {
		t.Run(tr.in, func(t *testing.T) {
			col, err := hexToRGBA(tr.in)
			if tr.out != col {
				t.Error("expected", tr.out)
				t.Error("got     ", col)
			}

			if err != nil && tr.shoulderr == false {
				t.Error("unexpected error:", err)
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
