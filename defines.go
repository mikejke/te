package main

// CTRL_KEY macro bitwise-ANDs a character with the value 00011111, in binary.
func CTRL_KEY(k rune) rune {
	return k & 0x1f
}
