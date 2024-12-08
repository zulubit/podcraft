package execs

import (
	"fmt"
	"os"

	"github.com/zulubit/podcraft/pkg/commands"
	"github.com/zulubit/podcraft/pkg/configfile"
	"github.com/zulubit/podcraft/pkg/replaceables"
	"github.com/zulubit/podcraft/pkg/validate"
	"github.com/zulubit/podcraft/pkg/walk"
)

func getCommands(filename string, prod bool) (*[]string, *configfile.Config, error) {
	quatoml, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("Error reading TOML: %v", err)
	}

	config, _, err := configfile.ParseConfigFromTOML(string(quatoml))
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing TOML: %v", err)
	}

	// we should do the replacables magic here, return the entire toml back and run parsefromtoml again
	newConfig, err := replaceables.ReplaceReplaceables(string(quatoml), *config, prod)
	if err != nil {
		return nil, nil, err
	}

	replacedConfig, meta, err := configfile.ParseConfigFromTOML(*newConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing TOML: %v", err)
	}

	err = validate.ValidateNoExtraKeys(meta)
	if err != nil {
		return nil, nil, fmt.Errorf("Extra keys found in TOML: %v", err)
	}

	data, err := walk.WalkQuadlets(replacedConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("Walking failed: %v", err)
	}

	comm, err := commands.ConstructCommands(*data)
	if err != nil {
		return nil, nil, fmt.Errorf("Commands could not be constructed: %v", err)
	}

	return comm, replacedConfig, nil

}
