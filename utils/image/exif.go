package image

import (
	"bufio"
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
	maxBytes := 65536
	bytesRead := 0

	// Searching for EXIF marker (0xFFE1)
	for bytesRead < maxBytes {
		b, err := reader.ReadByte()
		if err != nil {
			return 1
		}
		bytesRead++

		if b == 0xFF {
			nextByte, err := reader.ReadByte()
			if err != nil {
				return 1
			}
			bytesRead++

			switch nextByte {
			case 0xE1:
				println("Znaleziono marker EXIF")
				return processExifData(reader)
			case 0xD8:
				println("Znaleziono początek JPEG")
				continue
			case 0xDA:
				println("Znaleziono początek danych obrazu")
				return 1
			case 0xE0:
				// Skipping JFIF segment
				var skipBytes [2]byte
				if _, err := io.ReadFull(reader, skipBytes[:]); err != nil {
					return 1
				}
				length := int(skipBytes[0])<<8 | int(skipBytes[1])
				if _, err := io.ReadFull(reader, make([]byte, length-2)); err != nil {
					return 1
				}
				bytesRead += length
			}
		}
	}
	return 1
}

func processExifData(reader *bufio.Reader) int {
	// Reading segment length
	var lengthBytes [2]byte
	if _, err := io.ReadFull(reader, lengthBytes[:]); err != nil {
		return 1
	}

	// Checking header "Exif"
	exifHeader := make([]byte, 4)
	if _, err := io.ReadFull(reader, exifHeader); err != nil {
		return 1
	}
	if string(exifHeader) != "Exif" {
		return 1
	}

	// Skip 2 zeros
	if _, err := io.ReadFull(reader, make([]byte, 2)); err != nil {
		return 1
	}

	// Checking format (II or MM)
	var endianBytes [2]byte
	if _, err := io.ReadFull(reader, endianBytes[:]); err != nil {
		return 1
	}

	var littleEndian bool
	if string(endianBytes[:]) == "II" {
		littleEndian = true
	} else if string(endianBytes[:]) == "MM" {
		littleEndian = false
	} else {
		return 1
	}

	// Checking tag TIFF
	var tagMark [2]byte
	if _, err := io.ReadFull(reader, tagMark[:]); err != nil {
		return 1
	}
	if getUint16(tagMark[:], littleEndian) != 0x002A {
		return 1
	}

	// Reading offset IFD0
	var offsetBytes [4]byte
	if _, err := io.ReadFull(reader, offsetBytes[:]); err != nil {
		return 1
	}
	ifdOffset := getUint32(offsetBytes[:], littleEndian)

	// Skiping to IFD0
	if _, err := io.ReadFull(reader, make([]byte, ifdOffset-8)); err != nil {
		return 1
	}

	// Reading number of entries in IFD0
	var entriesBytes [2]byte
	if _, err := io.ReadFull(reader, entriesBytes[:]); err != nil {
		return 1
	}
	entries := getUint16(entriesBytes[:], littleEndian)

	// Searching for orientation tag (0x0112)
	for i := 0; i < int(entries); i++ {
		var entry [12]byte
		if _, err := io.ReadFull(reader, entry[:]); err != nil {
			return 1
		}

		tag := getUint16(entry[0:2], littleEndian)
		if tag == 0x0112 {
			return int(getUint16(entry[8:10], littleEndian))
		}
	}

	return 1
}
