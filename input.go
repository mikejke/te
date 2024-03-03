package main

import "fmt"

// editorMoveCursor...
func editorMoveCursor(k rune) {
	switch k {
	case 'w':
		Config.cy--
	case 'a':
		Config.cx--
	case 's':
		Config.cy++
	case 'd':
		Config.cx++
	}
}

// editorProcessKeypress waits for a keypress, and then handles it.
func editorProcessKeypress() {
	c, err := editorReadKey()
	if err != nil {
		die(1)
	}

	switch c {
	case CTRL_KEY('q'):
		fmt.Print("\x1b[2J")
		fmt.Print("\x1b[H")
		die(0)
	case 'w', 'a', 's', 'd':
		editorMoveCursor(c)
	}
}
