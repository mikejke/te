package main

func main() {
	err := enableRawMode(0)
	if err != nil {
		die(err)
	}
	NewConfig()

	for {
		editorRefreshScreen()
		editorProcessKeypress()
	}
}
