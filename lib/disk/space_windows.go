// +build windows

package disk

import (
	"syscall"
	"unsafe"
)

// Space returns total and free bytes available in a directory, e.g. `C:\`.
// It returns free space available to the user (including quota limitations),
// so it can be lower than the free space of the disk.
func Space(path string) (total, free uint64, err error) {
	kernel32, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return 0, 0, err
	}
	defer syscall.FreeLibrary(kernel32)
	GetDiskFreeSpaceEx, err := syscall.GetProcAddress(syscall.Handle(kernel32), "GetDiskFreeSpaceExW")
	if err != nil {
		return 0, 0, err
	}
	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)
	r1, _, e1 := syscall.Syscall6(uintptr(GetDiskFreeSpaceEx), 4,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)), 0, 0)
	if r1 == 0 {
		err := syscall.EINVAL
		if e1 != 0 {
			err = error(e1)
		}
		return 0, 0, err
	}
	total := uint64(lpTotalNumberOfBytes)
	free := uint64(lpFreeBytesAvailable)
	return total, free, nil
}
