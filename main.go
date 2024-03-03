package main

import (
	"os"

	"github.com/mikejke/go-te/editor"
)

var e = editor.Editor{}

func main() {
	err := e.EnableRawMode()
	if err != nil {
		e.Die(err)
	}
	e.InitEditor()

	if len(os.Args) >= 2 {
		e.OpenFile(os.Args[1])
	}

	for {
		e.RefreshScreen()
		e.HandleKeypress()
	}
}
