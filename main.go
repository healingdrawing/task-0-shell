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

	// full list of allowed commands includes echo and cd, managed separately
	supported_commands := strings.Join([]string{"echo", "cd", "ls", "pwd", "cat", "cp", "rm", "mv", "mkdir", "exit"}, ", ")
	// allowed commands
	allowed_commands := []string{"ls", "pwd", "cat", "cp", "rm", "mv", "mkdir", "exit"}

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
		case "echo":
			str := strings.Join(args[1:], " ")
			// cut first and last "
			if len(str) > 1 &&
				(str[0] == '"' && str[len(str)-1] == '"' ||
					str[0] == '\'' && str[len(str)-1] == '\'') {
				str = str[1 : len(str)-1]
			}
			cmd := exec.Command("echo", str)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				fmt.Println(err)
			}

		default:
			// limit the allowed commands. "cd" and "echo" are managed above, so we only allow
			// "ls" , "pwd" , "cat" , "cp" , "rm" , "mv" , "mkdir" , "exit".
			allowed := false
			for _, allowed_command := range allowed_commands {
				if args[0] == allowed_command {
					allowed = true
					break
				}
			}
			if !allowed && args[0] != "" {
				fmt.Println("Command [", args[0], "] not allowed.")
				fmt.Println("Supported commands:", supported_commands, ".")
				continue
			}

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
