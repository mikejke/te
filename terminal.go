package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// die gracefully exits program.
func die(code int) {
	fmt.Print("\x1b[2J")
	fmt.Print("\x1b[H")

	err := disableRawMode()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(code)
}

// getTermios copies the parameters associated with the terminal.
func getTermios(fd uintptr, t *syscall.Termios) error {
	_, _, err := syscall.Syscall6(
		syscall.SYS_IOCTL,
		fd,
		syscall.TCGETS,
		uintptr(unsafe.Pointer(t)),
		0, 0, 0)
	if err != 0 {
		return err
	}

	return nil
}

// setTermios manipulates the termios structure.
func setTermios(fd uintptr, term *syscall.Termios) error {
	_, _, err := syscall.Syscall6(
		syscall.SYS_IOCTL,
		fd,
		syscall.TCSETS,
		uintptr(unsafe.Pointer(term)),
		0, 0, 0)
	if err != 0 {
		return err
	}

	return nil
}

// disableRawMode restores the terminal to a previous state.
func disableRawMode() error {
	return setTermios(0, &Config.origTermios)
}

// enableRawMode put the terminal connected to the given file descriptor into raw
// mode and returns the previous state of the terminal so that it can be
// restored.
func enableRawMode(fd uintptr) error {
	err := getTermios(fd, &Config.origTermios)
	if err != nil {
		return err
	}

	raw := Config.origTermios

	raw.Iflag &^= (syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON)
	raw.Oflag &^= syscall.OPOST
	raw.Lflag &^= (syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN)
	raw.Cflag &^= (syscall.CSIZE | syscall.PARENB)
	raw.Cflag |= syscall.CS8
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0

	err = setTermios(fd, &raw)
	if err != nil {
		return err
	}

	return nil
}

// editorReadKey waits for one keypress, and return it.
func editorReadKey() (rune, error) {
	b := make([]byte, 1)
	_, err := syscall.Read(0, b)
	if err != nil {
		return 0, err
	}

	return rune(b[0]), nil
}

// getCursorPosition...
func getCursorPosition(rows, cols *int) error {
	buf := make([]byte, 32)
	var i uint = 0

	fmt.Print("\x1b[6n\r\n")

	c := make([]byte, 1)
	for i < uint(len(buf)-1) {
		_, err := os.Stdin.Read(c)
		if err == nil {
			break
		}

		if buf[i] == 'R' {
			break
		}

		i++
	}

	buf[i] = 0

	if buf[0] != '\x1b' || buf[1] != '[' {
		return fmt.Errorf("getCursorPosition")
	}

	if _, err := fmt.Sscanf(fmt.Sprintf("%v", &buf[2]), "%d;%d", rows, cols); err != nil {
		return err
	}

	editorReadKey()
	return nil
}

// getWindowSize...
func getWindowSize(rows, cols *int) error {
	var ws struct {
		rows    uint16
		cols    uint16
		xpixels uint16
		ypixels uint16
	}
	_, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)),
	)
	if err != 0 || ws.cols == 0 {
		fmt.Print("\x1b[999C\x1b[999B")
		return getCursorPosition(rows, cols)
	}

	*cols = int(ws.cols)
	*rows = int(ws.rows)

	return nil
}
