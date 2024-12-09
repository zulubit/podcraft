package execs

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/zulubit/podcraft/pkg/color"
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
		fmt.Printf("\n"+color.ColorYellow+"Running:"+color.ColorReset+" %s\n", c)

		cmd := exec.Command("bash", "-c", c)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Println(color.ColorRed + "Fail! deleting the created pod" + color.ColorReset)
			exec.Command("bash", "-c", "podman pod rm -f "+podName).Run()
			return err
		}
		fmt.Println(color.ColorGreen + "Success!" + color.ColorReset)
	}

	return nil
}
