package execs

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/zulubit/podcraft/pkg/configfile"
	"github.com/zulubit/podcraft/pkg/replaceables"
)

func PutQuadlets(filename string, prod bool, location string) error {
	if location == "" {
		location = "./podcraft"
	}

	quatoml, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Error reading TOML: %v", err)
	}

	config, _, err := configfile.ParseConfigFromTOML(string(quatoml))
	if err != nil {
		return fmt.Errorf("Error parsing TOML: %v", err)
	}

	newConfig, err := replaceables.ReplaceReplaceables(string(quatoml), *config, prod)
	if err != nil {
		return err
	}

	replacedConfig, _, err := configfile.ParseConfigFromTOML(*newConfig)
	if err != nil {
		return fmt.Errorf("Error parsing TOML: %v", err)
	}

	err = os.MkdirAll(location, 0755)
	if err != nil {
		return err
	}

	podFilename := path.Join(location, replacedConfig.Pod.Name+".pod")
	fmt.Printf("Writing pod to %s", podFilename)
	err = os.WriteFile(podFilename, []byte(replacedConfig.Pod.Quadlet), 0644)
	if err != nil {
		return err
	}

	for _, q := range replacedConfig.Quadlets {
		qp := path.Join(location, q.Name+"."+strings.ToLower(q.Type))
		fmt.Printf("\nWriting a quadlet to %s", qp)
		err = os.WriteFile(qp, []byte(q.Quadlet), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
