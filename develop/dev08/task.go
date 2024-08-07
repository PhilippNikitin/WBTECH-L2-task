package main

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

*/

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// cd command implementation
func cd(args []string) {
	if len(args) < 1 {
		fmt.Println("cd: missing argument")
		return
	}
	err := os.Chdir(args[0])
	if err != nil {
		fmt.Printf("cd: %s\n", err)
	}
}

// pwd command implementation
func pwd() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("pwd: %s\n", err)
	} else {
		fmt.Println(dir)
	}
}

// echo command implementation
func echo(args []string) {
	fmt.Println(strings.Join(args, " "))
}

// kill command implementation
func kill(args []string) {
	if len(args) < 1 {
		fmt.Println("kill: missing argument")
		return
	}
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Printf("kill: %s\n", err)
		return
	}
	err = syscall.Kill(pid, syscall.SIGKILL)
	if err != nil {
		fmt.Printf("kill: %s\n", err)
	}
}

// ps command implementation
func ps() {
	cmd := exec.Command("ps", "-e", "-o", "pid,cmd") // Use ps options for compatibility
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("ps: %s\n", err)
	}
}

func executePipeline(commands [][]string) {
	if len(commands) == 0 {
		return
	}

	var prevOutput bytes.Buffer

	for i, cmdArgs := range commands {
		if len(cmdArgs) == 0 {
			continue
		}

		// Create the command
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

		// Set the stdin for commands other than the first
		if i > 0 {
			cmd.Stdin = &prevOutput
		}

		// Set stdout to a new buffer for intermediate results
		var out bytes.Buffer
		cmd.Stdout = &out

		// Execute the command
		err := cmd.Run()
		if err != nil {
			fmt.Printf("command error: %s\n", err)
			return
		}

		// Reset the buffer for the next command
		prevOutput = out
	}

	// Print final output to stdout
	fmt.Print(prevOutput.String())
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Split the line into commands connected by pipes
		commands := strings.Split(line, "|")
		var parsedCommands [][]string
		for _, cmd := range commands {
			parsedCommands = append(parsedCommands, strings.Fields(cmd))
		}

		if len(parsedCommands) == 1 {
			args := parsedCommands[0]
			switch args[0] {
			case "cd":
				cd(args[1:])
			case "pwd":
				pwd()
			case "echo":
				echo(args[1:])
			case "kill":
				kill(args[1:])
			case "ps":
				ps()
			default:
				cmd := exec.Command(args[0], args[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					fmt.Printf("%s: %s\n", args[0], err)
				}
			}
		} else {
			executePipeline(parsedCommands)
		}
	}
}
