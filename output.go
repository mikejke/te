package main

import (
	"fmt"
)

func editorDrawRows(buf *string) {
	for y := 0; y < Config.srows; y++ {
		welcome := fmt.Sprintf("text editor -- version %s", VERSION)
		if y == Config.srows/3 {
			if len(welcome) > Config.scols {
				welcome = welcome[:Config.scols]
			}

			padding := (Config.scols - len(welcome)) / 2
			if padding > 0 {
				*buf += "~"
				padding--
			}

			for ; padding > 0; padding-- {
				*buf += " "
			}
			*buf += welcome

		} else {
			*buf += "~"
		}

		*buf += "\x1b[K"
		if y < Config.srows-1 {
			*buf += "\r\n"
		}
	}
}

func editorRefreshScreen() {
	// hide cursor
	buffer := "\x1b[?25l"
	// move cursor
	buffer += "\x1b[H"

	editorDrawRows(&buffer)

	buffer += fmt.Sprintf("\x1b[%d;%dH", Config.cy+1, Config.cx+1)

	// show cursor
	buffer += "\x1b[?25h"

	fmt.Print(buffer)
}
