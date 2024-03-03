package main

func main() {
	err := enableRawMode(0)
	if err != nil {
		die(1)
	}
	NewConfig()

	for {
		editorRefreshScreen()
		editorProcessKeypress()
	}
}
