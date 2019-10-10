package main

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"math/cmplx"
	"os"
	"sync"
)

func main() {
	//renderWithGoRoutine(os.Stdout)
	renderWithGoRoutine2(os.Stdout)
}

func renderWithGoRoutine2(out io.Writer) {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 1024, 1024
	)

	type point struct {
		px, py int
		z      complex128
	}

	points := make(chan point, 1024*1024)
	go func() {
		for py := 0; py < height; py++ {
			y := float64(py)/height*(ymax-ymin) + ymin
			for px := 0; px < width; px++ {
				x := float64(px)/width*(xmax-xmin) + xmin
				z := complex(x, y)
				points <- point{px, py, z}
			}
		}
		close(points)
	}()

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range points {
				// 画像の点 (px, py) は複素数値zを表している
				img.Set(p.px, p.py, mandelbrot(p.z))
			}
		}()
	}

	wg.Wait()
	png.Encode(out, img) // 注意: エラーを無視
}

func renderWithGoRoutine(out io.Writer) {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 1024, 1024
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// 平行な計算を制限するための係数セマフォ
	tokens := make(chan struct{}, 100)
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			go func(px, py int, z complex128) {
				tokens <- struct{}{}
				// 画像の点 (px, py) は複素数値zを表している
				img.Set(px, py, mandelbrot(z))
				<-tokens
			}(px, py, z)
		}
	}
	png.Encode(out, img) // 注意: エラーを無視
}

func render(out io.Writer) {
	const (
		xmin, ymin, xmax, ymax = -2, -2, +2, +2
		width, height          = 1024, 1024
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			z := complex(x, y)
			// 画像の点 (px, py) は複素数値zを表している
			img.Set(px, py, mandelbrot(z))
		}
	}
	png.Encode(out, img) // 注意: エラーを無視
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