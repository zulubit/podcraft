package execs

import (
	"fmt"
)

func PrintPodman(filename string, prod bool) error {
	comm, _, err := getCommands(filename, prod)
	if err != nil {
		return err
	}

	for _, cmd := range *comm {
		fmt.Println("\n" + cmd)
	}

	return nil
}
