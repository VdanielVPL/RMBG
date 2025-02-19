package image

import (
	"bytes"
	imgpkg "image"
	_ "image/gif"
	"image/jpeg"
	"image/png"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp"
)

func CropImage(image []byte, left, right, top, bottom float32) ([]byte, string) {

	// println("Cropping image...")
	imgbytes := bytes.NewReader(image)
	img, format, err := imgpkg.Decode(imgbytes)
	if err != nil {
		println(err.Error())
		return nil, format
	}
	// println("Decoded image")

	if format == "jpeg" {
		imgbytes.Seek(0, 0)
		orientation := ReadOrientation(imgbytes)
		// println("Orientation:", orientation)

		switch orientation {
		case 6: // 90 CW
			img = imaging.Rotate270(img)
		case 3: // 180
			img = imaging.Rotate180(img)
		case 8: // 270 CW
			img = imaging.Rotate90(img)
		}
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

	// println("Sub image")

	var buf bytes.Buffer
	if format == "jpeg" {
		err = jpeg.Encode(&buf, croppedImg, &jpeg.Options{Quality: 90})
		if err != nil {
			println(err.Error())
			return nil, format
		}
	} else {
		err = png.Encode(&buf, croppedImg)
		if err != nil {
			println(err.Error())
			return nil, format
		}
	}
	// println("Cropped image")
	return buf.Bytes(), format
}
