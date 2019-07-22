package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

// Goにおいてのオブジェクトはメソッドを持つ単なる変数であり、メソッドは特定の型に関連付けられた関数です。
// Subpixelsはメソッドを持っていると。。。
type Subpixels []color.Color

func (pixels Subpixels) average() color.Color { // (レシーバ　型) 関数名（引数）戻り値の型
	var r, g, b, a uint8

	n := uint32(len(pixels))

	for _, pixel := range pixels {
		rloc, gloc, bloc, aloc := pixel.RGBA()
		r += uint8(rloc / n)
		g += uint8(gloc / n)
		b += uint8(bloc / n)
		a += uint8(aloc / n)
	}
	return color.RGBA{r, g, b, a}
}

func main() {
	const (
		xmin, xmax, ymin, ymax = -2, +2, -2, +2
		width, height          = 1024, 1024
		offx, offy             = ((xmax - xmin) / width), ((ymax - ymin) / height)
	)

	offpmx := []float64{-offx, offx}
	offpmy := []float64{-offy, offy}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		for px := 0; px < height; px++ {
			// ここ
			x := float64(px)/height*(xmax-xmin) + xmin
			subpixels := make(Subpixels, 0)
			for n := 0; n < 2; n++ {
				for m := 0; m < 2; m++ {
					subZ := complex(x+offpmx[n], y+offpmy[m])
					subpixels = append(subpixels, mandelbrot(subZ))
				}
			}

			avgColor := subpixels.average()
			img.Set(px, py, avgColor)
		}
	}

	png.Encode(os.Stdout, img)
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.RGBA{0, 255 - contrast*n, 0, 255}
		}
	}
	return color.Black
}