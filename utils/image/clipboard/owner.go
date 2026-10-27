package clipboard

import (
	"bytes"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

var pngSignature = []byte("\x89PNG\r\n\x1a\n")

// Touched only on the owner thread, or while callMu is held and the owner thread
// is blocked in SendMessage, so no extra locking is needed.
var (
	pending []byte // image passed from CopyImage to the owner thread
	current []byte // image currently owned on the clipboard
	copyErr error
)

// ownerThread runs the message-only window that owns the clipboard.
// Clipboard ownership and window messages are bound to this OS thread.
func ownerThread(ready chan<- error) {
	runtime.LockOSThread()
	defer close(done)

	pCoInitializeEx.Call(0, 0x2) // COINIT_APARTMENTTHREADED, for WIC
	defer pCoUninitialize.Call()
	defer releaseWIC()

	hinst, _, _ := pGetModuleHandleW.Call(0)
	cls, _ := syscall.UTF16PtrFromString("RMBGClipboardOwner")

	wc := wndClassEx{lpfnWndProc: syscall.NewCallback(wndProc), hInstance: hinst, lpszClassName: cls}
	wc.cbSize = uint32(unsafe.Sizeof(wc))
	if r, _, err := pRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc))); r == 0 {
		ready <- err
		return
	}

	h, _, err := pCreateWindowExW.Call(0, uintptr(unsafe.Pointer(cls)), 0, 0, 0, 0, 0, 0, hwndMESSAGE, 0, hinst, 0)
	if h == 0 {
		ready <- err
		return
	}
	hwnd = h
	ready <- nil

	var m winMsg
	for {
		r, _, _ := pGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(r) <= 0 {
			break
		}
		pDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
}

func wndProc(h, msg, wp, lp uintptr) uintptr {
	switch msg {
	case wmCOPY:
		copyErr = setClipboard(h, pending)
		pending = nil
		return 0

	case wmRENDERFORMAT: // an app pastes, the clipboard is already opened by it
		render(wp)
		return 0

	case wmRENDERALLFORMATS: // we are closing while still owning the clipboard
		if r, _, _ := pOpenClipboard.Call(h); r != 0 {
			if owner, _, _ := pGetClipboardOwner.Call(); owner == h {
				render(cfPNG)
				render(cfDIBV5)
			}
			pCloseClipboard.Call()
		}
		return 0

	case wmDESTROYCLIPBOARD:
		current = nil
		return 0

	case wmCLOSEOWNER:
		pDestroyWindow.Call(h)
		return 0

	case wmDESTROY:
		pPostQuitMessage.Call(0)
		return 0
	}

	r, _, _ := pDefWindowProcW.Call(h, msg, wp, lp)
	return r
}

// setClipboard takes ownership of the clipboard and promises the formats without rendering them.
func setClipboard(h uintptr, data []byte) error {
	if !openClipboard(h) {
		return ErrBusy
	}
	defer pCloseClipboard.Call()

	if r, _, err := pEmptyClipboard.Call(); r == 0 {
		return err
	}
	current = data

	// Many apps take the first format they understand, so lossless "PNG" goes first.
	formats := []uintptr{cfPNG, cfDIBV5}
	if !bytes.HasPrefix(data, pngSignature) {
		formats = formats[1:]
	}

	// With a NULL handle SetClipboardData returns NULL even on success, only the error code counts.
	for _, f := range formats {
		if _, _, err := pSetClipboardData.Call(f, 0); err.(syscall.Errno) != 0 {
			current = nil
			return err
		}
	}
	return nil
}

// openClipboard retries for a moment, another app may be holding the clipboard.
func openClipboard(h uintptr) bool {
	for range 20 {
		if r, _, _ := pOpenClipboard.Call(h); r != 0 {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

// render produces a promised format from the current image.
func render(format uintptr) {
	if current == nil {
		return
	}

	var hMem uintptr
	switch format {
	case cfPNG:
		hMem = globalFromBytes(current)
	case cfDIBV5:
		hMem = renderDIBV5(current)
	}
	if hMem == 0 {
		return
	}

	if r, _, _ := pSetClipboardData.Call(format, hMem); r == 0 {
		pGlobalFree.Call(hMem) // on success the system owns the memory
	}
}
