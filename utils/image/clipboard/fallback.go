package clipboard

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/draw"
	"runtime/debug"

	imgutil "rmbg/utils/image"
)

// decodeGo is the fallback for images WIC can't decode, using the registered Go decoders.
func decodeGo(data []byte, orientation int) uintptr {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return 0
	}

	img = imgutil.ApplyOrientation(img, orientation)

	b := img.Bounds()
	w, h := b.Dx(), b.Dy()

	hMem, pixels := newDIBV5(w, h)
	if hMem == 0 {
		return 0
	}

	writePixels(pixels, img, w, h)
	whitenTransparent(pixels, w, h)
	pGlobalUnlock.Call(hMem)

	go debug.FreeOSMemory() // give the decoded image back to the OS without delaying the paste
	return hMem
}

// writePixels converts row by row (RGBA -> BGRA, bottom-up), without a full-size intermediate copy.
func writePixels(pixels uintptr, img image.Image, w, h int) {
	b := img.Bounds()
	stride := w * 4
	out := make([]byte, stride)

	put := func(y int, src []byte) {
		le := binary.LittleEndian
		for i := 0; i < stride; i += 4 { // swap R and B
			v := le.Uint32(src[i:])
			le.PutUint32(out[i:], v&0xFF00FF00|v<<16&0x00FF0000|v>>16&0x000000FF)
		}
		writeMemory(pixels+uintptr((h-1-y)*stride), out)
	}

	if n, ok := img.(*image.NRGBA); ok { // PNG with alpha: already straight RGBA
		for y := range h {
			put(y, n.Pix[n.PixOffset(b.Min.X, b.Min.Y+y):])
		}
		return
	}

	// Opaque images use RGBA (image/draw fast paths), the rest NRGBA to keep alpha straight.
	var row draw.Image
	var rowPix []byte
	if o, ok := img.(interface{ Opaque() bool }); ok && o.Opaque() {
		r := image.NewRGBA(image.Rect(0, 0, w, 1))
		row, rowPix = r, r.Pix
	} else {
		r := image.NewNRGBA(image.Rect(0, 0, w, 1))
		row, rowPix = r, r.Pix
	}

	for y := range h {
		draw.Draw(row, row.Bounds(), img, image.Pt(b.Min.X, b.Min.Y+y), draw.Src)
		put(y, rowPix)
	}
}
