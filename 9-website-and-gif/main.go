package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe("localhost:8000", nil)
	log.Fatal(err)
}

func handler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cyclesSlice, ok := r.Form["cycles"]
	if !ok {

	}
	cycles := 5
	if len(cyclesSlice) == 1 {
		var err error
		cycles, err = strconv.Atoi(cyclesSlice[0])
		if err != nil {
			log.Fatal(err)
		}
	}
	lissajous(w, cycles)
}

func lissajous(output io.Writer, cycles int) {
	const (
		nframes    = 64     // number of animation frames
		size       = 100    // radius of width and height
		resolution = 0.0001 // angular resolution
	)
	freq := rand.Float64()
	fmt.Fprintf(os.Stderr, "%v\n", freq)
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0

	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, []color.Color{color.Black, color.RGBA{0x00, 0xFF, 0x00, 0xFF}})
		for t := 0.0; t < 2*math.Pi*float64(cycles); t += resolution {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)

			const blackIndex = 1
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), blackIndex)
		}
		phase += 0.1

		const delay = 8 // delay between frames in 10ms units
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(output, &anim)
}
