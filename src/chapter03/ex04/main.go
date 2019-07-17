package main

import (
	"math"
	"fmt"
	"net/http"
	"strconv"
	"log"
	"io"
)

const (
	defaultWidth, defaultHeight = 600, 320
	cells                       = 100
	xyrange                     = 30.0
	angle                       = math.Pi / 6
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)
var width, height int
var xyscale, zscale float64

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var err error

		widthParam := r.FormValue("width");
		if widthParam != "" {
			width, err = strconv.Atoi(widthParam);
			if err != nil {
				fmt.Printf("Error while parsing param %s to int: %v", widthParam, err)
			}
		} else {
			width = defaultWidth
		}

		heightParam := r.FormValue("height");
		if heightParam != "" {
			height, err = strconv.Atoi(heightParam);
			if err != nil {
				fmt.Printf("Error while parsing param %s to int: %v", heightParam, err)
			}
		} else {
			height = defaultWidth
		}

		xyscale = float64(width) / 2 / xyrange
		zscale = float64(height) * 0.4

		w.Header().Set("Content-Type", "image/svg+xml")
		createSvg(w)
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func createSvg(out io.Writer) {
	fmt.Fprintf(out, "<svg	xmlns='http://www.w3.org/2000/svg' style='stroke: grey; strokewidth: 0.7' width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, belowZero := corner(i+1, j)
			bx, by, _ := corner(i, j)
			cx, cy, _ := corner(i, j+1)
			dx, dy, _ := corner(i+1, j+1)
			var color string
			if (belowZero) {
				color = "#ff0000"
			} else {
				color = "#0000ff"
			}
			fmt.Fprintf(out, "<polygon style='fill: %s' points='%g, %g %g, %g %g, %g %g, %g'/>\n",
				color, ax, ay, bx, by, cx, cy, dx, dy)

		}
	}
	fmt.Fprintln(out, "</svg>")
}
func isValid(number float64) bool {
	if math.IsInf(number, 0) || math.IsNaN(number) {
		return false
	} else {
		return true
	}
}

func corner(i, j int) (float64, float64, bool) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	z := f(x, y)

	sx := float64(width)/2 + (x-y)*cos30*xyscale
	sy := float64(height)/2 + (x+y)*sin30*xyscale - z*zscale
	belowZero := z < 0
	return sx, sy, belowZero
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y)
	result := math.Sin(r) / r
	if isValid(result) {
		return result
	} else {
		return 0
	}
}