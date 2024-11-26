package image

import (
	"bytes"
	imgpkg "image"
	_ "image/gif"
	"image/jpeg"
	"image/png"

	_ "golang.org/x/image/webp"
)

func CropImage(image []byte, left, right, top, bottom float32) []byte {

	println("Cropping image...")
	img, format, err := imgpkg.Decode(bytes.NewReader(image))
	if err != nil {
		println(err.Error())
		return nil
	}
	println("Decoded image")

	bounds := img.Bounds()

	x0 := bounds.Min.X + int(left)
	y0 := bounds.Min.Y + int(top)
	x1 := bounds.Max.X - int(right)
	y1 := bounds.Max.Y - int(bottom)

	rect := imgpkg.Rect(x0, y0, x1, y1)

	croppedImg := img.(interface {
		SubImage(r imgpkg.Rectangle) imgpkg.Image
	}).SubImage(rect)

	println("Sub image")

	var buf bytes.Buffer
	if format == "png" {
		err = png.Encode(&buf, croppedImg)
		if err != nil {
			println(err.Error())
			return nil
		}
	} else if format == "jpeg" {
		err = jpeg.Encode(&buf, croppedImg, &jpeg.Options{Quality: 90})
		if err != nil {
			println(err.Error())
			return nil
		}
	} else {
		println("Unsupported format: ", format)
		return nil
	}
	println("Cropped image")
	return buf.Bytes()
}
