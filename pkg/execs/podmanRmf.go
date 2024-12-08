package execs

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/zulubit/podcraft/pkg/color"
	"github.com/zulubit/podcraft/pkg/configfile"
	"github.com/zulubit/podcraft/pkg/replaceables"
)

func PodmanRmf(filename string, prod bool) error {
	quatoml, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Error reading TOML: %v", err)
	}

	config, _, err := configfile.ParseConfigFromTOML(string(quatoml))
	if err != nil {
		return fmt.Errorf("Error parsing TOML: %v", err)
	}

	// we should do the replacables magic here, return the entire toml back and run parsefromtoml again
	newConfig, err := replaceables.ReplaceReplaceables(string(quatoml), *config, prod)
	if err != nil {
		return err
	}

	replacedConfig, _, err := configfile.ParseConfigFromTOML(*newConfig)
	if err != nil {
		return fmt.Errorf("Error parsing TOML: %v", err)
	}

	cs := "podman pod rm -f " + replacedConfig.Pod.Name + " && podman network prune -f"

	cmd := exec.Command("bash", "-c", cs)

	fmt.Printf(color.ColorYellow+"Runnign:"+color.ColorReset+" %s\n\n", cs)

	// Set the command's stdout and stderr to the user's terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println("\nDeleted resources above.")

	return nil
}
