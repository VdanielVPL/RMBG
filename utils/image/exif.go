package image

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)

func ReadOrientation(image io.Reader) int {
	reader := bufio.NewReader(image)

	// Checking start of JPEG image (0xFFD8)
	var soi [2]byte
	if _, err := io.ReadFull(reader, soi[:]); err != nil || soi[0] != 0xFF || soi[1] != 0xD8 {
		return 1
	}

	// Going through segments until EXIF or start of image data
	for {
		b, err := reader.ReadByte()
		if err != nil || b != 0xFF {
			return 1
		}
		marker := byte(0xFF)
		for marker == 0xFF { // markers can be padded with extra 0xFF bytes
			if marker, err = reader.ReadByte(); err != nil {
				return 1
			}
		}
		if marker == 0xDA || marker == 0xD9 {
			return 1
		}

		var lengthBytes [2]byte
		if _, err := io.ReadFull(reader, lengthBytes[:]); err != nil {
			return 1
		}
		length := int(lengthBytes[0])<<8 | int(lengthBytes[1])
		if length < 2 {
			return 1
		}

		// APP1 is also used by XMP, only the one starting with "Exif" is EXIF
		if marker == 0xE1 {
			segment := make([]byte, length-2)
			if _, err := io.ReadFull(reader, segment); err != nil {
				return 1
			}
			if bytes.HasPrefix(segment, []byte("Exif\x00\x00")) {
				return processExifData(segment[6:])
			}
			continue
		}

		if _, err := reader.Discard(length - 2); err != nil {
			return 1
		}
	}
}

func processExifData(tiff []byte) int {
	if len(tiff) < 8 {
		return 1
	}

	// Checking format (II or MM)
	var order binary.ByteOrder
	switch string(tiff[0:2]) {
	case "II":
		order = binary.LittleEndian
	case "MM":
		order = binary.BigEndian
	default:
		return 1
	}

	// Checking tag TIFF
	if order.Uint16(tiff[2:4]) != 0x002A {
		return 1
	}

	// Reading offset IFD0 (counted from the start of the TIFF header)
	ifdOffset := int64(order.Uint32(tiff[4:8]))
	if ifdOffset < 8 || ifdOffset+2 > int64(len(tiff)) {
		return 1
	}

	// Reading number of entries in IFD0
	entries := int64(order.Uint16(tiff[ifdOffset:]))

	// Searching for orientation tag (0x0112)
	for i := int64(0); i < entries; i++ {
		start := ifdOffset + 2 + i*12
		if start+12 > int64(len(tiff)) {
			return 1
		}
		entry := tiff[start : start+12]

		if order.Uint16(entry[0:2]) == 0x0112 {
			orientation := int(order.Uint16(entry[8:10]))
			if orientation < 1 || orientation > 8 {
				return 1
			}
			return orientation
		}
	}

	return 1
}
