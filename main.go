package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"io"
	"math/rand"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

var (
	DPI         = 62
	FrameNumber = 48
	FontSizeMax = 33
	FontSizeMin = 32
	ImgH        = 100
	ImgW        = 180
)

var (
	ColorRed = color.RGBA{
		R: 0xff,
		G: 0,
		B: 0,
		A: 0xff,
	}
	ColorGreen = color.RGBA{
		R: 0,
		G: 0xff,
		B: 0,
		A: 0xff,
	}
	ColorBlue = color.RGBA{
		R: 0,
		G: 0,
		B: 0xff,
		A: 0xff,
	}
)

type Letter struct {
	letter       rune
	font         *freetype.Context
	size         int
	prevPosition fixed.Point26_6
	direction    int
}

type Frame = []*Letter

func main() {
	f, err := os.Open("font.ttf")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fontBytes, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Font %s was loaded\n", font.Name(truetype.NameIDFontFullName))

	// rgbaImage := image.NewRGBA(image.Rectangle{
	// 	Min: image.Point{0, 0},
	// 	Max: image.Point{ImgW, ImgH},
	// })
	pallete := color.Palette{color.White, color.Black, ColorBlue, ColorRed}
	baseIMG := image.NewPaletted(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{ImgW, ImgH},
	}, pallete)

	draw.Draw(baseIMG, baseIMG.Bounds(), image.White, image.Point{}, draw.Src)

	fmt.Printf("max bound: %d\n", baseIMG.Bounds().Max.Y)

	text := generateCapcha()
	fontContexts := generateRandomFonts(4, font, baseIMG)
	maxY := baseIMG.Bounds().Max.Y
	pt := freetype.Pt(0, int(fontContexts[0].PointToFixed(float64(maxY/4))>>6))
	letters := make([]*Letter, 0, len(text))
	dir := 1
	for i := range text {
		idx := rand.Intn(len(fontContexts))
		fontCtx := fontContexts[idx]
		letter := &Letter{
			letter:       rune(text[i]),
			font:         fontCtx,
			prevPosition: pt,
			direction:    (-1) * dir,
		}
		dir *= -1

		letters = append(letters, letter)
		pt, err = fontCtx.DrawString(string(text[i]), pt)
		if err != nil {
			panic(err)
		}
		y := FontSizeMin + rand.Intn((maxY/2 - FontSizeMin))
		pt.Y = fixed.Int26_6(fontCtx.PointToFixed(float64(y)))
	}

	frames := make([]*image.Paletted, 0, FrameNumber)
	for l := 0; l < FrameNumber; l++ {

		img := image.NewPaletted(image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{ImgW, ImgH},
		}, pallete)

		draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)

		for i := range letters {
			letters[i].font.SetDst(img)
			letters[i].prevPosition.Y += fixed.Int26_6(letters[i].direction) * fixed.Int26_6(letters[i].font.PointToFixed(float64(12)))
			max := fixed.Int26_6(letters[i].font.PointToFixed(float64(ImgH)))
			min := fixed.Int26_6(letters[i].font.PointToFixed(float64(FontSizeMax)))
			clamped, wasClamped := clamp(letters[i].prevPosition.Y, min, max)
			if wasClamped {
				letters[i].direction *= -1
			}
			letters[i].prevPosition.Y = clamped
			pt, _ = letters[i].font.DrawString(string(letters[i].letter), letters[i].prevPosition)
		}

		frames = append(frames, img)
	}

	f, err = os.Create("dapcha.gif")
	if err != nil {
		panic(err)
	}
	delays := make([]int, len(frames))
	for i := range frames {
		delays[i] = 0
	}
	print(len(delays))
	print(len(frames))
	g := gif.GIF{
		Image: frames,
		Delay: delays,
	}
	err = gif.EncodeAll(f, &g)
	if err != nil {
		panic(err)
	}
}

func generateCapcha() string {
	return "Hello, World!"
}

func generateRandomFonts(fontsNumber int, font *truetype.Font, baseIMG *image.Paletted) []*freetype.Context {
	result := make([]*freetype.Context, 0, fontsNumber)
	for i := 0; i < fontsNumber; i++ {
		c := freetype.NewContext()
		c.SetSrc(image.Black)
		c.SetFont(font)
		c.SetDPI(float64(DPI))
		size := FontSizeMin + rand.Intn(FontSizeMax-FontSizeMin)
		c.SetFontSize(float64(size))
		c.SetClip(baseIMG.Bounds())
		c.SetDst(baseIMG)
		result = append(result, c)
	}
	return result
}

func clamp(val, min, max fixed.Int26_6) (fixed.Int26_6, bool) {
	if val < min {
		return min, true
	}

	if val > max {
		return max, true
	}

	return val, false
}

func GIFDecode(w io.Writer, frames []*image.RGBA) {
}
