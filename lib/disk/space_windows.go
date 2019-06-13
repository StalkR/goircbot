// +build windows

package disk

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Space returns total and free bytes available in a directory, e.g. `C:\`.
// It returns free space available to the user (including quota limitations),
// so it can be lower than the free space of the disk.
func Space(path string) (total, free uint64, err error) {
	kernel32, err := windows.LoadLibrary("Kernel32.dll")
	if err != nil {
		return 0, 0, err
	}
	defer windows.FreeLibrary(kernel32)
	GetDiskFreeSpaceEx, err := windows.GetProcAddress(windows.Handle(kernel32), "GetDiskFreeSpaceExW")
	if err != nil {
		return 0, 0, err
	}
	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)
	r1, _, e1 := syscall.Syscall6(GetDiskFreeSpaceEx, 4,
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)), 0, 0)
	if r1 == 0 {
		err := syscall.EINVAL
		if e1 != 0 {
			err = syscall.Errno(e1)
		}
		return 0, 0, err
	}
	return uint64(lpTotalNumberOfBytes), uint64(lpFreeBytesAvailable), nil
}
