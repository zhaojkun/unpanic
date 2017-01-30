package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

var currentCmd *exec.Cmd

func main() {
	initSignal()
	var err error
	line := extractParam()
	for {
		currentCmd, err = runCommand(line)
		if err != nil {
			log.Println(err)
			continue
		}
		if err := wait(currentCmd); err != nil {
			log.Println(err)
		}
	}
}

func initSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			if currentCmd != nil {
				log.Println("Got signal:", sig)
				currentCmd.Process.Signal(sig)
			}
		}
	}()
}

func extractParam() string {
	return strings.Join(os.Args[1:], " ")
}

func runCommand(line string) (*exec.Cmd, error) {
	cmd := exec.Command("sh", "-c", line)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stdout, stderr)
	if err := cmd.Start(); err != nil {
		return cmd, err
	}
	return cmd, nil
}

func wait(cmd *exec.Cmd) error {
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
