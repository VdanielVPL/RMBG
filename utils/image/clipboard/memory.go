package clipboard

import "unsafe"

// globalAlloc returns a moveable global memory block and the address of its locked contents.
func globalAlloc(size int) (hMem, ptr uintptr) {
	hMem, _, _ = pGlobalAlloc.Call(gmemMOVEABLE, uintptr(size))
	if hMem == 0 {
		return 0, 0
	}

	ptr, _, _ = pGlobalLock.Call(hMem)
	if ptr == 0 {
		pGlobalFree.Call(hMem)
		return 0, 0
	}
	return hMem, ptr
}

func globalFromBytes(b []byte) uintptr {
	hMem, ptr := globalAlloc(len(b))
	if hMem == 0 {
		return 0
	}

	writeMemory(ptr, b)
	pGlobalUnlock.Call(hMem)
	return hMem
}

// readMemory copies memory that is not managed by Go to b.
func readMemory(b []byte, src uintptr) {
	pRtlMoveMemory.Call(uintptr(unsafe.Pointer(&b[0])), src, uintptr(len(b)))
}

// writeMemory copies b to memory that is not managed by Go.
func writeMemory(dst uintptr, b []byte) {
	pRtlMoveMemory.Call(dst, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)))
}
