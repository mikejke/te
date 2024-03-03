package editor

import (
	"syscall"
	"unsafe"
)

// getTermios copies the parameters associated with the terminal.
func getTermios(t *syscall.Termios) error {
	_, _, err := syscall.Syscall6(
		syscall.SYS_IOCTL,
		0,
		syscall.TCGETS,
		uintptr(unsafe.Pointer(t)),
		0,
		0,
		0,
	)
	if err != 0 {
		return err
	}

	return nil
}

// setTermios manipulates the termios structure.
func setTermios(t *syscall.Termios) error {
	_, _, err := syscall.Syscall6(
		syscall.SYS_IOCTL,
		0,
		syscall.TCSETS,
		uintptr(unsafe.Pointer(t)),
		0,
		0,
		0,
	)
	if err != 0 {
		return err
	}

	return nil
}
