package windows

import (
	"syscall"
	"unsafe"
	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
	errERROR_EINVAL     error = syscall.EINVAL
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return errERROR_EINVAL
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	return e
}

var (
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")
	procReadConsoleW = modkernel32.NewProc("ReadConsoleW")
	procReadConsoleInputW = modkernel32.NewProc("ReadConsoleInputW")
)

// ReadConsole reads characters from the console and writes the Unicode keycode(s) into buf
// buf should point to the element in an array where writing should begin
// toread specifies how many characters should be read
// read should point to a uint32 where the number of characters actually read should be stored
// inputControl should be set to NULL, as it currently has no effect
// See: https://docs.microsoft.com/en-us/windows/console/readconsole
func ReadConsole(console windows.Handle, buf *uint16, toread uint32, read *uint32, inputControl *byte) (err error) {
	r1, _, e1 := syscall.Syscall6(procReadConsoleW.Addr(), 5, uintptr(console), uintptr(unsafe.Pointer(buf)), uintptr(toread), uintptr(unsafe.Pointer(read)), uintptr(unsafe.Pointer(inputControl)), 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

// ReadConsoleInput reads keypresses from console.
// The data is read into the array starting at buf.
// toread is the amount of keypresses that should be read.
// The actual amount of keypresses read is stored in read.
// The difference between ReadConsole and ReadConsoleInput is that 
//  - ReadConsole only reads character insertion (reads any character key pressed)
//  - ReadConsoleInput reads any key (both key press and key release) as well as mouse, focus and window size change events
// See: https://docs.microsoft.com/en-us/windows/console/readconsoleinput
func ReadConsoleInput(console windows.Handle, rec *InputRecord, toread uint32, read *uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procReadConsoleInputW.Addr(), 4,
		uintptr(console), uintptr(unsafe.Pointer(rec)), uintptr(toread),
		uintptr(unsafe.Pointer(read)), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

