package execs

import (
	"os"
	"os/exec"
)

func TryStartPod(podName string) error {

	cmd := exec.Command("bash", "-c", "podman pod start "+podName)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
