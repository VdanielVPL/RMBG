package clipboard

import "syscall"

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	ole32    = syscall.NewLazyDLL("ole32.dll")
	codecs   = syscall.NewLazyDLL("windowscodecs.dll")

	pRegisterClassExW         = user32.NewProc("RegisterClassExW")
	pCreateWindowExW          = user32.NewProc("CreateWindowExW")
	pDefWindowProcW           = user32.NewProc("DefWindowProcW")
	pGetMessageW              = user32.NewProc("GetMessageW")
	pDispatchMessageW         = user32.NewProc("DispatchMessageW")
	pSendMessageW             = user32.NewProc("SendMessageW")
	pDestroyWindow            = user32.NewProc("DestroyWindow")
	pPostQuitMessage          = user32.NewProc("PostQuitMessage")
	pOpenClipboard            = user32.NewProc("OpenClipboard")
	pCloseClipboard           = user32.NewProc("CloseClipboard")
	pEmptyClipboard           = user32.NewProc("EmptyClipboard")
	pSetClipboardData         = user32.NewProc("SetClipboardData")
	pGetClipboardOwner        = user32.NewProc("GetClipboardOwner")
	pRegisterClipboardFormatW = user32.NewProc("RegisterClipboardFormatW")

	pGlobalAlloc      = kernel32.NewProc("GlobalAlloc")
	pGlobalLock       = kernel32.NewProc("GlobalLock")
	pGlobalUnlock     = kernel32.NewProc("GlobalUnlock")
	pGlobalFree       = kernel32.NewProc("GlobalFree")
	pRtlMoveMemory    = kernel32.NewProc("RtlMoveMemory")
	pGetModuleHandleW = kernel32.NewProc("GetModuleHandleW")

	pCoInitializeEx         = ole32.NewProc("CoInitializeEx")
	pCoUninitialize         = ole32.NewProc("CoUninitialize")
	pCoCreateInstance       = ole32.NewProc("CoCreateInstance")
	pWICConvertBitmapSource = codecs.NewProc("WICConvertBitmapSource")
)

const (
	cfDIBV5 = 17

	wmDESTROY          = 0x0002
	wmRENDERFORMAT     = 0x0305
	wmRENDERALLFORMATS = 0x0306
	wmDESTROYCLIPBOARD = 0x0307
	wmCOPY             = 0x8000 + 1 // WM_APP + 1
	wmCLOSEOWNER       = 0x8000 + 2 // WM_APP + 2

	gmemMOVEABLE = 0x0002
	hwndMESSAGE  = ^uintptr(2) // (HWND)-3, message-only window
)

type wndClassEx struct {
	cbSize, style                          uint32
	lpfnWndProc                            uintptr
	cbClsExtra, cbWndExtra                 int32
	hInstance, hIcon, hCursor, hbrBackgrnd uintptr
	lpszMenuName, lpszClassName            *uint16
	hIconSm                                uintptr
}

type winMsg struct {
	hwnd     uintptr
	message  uint32
	wParam   uintptr
	lParam   uintptr
	time     uint32
	pt       struct{ x, y int32 }
	lPrivate uint32
}
