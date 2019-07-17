package main

import (
	"os"
	"fmt"
	"math"
	"io/ioutil"
)

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
	method := "eggbox"
	if len(os.Args) == 2 {
		method = os.Args[1]
	}
	a := fmt.Sprintf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j, method)
			bx, by := corner(i, j, method)
			cx, cy := corner(i, j+1, method)
			dx, dy := corner(i+1, j+1, method)
			if math.IsNaN(ax) || math.IsNaN(ay) || math.IsNaN(bx) || math.IsNaN(by) ||
				math.IsNaN(cx) || math.IsNaN(cy) || math.IsNaN(dx) || math.IsNaN(dy) {
				continue
			}
			fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g'/>\n", ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	a += "</svg>"
	ioutil.WriteFile("./data.svg", []byte(a), 0666)
	fmt.Println("</svg>")
}

func corner(i, j int, method string) (float64, float64) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	var z float64
	switch method {
	case "eggbox":
		z = eggBox(x, y)
	case "moguls":
		z = moguls(x, y)
	case "saddle":
		z = saddle(x, y)
	default:
		fmt.Fprintln(os.Stderr, "Usage: ex02 eggbox|moguls|saddle")
		os.Exit(1)
	}

	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

func eggBox(x, y float64) float64 {
	return math.Cos(x) * math.Sin(y) / 5
}

func moguls(x, y float64) float64 {
	return x
}

func saddle(x, y float64) float64 {
	return y
}