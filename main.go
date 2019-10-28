package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {

	if runtime.GOOS == "windows" {
		fmt.Println("can't run on this environment")
	} else {
		path := os.Getenv("PATH")
		fmt.Println(path)
		hostname, err := os.Hostname()
		if err != nil {
			os.Exit(0)
		}
		fmt.Println(hostname)
		cmd := exec.Command("ls", "-l")
		output, stderr := cmd.Output()
		if stderr != nil {
			fmt.Println(stderr)
		}
		fmt.Println(string(output))

	}

}
