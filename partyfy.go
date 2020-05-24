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

func partyfy(r io.Reader) error {
	buf := bytes.Buffer{}

	if _, err := buf.ReadFrom(r); err != nil {
		return err
	}

	mimetype, err := getFileContentType(buf.Bytes())

	if err != nil {
		return err
	}

	img, err := getNRGBAImage(&buf, mimetype)

	if err != nil {
		return err
	}

	fmt.Println(img)

	return nil
}

func getNRGBAImage(buf *bytes.Buffer, mimetype string) (*image.NRGBA, error) {
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

	imgNRGBA, ok := img.(*image.NRGBA)

	if !ok {
		err = errors.New("error converting to NRGBA image")
	}

	return imgNRGBA, err
}

func getFileContentType(buf []byte) (string, error) {
	contentType := http.DetectContentType(buf)

	return contentType, nil
}

func getPartyColors() []color.NRGBA {
	return []color.NRGBA{
		{255, 141, 139, 255},
		{254, 214, 137, 255},
		{136, 255, 137, 255},
		{135, 255, 255, 255},
		{139, 181, 254, 255},
		{215, 140, 255, 255},
		{255, 140, 255, 255},
		{255, 104, 247, 255},
		{254, 108, 183, 255},
		{255, 105, 104, 255}}
}
