package execs

import (
	"fmt"
	"os"
	"os/exec"
)

func CreatePodman(filename string, prod bool) (*string, error) {

	comm, config, err := getCommands(filename, prod)
	if err != nil {
		return nil, err
	}
	err = tryRunCommands(comm, config.Pod.Name)
	if err != nil {
		return nil, err
	}

	return &config.Pod.Name, nil
}

func tryRunCommands(commandSlice *[]string, podName string) error {

	for _, c := range *commandSlice {
		fmt.Printf("\nRunning: %s\n", c)

		cmd := exec.Command("bash", "-c", c)
		// Set the command's stdout and stderr to the user's terminal
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Println("Fail! deleting the created pod")
			exec.Command("bash", "-c", "podman pod rm -f "+podName).Run()
			return err
		}
		fmt.Println("Success!")
	}

	return nil
}
