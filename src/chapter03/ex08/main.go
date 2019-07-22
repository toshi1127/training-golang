package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"math/big"
)

func main() {
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
			img.Set(px, py, mandelbrot64(z))
			img.Set(px, py, mandelbrot128(z))
			img.Set(px, py, mandelbrotBigFloat(z))
			img.Set(px, py, mandelbrotBigRat(z))
		}
	}
	png.Encode(os.Stdout, img) // NOTE: ignoring errors
}

func mandelbrot64(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	z64 := complex64(z)
	var v complex64
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z64
		if cmplx.Abs(complex128(v)) > 2 {
			return color.RGBA{255 - contrast*n, 0, 0, 255}
		}
	}
	return color.Black
}

func mandelbrot128(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.RGBA{255 - contrast*n, 0, 0, 255}
		}
	}
	return color.Black
}

func mandelbrotBigFloat(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	realZ := big.NewFloat(real(z))
	imagZ := big.NewFloat(imag(z))
	realV := new(big.Float)
	imagV := new(big.Float)
	for n := uint8(0); n < iterations; n++ {
		tempR := new(big.Float)
		tempI := new(big.Float)
		tempR.Mul(realV, realV).Sub(tempR, (&big.Float{}).Mul(imagV, imagV)).Add(tempR, realZ)
		tempI.Mul(realV, imagV).Mul(tempI, big.NewFloat(2)).Add(tempI, imagZ)
		realV, imagV = tempR, tempI
		sum := new(big.Float)
		sum.Mul(realV, realV).Add(sum, (&big.Float{}).Mul(imagV, imagV))
		if sum.Cmp(big.NewFloat(4)) == 1 {
			return color.Gray{255 - contrast*n}
		}
	}
	return color.Black
}

func mandelbrotBigRat(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	zz := big.NewRat(int64(real(z)), int64(imag(z)))
	v := big.NewRat(0, 1)
	for n := uint8(0); n < iterations; n++ {
		v.Mul(v, v)
		v.Add(v, zz)
		if v.Abs(v).Cmp(big.NewRat(2, 1)) >= 1 {
			return color.RGBA{255 - contrast*n, 0, 0, 255}
		}
	}
	return color.Black
}