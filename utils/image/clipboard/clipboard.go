// Package clipboard copies images to the Windows clipboard using delayed rendering:
// copying only promises "PNG" and CF_DIBV5, the data is produced when an app pastes.
// "PNG" gives the original bytes (with alpha), CF_DIBV5 is decoded on demand by WIC
// and Windows synthesizes CF_DIB/CF_BITMAP from it (Chromium/Electron need CF_DIB).
package clipboard

import (
	"errors"
	"sync"
	"syscall"
	"unsafe"
)

var (
	ErrNotReady = errors.New("clipboard not initialized")
	ErrBusy     = errors.New("clipboard is used by another application")
)

var (
	initOnce sync.Once
	callMu   sync.Mutex // serializes CopyImage/Close
	hwnd     uintptr
	cfPNG    uintptr
	done     chan struct{}
)

// Init starts the hidden clipboard owner window. Call once at startup.
func Init() error {
	var err error

	initOnce.Do(func() {
		name, _ := syscall.UTF16PtrFromString("PNG")
		cfPNG, _, _ = pRegisterClipboardFormatW.Call(uintptr(unsafe.Pointer(name)))

		ready := make(chan error)
		done = make(chan struct{})
		go ownerThread(ready)
		err = <-ready
	})

	if hwnd == 0 && err == nil {
		err = ErrNotReady
	}
	return err
}

// CopyImage puts a PNG, JPEG or WebP image on the clipboard.
// data is kept, not copied, so it must not be modified afterwards.
func CopyImage(data []byte) error {
	if len(data) == 0 {
		return errors.New("no image data")
	}

	callMu.Lock()
	defer callMu.Unlock()

	if hwnd == 0 {
		return ErrNotReady
	}

	pending = data
	pSendMessageW.Call(hwnd, wmCOPY, 0, 0)
	return copyErr
}

// Close destroys the owner window. An image still on the clipboard is rendered
// before that, so it stays pasteable after the app exits. Call on shutdown.
func Close() {
	callMu.Lock()
	defer callMu.Unlock()

	if hwnd == 0 {
		return
	}

	pSendMessageW.Call(hwnd, wmCLOSEOWNER, 0, 0)
	<-done
	hwnd = 0
}
