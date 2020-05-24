package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
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
