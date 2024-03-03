package main

import "fmt"

// editorMoveCursor...
func editorMoveCursor(k int) {
	switch k {
	case ARROW_LEFT:
		Config.cx = Clamp(Config.cx-1, 0, Config.scols-1)
	case ARROW_RIGHT:
		Config.cx = Clamp(Config.cx+1, 0, Config.scols-1)
	case ARROW_UP:
		Config.cy = Clamp(Config.cy-1, 0, Config.srows-1)
	case ARROW_DOWN:
		Config.cy = Clamp(Config.cy+1, 0, Config.srows-1)
	}
}

func Clamp(f, low, high int) int {
	if f < low {
		return low
	}
	if f > high {
		return high
	}
	return f
}

// editorProcessKeypress waits for a keypress, and then handles it.
func editorProcessKeypress() {
	c := editorReadKey()

	switch c {
	case CTRL_KEY('q'):
		fmt.Print("\x1b[2J")
		fmt.Print("\x1b[H")
		die(nil)

	case HOME:
		Config.cx = 0

	case END:
		Config.cx = Config.scols - 1

	case PAGE_UP, PAGE_DOWN:
		for i := Config.srows; i > 0; i-- {
			if c == PAGE_UP {
				editorMoveCursor(ARROW_UP)
			} else {
				editorMoveCursor(ARROW_DOWN)
			}
		}

	case ARROW_LEFT, ARROW_RIGHT, ARROW_UP, ARROW_DOWN:
		editorMoveCursor(c)
	}
}
