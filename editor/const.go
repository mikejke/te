package editor

const (
	VERSION  = "0.0.1"
	TAB_STOP = 8
)

const (
	ARROW_LEFT int = iota + 1000
	ARROW_RIGHT
	ARROW_UP
	ARROW_DOWN
	DEL
	HOME
	END
	PAGE_UP
	PAGE_DOWN
)

// CTRL_KEY macro bitwise-ANDs a character with the value 00011111, in binary.
var CTRL_KEY = func(k int) int {
	return k & 0x1f
}
