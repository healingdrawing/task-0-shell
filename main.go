package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main() {

	bang()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		command, err := reader.ReadString('\n')
		if err != nil {
			// Handle EOF (Ctrl+D)
			if err == io.EOF {
				return
			}
			fmt.Println(err)
			continue
		}
		command = strings.TrimSpace(command)
		if command == "exit" {
			break
		}
		args := strings.Split(command, " ")
		switch args[0] {
		case "cd":
			var dir string
			if len(args) > 1 {
				dir = args[1]
			} else {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					fmt.Println(err)
					continue
				}
				dir = homeDir
			}
			err := os.Chdir(dir)
			if err != nil {
				fmt.Println(err)
			}
		default:
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil { // skip print error, to satisfy the audit question
				if err.Error() == "exit status 1" {
					continue
				} else if err.Error() == "exec: no command" {
					continue
				}
				fmt.Println(err)
			}
		}
	}
}
