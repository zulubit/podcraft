package validate

import (
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/zulubit/podcraft/pkg/configfile"
)

// validateUnits checks if the required fields are present in the parsed Units.
func ValidateUnits(units *configfile.Config) error {
	if units.Pod.Name == "" {
		return errors.New("pod name is missing or invalid")
	}

	if len(units.Quadlets) == 0 {
		return errors.New("no quadlets defined")
	}

	return nil
}

// ValidateNoExtraKeys checks if there are extra keys in the TOML metadata
func ValidateNoExtraKeys(metaData *toml.MetaData) error {
	if metaData == nil {
		return fmt.Errorf("metadata is nil")
	}

	if undecoded := metaData.Undecoded(); len(undecoded) > 0 {
		return fmt.Errorf("extra keys found in TOML: %v", undecoded)
	}

	return nil
}

// validate missing keys
func ValidateMissingKeys(config *configfile.Config) error {
	if config.Pod.Name == "" || config.Pod.Quadlet == "" {
		return fmt.Errorf("[main_pod] configuration is incomplete. Both name and quadlet must be set")
	}

	for _, q := range config.Quadlets {
		if q.Name == "" || q.Quadlet == "" || q.Type == "" {
			return fmt.Errorf("[[quadlet]] configuration is incomplete. Name, quadlet and type must be set")
		}
	}

	return nil
}

// validate unit type
func ValidateUnitType(unitType string) error {
	if unitType == "Pod" || unitType == "pod" {
		return errors.New("Pod definitions are not allowed in the quadlets section, please use the pod section only.")
	}
	if unitType != "Network" && unitType != "Container" && unitType != "Volume" && unitType != "Image" && unitType != "Build" && unitType != "Kube" {
		return fmt.Errorf("invalid quadlets type present in the [[quadlets]] section: %s \nAllowed values: Network, Container, Volume, Image, Build\nThis option is case sensitive!", unitType)
	}
	return nil
}

// validates containers use the right pod
func ValidateContainerPod(containerPod string, podName string, containerName string) error {
	if containerPod != podName {
		return fmt.Errorf("pod name mismatch in container '%s' and pod '%s'. Every container must be part of the top level pod.\n", containerName, podName)
	}
	return nil
}

// validate containers have names matching quadlet names
func ValidateContainerName(containerName string, quadletName string) error {
	if containerName != quadletName {
		return fmt.Errorf("container name '%s' does not match quadlet name '%s'.\nConainerName might also be missing.\nReplace/add your container name with 'ConainerName=%s'\n", containerName, quadletName, quadletName)
	}
	return nil
}
