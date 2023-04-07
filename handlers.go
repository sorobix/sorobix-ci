package main

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"os/exec"
	"time"
)

func deployer(repo Repository) (string, error) {
	id := uuid.New()
	return id.String(), executor(id.String(), repo)
}

func executor(uuid string, repo Repository) error {
	fmt.Println("executor initiating")
	cmd := exec.Command(SCRIPT)
	stderr, err := cmd.StderrPipe()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	results, _ := io.ReadAll(stdout)
	errors, _ := io.ReadAll(stderr)
	d := Deployment{
		ID:     uuid,
		StdOut: results,
		StdErr: errors,
		Time:   time.Now().String(),
	}
	fmt.Println(results)
	err = repo.InsertDeployment(d)
	if err != nil {
		return err
	}
	return nil
}
