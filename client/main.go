package main

import (
	"ds-proj/client/commands"
	"fmt"
	"os"
)

// build client CLI 'go build -o dfs ./client'
// once inside the shell, can do 'getfile 127.0.0.1:8080 <filename>'
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify operation to be performed. \n Available operations are: getfile, postfile, deletefile")
		os.Exit(0)
	}

	switch os.Args[1] {
	case "exit":
		os.Exit(0)

	case "getfile":
		if len(os.Args) < 4 {
			fmt.Println("Missing two args. getfile requires masterIP and file name.")
			os.Exit(0)
		} else {
			commands.GetFile(os.Args[2], os.Args[3])
		}
	case "postfile":
		if len(os.Args) < 4 {
			fmt.Println("Missing two args. postfile requires masterIP and file name.")
			os.Exit(0)
		} else {
			// you need to manually create a file in your own directory first
			commands.PostFile(os.Args[2], os.Args[3])
		}
	case "deletefile":
		if len(os.Args) < 4 {
			fmt.Println("Missing two args. deletefile requires masterIP and file name.")
			os.Exit(0)
		} else {
			commands.DeleteFile(os.Args[2], os.Args[3])
		}
	case "ls":
		if len(os.Args) < 3 {
			fmt.Println("Missing two args. ls requires masterIP and directory.")
			os.Exit(0)
		} else {
			if len(os.Args) == 3 {
				commands.ListDir(os.Args[2], "")
			} else {
				commands.ListDir(os.Args[2], os.Args[3])
			}
		}
	default:
		fmt.Println("Invalid command.")
	}
}
