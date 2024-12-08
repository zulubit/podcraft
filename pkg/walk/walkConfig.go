package walk

import (
	"errors"
	"fmt"

	"github.com/zulubit/podcraft/pkg/configfile"
	"github.com/zulubit/podcraft/pkg/readini"
	"github.com/zulubit/podcraft/pkg/validate"
)

type Actionables []Actionable

type Actionable struct {
	Type     string
	MainName string
	Options  map[string][]string
}

func WalkQuadlets(config *configfile.Config) (*Actionables, error) {
	acc := Actionables{}

	// Parse pod options
	podOptions, err := readini.ReadDataFromIni(config.Pod.Quadlet)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pod quadlet: %w", err)
	}

	podData, ok := podOptions["Pod"]
	if !ok {
		return nil, fmt.Errorf("missing [Pod] section in pod quadlet")
	}

	podActionable := Actionable{
		Type:     "Pod",
		MainName: config.Pod.Name,
		Options:  podData,
	}
	acc = append(acc, podActionable)

	// Parse quadlets
	for _, d := range config.Quadlets {
		// Validate unit type
		if err := validate.ValidateUnitType(d.Type); err != nil {
			return nil, fmt.Errorf("invalid unit type for quadlet %s: %w", d.Name, err)
		}

		quadOptions, err := readini.ReadDataFromIni(d.Quadlet)
		if err != nil {
			return nil, fmt.Errorf("failed to parse quadlet %s: %w", d.Name, err)
		}

		quadData, ok := quadOptions[d.Type]
		if !ok {
			return nil, fmt.Errorf("missing [%s] section in quadlet %s", d.Type, d.Name)
		}

		// Validate container pod association
		if d.Type == "Container" {
			podName := ""
			if podList, exists := quadData["Pod"]; exists && len(podList) > 0 {
				podName = podList[0] // Assume the first Pod value for validation
			}

			containerName := ""
			if containerList, exists := quadData["ContainerName"]; exists && len(containerList) > 0 {
				containerName = containerList[0] // Assume the first Container value for validation
			}

			err := validate.ValidateContainerPod(podName, config.Pod.Name, d.Name)
			if err != nil {
				return nil, err
			}

			err = validate.ValidateContainerName(containerName, d.Name)
			if err != nil {
				return nil, err
			}

			if Image, exists := quadData["Image"]; !exists || len(Image) == 0 {
				return nil, errors.New("[Container] units must have an 'Image' defined")
			}
		}

		quadActionable := Actionable{
			Type:     d.Type,
			MainName: d.Name,
			Options:  quadData,
		}

		acc = append(acc, quadActionable)
	}

	return &acc, nil
}
