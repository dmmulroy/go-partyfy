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
	"net/http"
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

	for idx := 0; idx < len(pixels); idx += 4 {
		pixel := color.RGBA{pixels[idx], pixels[idx+1], pixels[idx+2], pixels[idx+3]}

		fmt.Println(pixel)
	}

	// quant := quantize.MedianCutQuantizer{}

	// palette := quant.Quantize(make([]color.Color, 0, 256), img)

	// frames := make([]image.Paletted, len(colors))

	// fmt.Println(NRGBAimg.Pix)

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
