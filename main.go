package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mandykoh/prism/linear"
	"github.com/mandykoh/prism/srgb"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var dpiFlag = flag.Float64("dpi", 92, "pixel per inch resolution")
var dstFlag = flag.String("dst", ".", "path where to store the output png images")
var fontFlag = flag.String("ttf", "", "ttf font file to use (default: goregular)")
var textFlag = flag.String("text", "This is a test", "text to render")
var ptFlag = flag.Float64("pt", 16, "pt size of font")
var srgbFlag = flag.Bool("srgb", false, "use sRGB gamma correction")
var gammaFlag = flag.Float64("gamma", 1., "correct for specified gamma value")

func main() {
	flag.Parse()
	// load the font
	ttfFont := goregular.TTF
	if *fontFlag != "" {
		var err error
		ttfFont, err = ioutil.ReadFile(*fontFlag)
		if err != nil {
			log.Fatal(err)
		}
	}
	// parse the font
	f, err := opentype.Parse(ttfFont)
	if err != nil {
		log.Fatalf("font parse: %v", err)
	}
	// build the font face (collection of glyphs for specified size and DPI)
	dpi := *dpiFlag
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    *ptFlag,
		DPI:     float64(dpi),
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("create face: %v", err)
	}
	// compute image dimension
	text := *textFlag

	margin := int(math.Ceil(*ptFlag * dpi / 72.))
	metrics := face.Metrics()
	height := (int(metrics.Height)+63)/64 + 2*margin // int(Ceil)
	startingDotX := fixed.I(margin)
	startingDotY := (metrics.Ascent+63)&^63 + fixed.I(margin) // Ceil
	width := (int(font.MeasureString(face, text))+63)/64 + 2*margin
	// Create an image with 16bit gray colors of specified size filled with white
	dst := image.NewGray16(image.Rect(0, 0, width, height))
	white := color.Gray{255}
	draw.Draw(dst, dst.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)
	// Instantiate a font face drawer and draw text
	d := font.Drawer{
		Dst:  dst,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.Point26_6{X: startingDotX, Y: startingDotY},
	}
	d.DrawString(text)
	// Convert the image to 8bit gray using the specified gamma correction
	out := image.NewGray(dst.Bounds())
	var gammaStr string
	if *srgbFlag {
		gammaStr = "srgb"
		srgb.EncodeImage(out, dst, runtime.NumCPU())
	} else if *gammaFlag == 1. {
		gammaStr = "gamma_1"
		draw.Draw(out, out.Rect, dst, dst.Rect.Min, draw.Src)
	} else {
		gammaStr = fmt.Sprintf("gamma_%g", *gammaFlag)
		gammaEncodeImage(out, dst, *gammaFlag, runtime.NumCPU())
	}
	// get font name
	ttfFontName := "goregular"
	if *fontFlag != "" {
		ttfFontName = filepath.Base(*fontFlag)
		if pos := strings.LastIndexByte(ttfFontName, '.'); pos != -1 {
			ttfFontName = ttfFontName[:pos]
		}
	}
	// save image to png file with some parameters in the name
	fileName := fmt.Sprintf("output_%gpt_%gdpi_%s_%s.png", *ptFlag, dpi, gammaStr, ttfFontName)
	if *dstFlag != "." {
		fileName = filepath.Join(*dstFlag, fileName)
	}
	if err := writeImage(fileName, out); err != nil {
		log.Fatal("save png:", err)
	}
}

func writeImage(path string, img image.Image) error {
	_ = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	imgFile, err := os.Create(path)
	if err != nil {
		return nil
	}
	defer imgFile.Close()
	return png.Encode(imgFile, img)
}

func gammaEncodeImage(dst draw.Image, src image.Image, gamma float64, parallelism int) {
	if gamma == 0. {
		log.Fatal("invalid gamma value 0.")
	}
	cor := 1 / gamma
	linear.TransformImageColor(dst, src, parallelism, func(c color.Color) color.RGBA64 {
		rgb, alpha := linear.RGBFromLinear(c)
		return color.RGBA64{
			A: uint16(alpha * 65535),
			R: uint16(math.Pow(float64(rgb.R), cor) * 65535),
			G: uint16(math.Pow(float64(rgb.G), cor) * 65535),
			B: uint16(math.Pow(float64(rgb.B), cor) * 65535),
		}
	})
}
