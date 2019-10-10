package main

import (
	"fmt"
	"math"
	"io/ioutil"
)

// math.Sin(r)/rの値が有限であるかをチェックする

const (
	width, height = 600, 320
	cells         = 100
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func main() {
	var a string
	a = fmt.Sprintf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			fmt.Printf("<polygon points='%g, %g %g, %g %g, %g %g, %g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	a += "</svg>"
	ioutil.WriteFile("./data.svg", []byte(a), 0666)
	fmt.Println("</svg>")
}

func isValid(number float64) bool {
	if math.IsInf(number, 0) || math.IsNaN(number) {
		return false
	} else {
		return true
	}
}

func corner(i, j int) (float64, float64) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	z := f(x, y)

	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y)
	result := math.Sin(r)/r
	if isValid(result) {
		return result
	} else {
		return 0 // 0を返してはいけないぞ
	}
}