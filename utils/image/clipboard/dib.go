package clipboard

import (
	"bytes"
	"encoding/binary"

	imgutil "rmbg/utils/image"
)

const dibHeaderSize = 124 // BITMAPV5HEADER

var jpegMagic = []byte{0xFF, 0xD8}

// renderDIBV5 decodes the image into clipboard memory, with WIC or the Go decoders as a fallback.
func renderDIBV5(data []byte) uintptr {
	orientation := 1
	if bytes.HasPrefix(data, jpegMagic) { // apply EXIF orientation like the preview does
		orientation = imgutil.ReadOrientation(bytes.NewReader(data))
	}

	if hMem := decodeWIC(data, orientation); hMem != 0 {
		return hMem
	}
	return decodeGo(data, orientation) // e.g. WebP without the Windows codec installed
}

// newDIBV5 allocates a 32-bit bottom-up BGRA DIBV5 with straight (not premultiplied) alpha
// and returns its handle (still locked) and the address of the pixel area.
func newDIBV5(w, h int) (hMem, pixels uintptr) {
	if w <= 0 || h <= 0 {
		return 0, 0
	}

	hMem, ptr := globalAlloc(dibHeaderSize + w*h*4)
	if hMem == 0 {
		return 0, 0
	}

	var hdr [dibHeaderSize]byte
	le := binary.LittleEndian
	le.PutUint32(hdr[0:], dibHeaderSize)
	le.PutUint32(hdr[4:], uint32(w))
	le.PutUint32(hdr[8:], uint32(h)) // positive = bottom-up; Word can't paste negative heights
	le.PutUint16(hdr[12:], 1)        // planes
	le.PutUint16(hdr[14:], 32)       // bpp
	le.PutUint32(hdr[16:], 3)        // BI_BITFIELDS
	le.PutUint32(hdr[20:], uint32(w*h*4))
	le.PutUint32(hdr[40:], 0x00FF0000) // R mask
	le.PutUint32(hdr[44:], 0x0000FF00) // G mask
	le.PutUint32(hdr[48:], 0x000000FF) // B mask
	le.PutUint32(hdr[52:], 0xFF000000) // A mask
	le.PutUint32(hdr[56:], 0x73524742) // LCS_sRGB
	le.PutUint32(hdr[108:], 4)         // LCS_GM_IMAGES
	writeMemory(ptr, hdr[:])

	return hMem, ptr + dibHeaderSize
}

// whitenTransparent makes fully transparent pixels white instead of (usually) black,
// so apps that ignore alpha (e.g. Paint) paste a white background. Alpha-aware apps don't see the change.
func whitenTransparent(pixels uintptr, w, h int) {
	stride := w * 4
	row := make([]byte, stride)
	for y := range h {
		p := pixels + uintptr(y*stride)
		readMemory(row, p)
		changed := false
		for i := 0; i < stride; i += 4 {
			if row[i+3] == 0 {
				row[i], row[i+1], row[i+2] = 0xFF, 0xFF, 0xFF
				changed = true
			}
		}
		if changed {
			writeMemory(p, row)
		}
	}
}
