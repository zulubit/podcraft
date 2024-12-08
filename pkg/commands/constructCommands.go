package commands

import (
	"github.com/zulubit/podcraft/pkg/commandmap"
	"github.com/zulubit/podcraft/pkg/walk"
)

func ConstructCommands(actionables walk.Actionables) (*[]string, error) {
	var pod string
	var networks []string
	var volumes []string
	var builds []string
	var images []string
	var containers []string

	var commands []string

	for _, a := range actionables {
		switch a.Type {
		case "Pod":
			c, err := commandmap.GeneratePodmanPodCommand(a.MainName, a.Options)
			if err != nil {
				return nil, err
			}
			pod = c
		case "Network":
			c, err := commandmap.GeneratePodmanNetworkCommand(a.MainName, a.Options)
			if err != nil {
				return nil, err
			}
			networks = append(networks, c)
		case "Volume":
			c, err := commandmap.GeneratePodmanVolumeCommand(a.MainName, a.Options)
			if err != nil {
				return nil, err
			}
			volumes = append(volumes, c)
		case "Build":
			c, err := commandmap.GeneratePodmanBuildCommand(a.Options)
			if err != nil {
				return nil, err
			}
			builds = append(builds, c)
		case "Image":
			c, err := commandmap.GeneratePodmanImageCommand(a.Options)
			if err != nil {
				return nil, err
			}
			images = append(images, c)
		case "Container":
			c, err := commandmap.GeneratePodmanContainerCommand(a.Options)
			if err != nil {
				return nil, err
			}
			containers = append(containers, c)
		}
	}

	commands = append(commands, networks...)
	commands = append(commands, volumes...)
	commands = append(commands, builds...)
	commands = append(commands, images...)
	commands = append(commands, pod)
	commands = append(commands, containers...)

	return &commands, nil
}
