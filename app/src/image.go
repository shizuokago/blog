package blog

import (
	"bytes"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"

	"io"
	"io/ioutil"
	"math"

	"github.com/dddaisuke/graphics-go/graphics"
	//"github.com/dddaisuke/graphics-go/graphics/interp"
	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
)

var affines map[int]graphics.Affine = map[int]graphics.Affine{
	1: graphics.I,
	2: graphics.I.Scale(-1, 1),
	3: graphics.I.Scale(-1, -1),
	4: graphics.I.Scale(1, -1),
	5: graphics.I.Rotate(toRadian(90)).Scale(-1, 1),
	6: graphics.I.Rotate(toRadian(90)),
	7: graphics.I.Rotate(toRadian(-90)).Scale(-1, 1),
	8: graphics.I.Rotate(toRadian(-90)),
}

func toRadian(d float64) float64 {
	return math.Pi * d / 180
}

func init() {
}

func convertImage(r io.Reader) ([]byte, bool, error) {

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, false, err
	}

	var img image.Image

	// get orientation
	buff := bytes.NewBuffer(b)
	//o, err := readOrientation(buff)

	/*
		o := 1
		if err == nil {

			img, _, err = image.Decode(buff)
			if err != nil {
				return nil, err
			}
			bounds := img.Bounds()

			if o >= 5 && o <= 8 {
				s := bounds.Size()
				bounds = image.Rectangle{bounds.Min, image.Point{s.Y, s.X}}
			}

			d := image.NewRGBA64(bounds)
			affine := affines[o]
			affine.TransformCenter(d, bounds, interp.Bilinear)

			img = d
		}
	*/

	cnv := false

	//over 1mb
	if len(b) > (1 * 1024 * 1024) {
		if img == nil {
			img, _, err = image.Decode(buff)
			if err != nil {
				return nil, false, err
			}
		}
		img = resize.Resize(1000, 0, img, resize.Lanczos3)
		cnv = true
	}

	if cnv {
		buffer := new(bytes.Buffer)
		if err := jpeg.Encode(buffer, img, nil); err != nil {
			return nil, cnv, err
		}
		b = buffer.Bytes()
	}

	return b, cnv, nil
}

func readOrientation(r io.Reader) (o int, err error) {
	e, err := exif.Decode(r)
	if err != nil {
		return 0, err
	}
	tag, err := e.Get(exif.Orientation)
	if err != nil {
		return 0, err
	}
	o, err = tag.Int(0)
	if err != nil {
		return 0, err
	}
	return o, nil
}
