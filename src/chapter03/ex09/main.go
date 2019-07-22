package main

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math/cmplx"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		x := 2
		if qx, err := strconv.Atoi(q.Get("x")); err == nil && qx > 0 {
			x = qx
		}
		y := 2
		if qy, err := strconv.Atoi(q.Get("y")); err == nil && qy > 0 {
			y = qy
		}
		scale := 1.0
		if qs, err := strconv.ParseFloat(q.Get("scale"), 64); err == nil && qs > 0 {
			scale = qs
		}

		w.Header().Set("Content-Type", "image/png")
		draw(w, x, y, scale)
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func draw(w io.Writer, paramX, paramY int, scale float64) {
	var (
		xmin, ymin, xmax, ymax = -paramX, -paramY, +paramX, +paramY
		width, height          = 256 * scale, 256 * scale
	)

	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	for py := 0; float64(py) < height; py++ {
		y := float64(py)/height*(float64(ymax-ymin)) + float64(ymin)
		for px := 0; px < int(width); px++ {
			x := float64(px)/width*float64((xmax-xmin)) + float64(xmin)
			z := complex(float64(x), float64(y))

			img.Set(px, py, mandelbrot(z))
		}
	}
	png.Encode(w, img)
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - contrast*n}
		}
	}
	return color.Black
}