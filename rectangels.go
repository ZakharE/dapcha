package main

import (
	"bufio"
	"image"
	"image/draw"
	"image/png"
	"os"
)

func main1() {
	rgbaImage := image.NewRGBA(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{800, 600},
	})

	f, err := os.Create("rectangles.png")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	draw.Draw(rgbaImage, rgbaImage.Bounds(), image.White, image.Point{}, draw.Src)
	for i := 0; i < 800; i += 20 {
		pt := image.Point{i, 0}
		vLine := image.NewRGBA(image.Rectangle{
			Min: image.Point{i, 0},
			Max: image.Point{i + 5, 600},
		})

		draw.Draw(rgbaImage, vLine.Bounds(), image.Black, pt, draw.Src)

	}

	b := bufio.NewWriter(f)
	png.Encode(b, rgbaImage)
	b.Flush()
}
