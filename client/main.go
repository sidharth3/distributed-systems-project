package main

import (
	"bufio"
	"ds-proj/client/commands"
	"ds-proj/client/helpers"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//once inside the shell, can do 'getfile 127.0.0.1:8080 <filename>'
func main() {
	helpers.MockClient()
	os.MkdirAll(helpers.StorageDir(), os.ModePerm)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(("$~ "))
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		err = run(command)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func run(command string) error {
	command = strings.TrimSuffix(command, "\n")
	arrCommand := strings.Fields(command)

	if len(arrCommand) < 3 {
		return errors.New("Missing two args. getfile requires masterIP and file name.")
	}

	switch arrCommand[0] {
	case "exit":
		os.Exit(0)

	case "getfile":
		commands.GetFile(arrCommand[1], arrCommand[2])
		return nil
	case "postfile":
		// you need to manually create a file in your own directory first
		commands.PostFile(arrCommand[1], arrCommand[2])
		return nil
	default:
		fmt.Println("Invalid command.")
		return nil
	}
	cmd := exec.Command(arrCommand[0], arrCommand[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
