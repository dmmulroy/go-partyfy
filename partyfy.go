package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"

	"github.com/ericpauley/go-quantize/quantize"
)

const (
	PNG  = "image/png"
	JPG  = "image/jpg"
	JPEG = "image/jpeg"
	GIF  = "image/gif"
)

var colors = []color.NRGBA{
	{255, 141, 139, 255},
	{254, 214, 137, 255},
	{136, 255, 137, 255},
	{135, 255, 255, 255},
	{139, 181, 254, 255},
	{215, 140, 255, 255},
	{255, 140, 255, 255},
	{255, 104, 247, 255},
	{254, 108, 183, 255},
	{255, 105, 104, 255},
}

func partyfy(r io.Reader) error {
	buf := bytes.Buffer{}

	if _, err := buf.ReadFrom(r); err != nil {
		return err
	}

	mimetype, err := getFileContentType(buf.Bytes())

	if err != nil {
		return err
	}

	img, err := getImage(&buf, mimetype)

	if err != nil {
		return err
	}

	pixels, err := getPixels(img)

	if err != nil {
		return err
	}

	quant := quantize.MedianCutQuantizer{}
	palette := quant.Quantize(make([]color.Color, 0, 256), img)

	frames := make([]*image.Paletted, len(colors))

	for idx := range frames {
		frames[idx] = &image.Paletted{
			Palette: palette,
			Rect:    img.Bounds(),
			Pix:     make([]uint8, len(pixels)),
		}
	}

	for idx := 0; idx < len(pixels); idx += 4 {
		pixel := color.NRGBA{pixels[idx], pixels[idx+1], pixels[idx+2], pixels[idx+3]}

		gp := mix(grayscale(pixel), colors[0], 60)

		if gp.R != 0x0 {
			fmt.Printf("%#v\n", gp)

		}

		pixels[idx] = gp.R
		pixels[idx+1] = gp.G
		pixels[idx+2] = gp.B
		pixels[idx+3] = gp.A
	}

	x := image.NewNRGBA(img.Bounds())
	x.Pix = pixels

	f, err := os.Create("test.gif")

	if err != nil {
		return err
	}

	defer f.Close()

	err = gif.Encode(f, x, &gif.Options{Quantizer: quantize.MedianCutQuantizer{}})

	if err != nil {
		return err
	}
	return nil
}

func getImage(buf *bytes.Buffer, mimetype string) (image.Image, error) {
	var img image.Image
	var err error

	switch mimetype {
	case PNG:
		img, err = png.Decode(buf)
	case JPEG:
	case JPG:
		img, err = jpeg.Decode(buf)
	case GIF:
		img, err = gif.Decode(buf)
	default:
		err = errors.New("unsupportd mimetype")
	}

	return img, err
}

func getPixels(i image.Image) ([]uint8, error) {
	img, ok := i.(*image.NRGBA)

	if !ok {
		return nil, errors.New("error getting pixels")
	}

	return img.Pix, nil
}

func getFileContentType(buf []byte) (string, error) {
	contentType := http.DetectContentType(buf)

	return contentType, nil
}

func grayscale(pixel color.NRGBA) color.NRGBA {
	gray := uint8(math.Round((0.21*float64(pixel.R) + 0.72*float64(pixel.G) + 0.07*float64(pixel.B))))

	return color.NRGBA{gray, gray, gray, pixel.A}
}

func mix(a color.NRGBA, b color.NRGBA, opacity uint8) color.NRGBA {
	return color.NRGBA{
		uint8((float64(b.R)-float64(a.R))*(float64(opacity)/100) + float64(a.R)),
		uint8((float64(b.G)-float64(a.G))*(float64(opacity)/100) + float64(a.G)),
		uint8((float64(b.B)-float64(a.B))*(float64(opacity)/100) + float64(a.B)),
		a.A,
	}
}
