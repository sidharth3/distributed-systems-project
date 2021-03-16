package main

import (
	"ds-proj/client/commands"
)

func main() {
	commands.GetFileMaster("127.0.0.1:8080", "test_file.txt")
}
