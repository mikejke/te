package editor

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type Cursor struct {
	X, Y int
}

type Screen struct {
	Rows, Cols int
}

type EditorRow struct {
	String []string
	Render []string
}

type Editor struct {
	Cursor
	Screen
	Row            EditorRow
	Rowoff, Coloff int
	Termios        syscall.Termios
}

func (e *Editor) InitEditor() {
	if err := e.getWindowSize(&e.Screen.Rows, &e.Screen.Cols); err != nil {
		e.Die(err)
	}
}

func (e *Editor) updateRow() {
	var tabs int
	for i := 0; i < len(e.Row.String); i++ {
		if e.Row.String[i] == "\t" {
			tabs++
		}
	}

	e.Row.Render = make([]string, len(e.Row.String)+tabs*(TAB_STOP-1))
	for i, j := 0, 0; i < len(e.Row.String); i++ {
		if e.Row.String[i] == "\t" {
			e.Row.Render[j] = " "
			j++
			for j%TAB_STOP != 0 {
				e.Row.Render[j] = " "
				j++
			}
		} else {
			e.Row.Render[j] = e.Row.String[i]
			j++
		}
	}
}

// OpenFile...
func (e *Editor) OpenFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		e.Row.String = append(e.Row.String, scanner.Text())
	}
	e.updateRow()

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// printWelcomMessage...
func (e *Editor) printWelcomeMessage(buf *string) {
	welcome := fmt.Sprintf("text editor -- version %s", VERSION)

	if len(welcome) > e.Screen.Cols {
		welcome = welcome[:e.Screen.Cols]
	}

	padding := (e.Screen.Cols - len(welcome)) / 2
	if padding > 0 {
		*buf += "~"
		padding--
	}

	for ; padding > 0; padding-- {
		*buf += " "
	}

	*buf += welcome
}

func (e *Editor) Scroll() {
	if e.Cursor.Y < e.Rowoff {
		e.Rowoff = e.Cursor.Y
	}

	if e.Cursor.Y >= e.Rowoff+e.Screen.Rows {
		e.Rowoff = e.Cursor.Y - e.Screen.Rows + 1
	}

	if e.Cursor.X < e.Coloff {
		e.Coloff = e.Cursor.X
	}
	if e.Cursor.X >= e.Coloff+e.Screen.Cols {
		e.Coloff = e.Cursor.X - e.Screen.Cols + 1
	}
}

// DrawRows...
func (e *Editor) DrawRows(buf *string) {
	for y := 0; y < e.Screen.Rows; y++ {
		filerow := y + e.Rowoff
		if filerow >= len(e.Row.String) {
			if len(e.Row.String) == 0 && y == e.Screen.Rows/3 {
				e.printWelcomeMessage(buf)
			} else {
				*buf += "~"
			}
		} else {
			len := len(e.Row.Render[filerow]) - e.Coloff
			if len > e.Screen.Cols {
				*buf += e.Row.Render[filerow][:e.Screen.Cols]
			} else {
				*buf += e.Row.Render[filerow]
			}
		}

		*buf += "\x1b[K"
		if y < e.Screen.Rows-1 {
			*buf += "\r\n"
		}
	}
}

func (e *Editor) RefreshScreen() {
	e.Scroll()

	// hide cursor
	buffer := "\x1b[?25l"
	// move cursor
	buffer += "\x1b[H"

	e.DrawRows(&buffer)

	buffer += fmt.Sprintf("\x1b[%d;%dH", (e.Cursor.Y-e.Rowoff)+1, (e.Cursor.X-e.Coloff)+1)

	// show cursor
	buffer += "\x1b[?25h"

	fmt.Print(buffer)
}

// MoveCursor...
func (e *Editor) MoveCursor(k int) {
	var row *string
	if e.Cursor.Y < len(e.Row.String) {
		row = &e.Row.String[e.Cursor.Y]
	}

	switch k {
	case ARROW_LEFT:
		if e.Cursor.X != 0 {
			e.Cursor.X--
		} else if e.Cursor.Y > 0 {
			e.Cursor.Y--
			e.Cursor.X = len(e.Row.String[e.Cursor.Y])
		}
	case ARROW_RIGHT:
		if row != nil {
			if e.Cursor.X < len(*row) {
				e.Cursor.X++
			} else if e.Cursor.X == len(*row) {
				e.Cursor.Y++
				e.Cursor.X = 0
			}
		}
	case ARROW_UP:
		if e.Cursor.Y != 0 {
			e.Cursor.Y--
		}
	case ARROW_DOWN:
		if e.Cursor.Y < len(e.Row.String) {
			e.Cursor.Y++
		}
	}

	row = nil
	if e.Cursor.Y < len(e.Row.String) {
		row = &e.Row.String[e.Cursor.Y]
	}
	var rowlen int
	if row != nil {
		rowlen = len(*row)
	}

	if e.Cursor.X > rowlen {
		e.Cursor.X = rowlen
	}
}

// HandleKeypress waits for a keypress, and then handles it.
func (e *Editor) HandleKeypress() {
	c := e.ReadKey()

	switch c {
	case CTRL_KEY('q'):
		fmt.Print("\x1b[2J")
		fmt.Print("\x1b[H")
		e.Die(nil)

	case HOME:
		e.Cursor.Y = 0

	case END:
		e.Cursor.X = e.Screen.Cols - 1

	case PAGE_UP, PAGE_DOWN:
		for i := e.Screen.Rows; i > 0; i-- {
			if c == PAGE_UP {
				e.MoveCursor(ARROW_UP)
			} else {
				e.MoveCursor(ARROW_DOWN)
			}
		}

	case ARROW_LEFT, ARROW_RIGHT, ARROW_UP, ARROW_DOWN:
		e.MoveCursor(c)
	}
}

// ReadKey waits for one keypress, and return it.
func (e *Editor) ReadKey() int {

	b := make([]byte, 1)
	if n, err := os.Stdin.Read(b); n != 1 || err != nil {
		return '\x1b'
	}

	c := int(b[0])

	if c == '\x1b' {
		seq := make([]byte, 3)
		if _, err := os.Stdin.Read(seq[:2]); err != nil {
			return '\x1b'
		}

		if seq[0] == '[' {
			if seq[1] >= '0' && seq[1] <= '9' {
				if _, err := os.Stdin.Read(seq[2:]); err != nil {
					return '\x1b'
				}
				if seq[2] == '~' {
					switch seq[1] {
					case '1', '7':
						return HOME
					case '3':
						return DEL
					case '4', '8':
						return END
					case '5':
						return PAGE_UP
					case '6':
						return PAGE_DOWN
					}
				}
			} else {
				switch seq[1] {
				case 'D':
					return ARROW_LEFT
				case 'C':
					return ARROW_RIGHT
				case 'A':
					return ARROW_UP
				case 'B':
					return ARROW_DOWN
				case 'H':
					return HOME
				case 'F':
					return END
				}
			}
		} else if seq[0] == 'O' {
			switch seq[1] {
			case 'H':
				return HOME
			case 'F':
				return END
			}
		}
		return '\x1b'
	}

	return c
}

// getCursorPosition...
func (e *Editor) getCursorPosition(rows, cols *int) error {
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

	e.ReadKey()
	return nil
}

// getWindowSize...
func (e *Editor) getWindowSize(rows, cols *int) error {
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
		return e.getCursorPosition(rows, cols)
	}

	*cols = int(ws.cols)
	*rows = int(ws.rows)

	return nil
}

// DisableRawMode restores the terminal to a previous state.
func (e *Editor) DisableRawMode() error {
	return setTermios(&e.Termios)
}

// EnableRawMode put the terminal connected to the given file descriptor into raw
// mode and returns the previous state of the terminal so that it can be
// restored.
func (e *Editor) EnableRawMode() error {
	err := getTermios(&e.Termios)
	if err != nil {
		return err
	}

	raw := e.Termios

	raw.Iflag &^= (syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK | syscall.ISTRIP | syscall.INLCR | syscall.IGNCR | syscall.ICRNL | syscall.IXON)
	raw.Oflag &^= syscall.OPOST
	raw.Lflag &^= (syscall.ECHO | syscall.ECHONL | syscall.ICANON | syscall.ISIG | syscall.IEXTEN)
	raw.Cflag &^= (syscall.CSIZE | syscall.PARENB)
	raw.Cflag |= syscall.CS8
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0

	err = setTermios(&raw)
	if err != nil {
		return err
	}

	return nil
}

// Die gracefully exits program.
func (e *Editor) Die(err error) {
	fmt.Print("\x1b[2J")
	fmt.Print("\x1b[H")

	if err := e.DisableRawMode(); err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
