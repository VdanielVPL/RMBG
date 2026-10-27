package clipboard

import (
	"runtime"
	"syscall"
	"unsafe"
)

// Windows Imaging Component: the native decoder, decodes straight into clipboard memory.

type guid struct {
	d1     uint32
	d2, d3 uint16
	d4     [8]byte
}

var (
	clsidWICImagingFactory  = guid{0xcacaf262, 0x9370, 0x4615, [8]byte{0xa1, 0x3b, 0x9f, 0x55, 0x39, 0xda, 0x4c, 0x0a}}
	iidIWICImagingFactory   = guid{0xec5ec8a9, 0xc395, 0x4314, [8]byte{0x9c, 0x77, 0x54, 0xd7, 0xa9, 0x35, 0xff, 0x70}}
	wicPixelFormat32bppBGRA = guid{0x6fddc324, 0x4e03, 0x4bfe, [8]byte{0xb1, 0x85, 0x3d, 0x77, 0x76, 0x8d, 0xc9, 0x0f}} // straight alpha
)

// COM vtable indexes of the methods we call.
const (
	vtRelease = 2 // IUnknown

	vtFactoryCreateDecoderFromStream = 4  // IWICImagingFactory
	vtFactoryCreateFlipRotator       = 13 // IWICImagingFactory
	vtFactoryCreateStream            = 14 // IWICImagingFactory
	vtStreamInitializeFromMemory     = 16 // IWICStream
	vtDecoderGetFrame                = 13 // IWICBitmapDecoder
	vtSourceGetSize                  = 3  // IWICBitmapSource
	vtSourceCopyPixels               = 7  // IWICBitmapSource
	vtFlipRotatorInitialize          = 8  // IWICBitmapFlipRotator
)

// WICBitmapTransformOptions
const (
	wicRotate90       = 1
	wicRotate180      = 2
	wicRotate270      = 3
	wicFlipHorizontal = 0x8
	wicFlipVertical   = 0x10
)

// comObject is the memory layout of any COM interface: a pointer to its method table.
type comObject struct {
	vtbl *[32]uintptr
}

var wicFactory *comObject // created on first use, owner thread only

// call calls method idx of the object and returns the HRESULT (0 = S_OK).
// uintptrescapes keeps out-variables passed as uintptr(unsafe.Pointer(&x)) on the heap (like LazyProc.Call),
// otherwise the stack can move during the call and COM writes the result to the old place.
//
//go:uintptrescapes
func (o *comObject) call(idx int, args ...uintptr) uintptr {
	r, _, _ := syscall.SyscallN(o.vtbl[idx], append([]uintptr{uintptr(unsafe.Pointer(o))}, args...)...)
	return r
}

func (o *comObject) release() {
	o.call(vtRelease)
}

func getWICFactory() *comObject {
	if wicFactory == nil {
		hr, _, _ := pCoCreateInstance.Call(
			uintptr(unsafe.Pointer(&clsidWICImagingFactory)), 0, 1, // CLSCTX_INPROC_SERVER
			uintptr(unsafe.Pointer(&iidIWICImagingFactory)), uintptr(unsafe.Pointer(&wicFactory)))
		if hr != 0 {
			wicFactory = nil
		}
	}
	return wicFactory
}

func releaseWIC() {
	if wicFactory != nil {
		wicFactory.release()
		wicFactory = nil
	}
}

// wicTransform maps EXIF orientation to the WIC transform that also turns rows bottom-up for the DIB.
// WIC flips before rotating, so 90 and 270 are swapped compared to plain rotation.
func wicTransform(orientation int) uintptr {
	switch orientation {
	case 2:
		return wicFlipHorizontal | wicFlipVertical
	case 3:
		return wicRotate180 | wicFlipVertical
	case 4:
		return wicRotate180 | wicFlipHorizontal | wicFlipVertical
	case 5:
		return wicRotate90 | wicFlipHorizontal | wicFlipVertical
	case 6:
		return wicRotate270 | wicFlipVertical
	case 7:
		return wicRotate270 | wicFlipHorizontal | wicFlipVertical
	case 8:
		return wicRotate90 | wicFlipVertical
	}
	return wicFlipVertical
}

func decodeWIC(data []byte, orientation int) uintptr {
	factory := getWICFactory()
	if factory == nil {
		return 0
	}
	defer runtime.KeepAlive(data) // WIC reads data in place until the stream is released

	var stream, decoder, frame, converter, transformer *comObject

	if factory.call(vtFactoryCreateStream, uintptr(unsafe.Pointer(&stream))) != 0 {
		return 0
	}
	defer stream.release()

	if stream.call(vtStreamInitializeFromMemory, uintptr(unsafe.Pointer(&data[0])), uintptr(len(data))) != 0 {
		return 0
	}

	if factory.call(vtFactoryCreateDecoderFromStream, uintptr(unsafe.Pointer(stream)), 0, 0, uintptr(unsafe.Pointer(&decoder))) != 0 {
		return 0
	}
	defer decoder.release()

	if decoder.call(vtDecoderGetFrame, 0, uintptr(unsafe.Pointer(&frame))) != 0 {
		return 0
	}
	defer frame.release()

	if hr, _, _ := pWICConvertBitmapSource.Call(uintptr(unsafe.Pointer(&wicPixelFormat32bppBGRA)), uintptr(unsafe.Pointer(frame)), uintptr(unsafe.Pointer(&converter))); hr != 0 {
		return 0
	}
	defer converter.release()

	if factory.call(vtFactoryCreateFlipRotator, uintptr(unsafe.Pointer(&transformer))) != 0 {
		return 0
	}
	defer transformer.release()

	if transformer.call(vtFlipRotatorInitialize, uintptr(unsafe.Pointer(converter)), wicTransform(orientation)) != 0 {
		return 0
	}

	var w, h uint32
	if transformer.call(vtSourceGetSize, uintptr(unsafe.Pointer(&w)), uintptr(unsafe.Pointer(&h))) != 0 {
		return 0
	}

	hMem, pixels := newDIBV5(int(w), int(h))
	if hMem == 0 {
		return 0
	}

	// Whole image in one call: row by row makes some sources (rotator, some codecs)
	// decode the whole image again for every row.
	stride := uintptr(w) * 4
	if transformer.call(vtSourceCopyPixels, 0, stride, stride*uintptr(h), pixels) != 0 {
		pGlobalUnlock.Call(hMem)
		pGlobalFree.Call(hMem)
		return 0
	}
	whitenTransparent(pixels, int(w), int(h))

	pGlobalUnlock.Call(hMem)
	return hMem
}
