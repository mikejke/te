package main

import "syscall"

type EditorConfig struct {
	cx, cy      int
	srows       int
	scols       int
	origTermios syscall.Termios
}

var Config EditorConfig

func NewConfig() {
	if err := getWindowSize(&Config.srows, &Config.scols); err != nil {
		die(err)
	}
}
