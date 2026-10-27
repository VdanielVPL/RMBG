package image

import (
	"bytes"
	imgpkg "image"
	"image/jpeg"
	"image/png"

	_ "golang.org/x/image/webp"
)

func CropImage(image []byte, left, right, top, bottom float32) ([]byte, string, error) {

	img, format, err := decodeOriented(image)
	if err != nil {
		return nil, format, err
	}

	bounds := img.Bounds()

	x0 := bounds.Min.X + int(left)
	y0 := bounds.Min.Y + int(top)
	x1 := bounds.Max.X - int(right)
	y1 := bounds.Max.Y - int(bottom)

	rect := imgpkg.Rect(x0, y0, x1, y1)

	croppedImg := img.(interface {
		SubImage(r imgpkg.Rectangle) imgpkg.Image
	}).SubImage(rect)

	var buf bytes.Buffer
	if format == "jpeg" {
		err = jpeg.Encode(&buf, croppedImg, &jpeg.Options{Quality: 100})
		if err != nil {
			return nil, format, err
		}
	} else {
		err = png.Encode(&buf, croppedImg)
		if err != nil {
			return nil, format, err
		}
	}
	return buf.Bytes(), format, nil
}
