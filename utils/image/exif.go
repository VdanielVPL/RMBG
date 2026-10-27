package image

import (
	"bufio"
	"bytes"
	"io"
)

func getUint16(data []byte, littleEndian bool) uint16 {
	if littleEndian {
		return uint16(data[0]) | uint16(data[1])<<8
	}
	return uint16(data[1]) | uint16(data[0])<<8
}

func getUint32(data []byte, littleEndian bool) uint32 {
	if littleEndian {
		return uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24
	}
	return uint32(data[3]) | uint32(data[2])<<8 | uint32(data[1])<<16 | uint32(data[0])<<24
}

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
	var littleEndian bool
	if string(tiff[0:2]) == "II" {
		littleEndian = true
	} else if string(tiff[0:2]) == "MM" {
		littleEndian = false
	} else {
		return 1
	}

	// Checking tag TIFF
	if getUint16(tiff[2:4], littleEndian) != 0x002A {
		return 1
	}

	// Reading offset IFD0 (counted from the start of the TIFF header)
	ifdOffset := int64(getUint32(tiff[4:8], littleEndian))
	if ifdOffset < 8 || ifdOffset+2 > int64(len(tiff)) {
		return 1
	}

	// Reading number of entries in IFD0
	entries := int64(getUint16(tiff[ifdOffset:], littleEndian))

	// Searching for orientation tag (0x0112)
	for i := int64(0); i < entries; i++ {
		start := ifdOffset + 2 + i*12
		if start+12 > int64(len(tiff)) {
			return 1
		}
		entry := tiff[start : start+12]

		if getUint16(entry[0:2], littleEndian) == 0x0112 {
			orientation := int(getUint16(entry[8:10], littleEndian))
			if orientation < 1 || orientation > 8 {
				return 1
			}
			return orientation
		}
	}

	return 1
}
