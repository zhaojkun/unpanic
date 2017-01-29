package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	line := strings.Join(os.Args[1:], " ")
	for {
		cmd := exec.Command("sh", "-c", line)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println(err)
			continue
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Println(err)
			continue
		}
		go io.Copy(os.Stdout, stdout)
		go io.Copy(os.Stdout, stderr)
		if err := cmd.Start(); err != nil {
			log.Println(err)
			continue
		}
		if err := cmd.Wait(); err != nil {
			log.Println(err)
		}
	}
}
