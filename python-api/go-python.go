package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

func parseArgs(args []string) (string, string) {
	// check if there is at least 2 arguments
	if len(args) < 2 {
		fmt.Println("Please provide an engine and a process")
		return "", ""
	}
	engine := args[0]
	process := args[1]
	// check if process end with .py
	if process[len(process)-3:] != ".py" {
		fmt.Println("Please provide a python file")
		return "", ""
	}
	return engine, process

}

func main() {

	// parse args
	// engine, process
	engine, process := parseArgs(os.Args[1:])
	if engine == "" || process == "" {
		return
	}

	// init cmd
	cmd := exec.Command(engine, process)

	// wait during 5 seconds
	time.Sleep(2 * time.Second)
	// Use a pipe to pass data to the script's stdin
	stdin, _ := cmd.StdinPipe() // StdinPipe returns a pipe that will be connected to the command's
	// standard input when the command starts.
	go func() {
		defer stdin.Close()
		fmt.Fprint(stdin, "Hello Python")
	}()

	// Use a pipe to receive data from the script's stdout
	stdout, _ := cmd.StdoutPipe()
	go func() {
		bytes, _ := ioutil.ReadAll(stdout)
		if len(bytes) > 0 {
			fmt.Printf("Python script output: %v", string(bytes))
		}
	}()

	// Run the script
	cmd.Run()
}
