package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// URLからパラメータ値を読み取れるようにする

var palette = []color.Color{color.White, color.Black}

const (
	whilteIndex = 0
	blackIndex  = 1
)

// lissajousに渡すパラメータ
type Param struct {
	cycles float64
	res float64
	size int
	nframes int
	delay int
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// パラメータの初期値をセット
		p := Param{5.0, 0.001, 100, 64, 8}
		// クエリがセットされていた場合は上書きする
		for k, vals := range r.URL.Query() {
			v := vals[0]
			switch k {
			case "cycles":
				p.cycles, _ = strconv.ParseFloat(v, 64)
			case "res":
				p.res, _ = strconv.ParseFloat(v, 64)
			case "size":
				p.size, _ = strconv.Atoi(v)
			case "nframes":
				p.nframes, _ = strconv.Atoi(v)
			case "delay":
				p.delay, _ = strconv.Atoi(v)
			}
		}

		lissajous(w, p)
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func lissajous(out io.Writer, p Param) {
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: p.nframes}
	phase := 0.0
	for i := 0; i < p.nframes; i++ {
		rect := image.Rect(0, 0, 2*p.size+1, 2*p.size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < p.cycles*2*math.Pi; t += p.res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(p.size+int(x*float64(p.size)+0.5), p.size+int(y*float64(p.size)+0.5), blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, p.delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
