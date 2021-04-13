package main

import (
	"ds-proj/client/commands"
	"fmt"
	"os"
)

// build client CLI 'go build -o dfs ../client'
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify commands to be performed.")
		fmt.Println("Available dfs commands are: -cat, -get, -put, -rm, -ls, -mkdir")
		os.Exit(0)
	}

	switch os.Args[1] {
	case "exit":
		os.Exit(0)

	case "-cat":
		// example command: ./dfs -cat 127.0.0.1:8080 example.txt
		if len(os.Args) < 4 {
			fmt.Println("Missing two args. cat requires masterIP and file name.")
			os.Exit(0)
		} else {
			commands.GetFile(os.Args[2], os.Args[3])
		}
	case "-get":
		// ./dfs -get 127.0.0.1:8080 /gfspath/to/example.txt localpath/to - DONE
		if len(os.Args) < 5 {
			fmt.Println("Missing three args. get requires masterIP, dfs filepath and local filepath.")
			os.Exit(0)
		} else {
			commands.DownloadFile(os.Args[2], os.Args[3], os.Args[4])
		}
	case "-put":
		// example command: ./dfs -put 127.0.0.1:8080 example.txt
		// ./dfs -put 127.0.0.1:8080 example.txt /gfspath/to - DONE
		if len(os.Args) < 4 {
			fmt.Println("Missing args. put requires masterIP, file name and dfs destination path.")
			os.Exit(0)
		} else {
			// you need to manually create a file in your own directory first
			commands.PostFile(os.Args[2], os.Args[3], os.Args[4])
		}
	case "-rm":
		if len(os.Args) < 4 {
			fmt.Println("Missing two args. rm requires masterIP and file name.")
			os.Exit(0)
		} else {
			commands.DeleteFile(os.Args[2], os.Args[3])
		}
	case "-ls":
		// example command: ./dfs -ls 127.0.0.1:8080 /gfspath/to
		// TODO: ./dfs -ls 127.0.0.1:8080 /gfspath/to (. for root dir)
		if len(os.Args) < 3 {
			fmt.Println("Missing two args. ls requires masterIP and directory.")
			os.Exit(0)
		} else {
			if len(os.Args) == 3 {
				commands.ListDir(os.Args[2], ".")
			} else {
				commands.ListDir(os.Args[2], os.Args[3])
			}
		}
	case "-mkdir":
		// TODO: ./dfs -mkdir 127.0.0.1:8080 /gfspath/to
		if len(os.Args) < 3 {
			fmt.Println("Missing two args. ls requires masterIP and directory.")
			os.Exit(0)
		} else {
			fmt.Println("Have not implemented.")
			os.Exit(0)
		}
	default:
		fmt.Println("Invalid command.")
	}
}
