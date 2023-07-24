package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

func main() {
	f, err := os.Open("font.ttf")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fontBytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Font %s was loaded\n", font.Name(truetype.NameIDFontFullName))

	rgbaImage := image.NewRGBA(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{800, 600},
	})

	draw.Draw(rgbaImage, rgbaImage.Bounds(), image.White, image.ZP, draw.Src)

	res, err := os.Create("result.png")
	if err != nil {
		panic(err)
	}

	text := generateCapcha()
	c := freetype.NewContext()
	c.SetSrc(image.Black)
	c.SetFont(font)
	c.SetDPI(128)
	c.SetFontSize(20)
	c.SetClip(rgbaImage.Bounds())
	c.SetDst(rgbaImage)

	c1 := freetype.NewContext()
	c1.SetSrc(image.NewUniform(color.RGBA{
		R: 0xff,
		G: 0,
		B: 0,
		A: 0xff,
	}))
	c1.SetFont(font)
	c1.SetDPI(128)
	c1.SetFontSize(20)
	c1.SetClip(rgbaImage.Bounds())
	c1.SetDst(rgbaImage)

	pt := freetype.Pt(10, int(c.PointToFixed(20)>>6))

	for i, r := range text {
		if i%2 == 0 {
			pt, err = c1.DrawString(string(r), pt)
		} else {
			pt, err = c.DrawString(string(r), pt)
		}

		if err != nil {
			panic(err)
		}
		y := rand.Intn(rgbaImage.Bounds().Max.Y)
		pt.Y = freetype.Pt(0, int(c.PointToFixed(float64(y))>>6)).Y
	}

	err = png.Encode(res, rgbaImage)
	if err != nil {
		panic(err)
	}
}

func generateCapcha() string {
	return "Hello, World!"
}
