package image

import (
	"syscall"
	"unsafe"
)

func CopyToClipboard(imageData []byte) error {

	if imageData != nil {

		err := setClipboardImage(imageData)
		if err != nil {
			return err
		}
		return nil

	}
	return nil
}

// Forgive me for what you are about to see...
var CF_PNG uint32

func init() {
	formatName, _ := syscall.UTF16PtrFromString("PNG")
	CF_PNG = RegisterClipboardFormat(formatName)
}

func RegisterClipboardFormat(formatName *uint16) uint32 {
	user32 := syscall.NewLazyDLL("user32.dll")
	registerClipboardFormat := user32.NewProc("RegisterClipboardFormatW")
	ret, _, _ := registerClipboardFormat.Call(uintptr(unsafe.Pointer(formatName)))
	return uint32(ret)
}

func setClipboardImage(pngData []byte) error {
	user32 := syscall.NewLazyDLL("user32.dll")
	kernel32 := syscall.NewLazyDLL("kernel32.dll")

	openClipboard := user32.NewProc("OpenClipboard")
	emptyClipboard := user32.NewProc("EmptyClipboard")
	setClipboardData := user32.NewProc("SetClipboardData")
	closeClipboard := user32.NewProc("CloseClipboard")

	globalAlloc := kernel32.NewProc("GlobalAlloc")
	globalLock := kernel32.NewProc("GlobalLock")
	globalUnlock := kernel32.NewProc("GlobalUnlock")

	const GMEM_MOVEABLE = 0x0002

	// Opening clipboard
	ret, _, err := openClipboard.Call(0)
	if ret == 0 {
		return err
	}
	defer closeClipboard.Call()

	// Emptying clipboard
	ret, _, err = emptyClipboard.Call()
	if ret == 0 {
		return err
	}

	// Allocating global mem
	hMem, _, err := globalAlloc.Call(GMEM_MOVEABLE, uintptr(len(pngData)))
	if hMem == 0 {
		return err
	}

	// Locking mem and copying data
	memPtr, _, err := globalLock.Call(hMem)
	if memPtr == 0 {
		return err
	}
	defer globalUnlock.Call(hMem)

	libc := syscall.NewLazyDLL("msvcrt.dll")
	memcpy := libc.NewProc("memcpy")
	memcpy.Call(memPtr, uintptr(unsafe.Pointer(&pngData[0])), uintptr(len(pngData)))

	// Setting clipboard data
	ret, _, err = setClipboardData.Call(uintptr(CF_PNG), hMem)
	if ret == 0 {
		return err
	}

	return nil
}
