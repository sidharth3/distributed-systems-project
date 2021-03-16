package main

import (
	"ds-proj/client/commands"
)

func main() {
	commands.Getfile_master("127.0.0.1:8080", "test_file.txt")
}
