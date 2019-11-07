package main

import (
	"flag"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var toJPEG = flag.Bool("jpeg", false, "jpeg")
var toPNG = flag.Bool("png", false, "png")
var toGIF = flag.Bool("gif", false, "gif")

func main() {
	flag.Parse()

	for _, fname := range flag.Args() {
		handleFile(fname)
	}
}

func handleFile(fname string) {
	f, err := os.Open(fname)
	if err != nil {
		log.Print(err)
		return
	}
	defer f.Close()

	img, kind, err := image.Decode(f)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print("Input format = ", kind)

	basename := strings.TrimSuffix(fname, filepath.Ext(fname))

	if *toJPEG && kind != "jpeg" {
		outFname := basename + ".jpeg"
		out, err := os.Create(outFname)
		if err != nil {
			log.Print(err)
		} else {
			defer out.Close()
			err = jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
			if err != nil {
				log.Print(err)
			}
			log.Printf("created %q", outFname)
		}
	}
	if *toPNG && kind != "png" {
		outFname := basename + ".png"
		out, err := os.Create(outFname)
		if err != nil {
			log.Print(err)
		} else {
			defer out.Close()
			err = png.Encode(out, img)
			if err != nil {
				log.Print(err)
			}
			log.Printf("created %q", outFname)
		}
	}
	if *toGIF && kind != "gif" {
		outFname := basename + ".gif"
		out, err := os.Create(outFname)
		if err != nil {
			log.Print(err)
		} else {
			defer out.Close()
			err = gif.Encode(out, img, &gif.Options{NumColors: 256})
			if err != nil {
				log.Print(err)
			}
			log.Printf("created %q", outFname)
		}
	}
}