package execs

import (
	"fmt"
)

func PrintPodman(filename string, prod bool) error {
	comm, _, err := getCommands(filename, prod)
	if err != nil {
		return err
	}

	// Print each command on its own line
	for _, cmd := range *comm {
		fmt.Println("\n" + cmd)
	}

	return nil
}
