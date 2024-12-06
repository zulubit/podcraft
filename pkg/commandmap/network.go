package commandmap

import (
	"fmt"
	"strings"
)

type NetworkField struct {
	Type        string // "default" or "special"
	Flag        string // Podman network flag or empty for special handling
	Description string // Description of the field
}

var NetworkToPodman = map[string]NetworkField{
	// Default mappings
	"ContainersConfModule": {"default", "--module %s", ""},
	"DisableDNS":           {"default", "--disable-dns", ""},
	"DNS":                  {"default", "--dns %s", ""},
	"Driver":               {"default", "--driver %s", ""},
	"Gateway":              {"default", "--gateway %s", ""},
	"GlobalArgs":           {"default", "--log-level %s", ""},
	"Internal":             {"default", "--internal", ""},
	"IPAMDriver":           {"default", "--ipam-driver %s", ""},
	"IPRange":              {"default", "--ip-range %s", ""},
	"IPv6":                 {"default", "--ipv6", ""},
	"Label":                {"default", "--label %s", ""},
	"Options":              {"default", "--opt %s", ""},
	"PodmanArgs":           {"default", "%s", "Additional Podman arguments provided."},
	"Subnet":               {"default", "--subnet %s", ""},

	// Special mappings requiring custom handling
	"NetworkName": {"special", "", "Specifies the name of the network."}, // Required positional argument
}

// formatNetworkFlag formats a single network-related flag with its corresponding value.
func formatNetworkFlag(key, value string) (string, error) {
	field, ok := NetworkToPodman[key]
	if !ok {
		return "", fmt.Errorf("unknown key: %s", key)
	}

	if field.Type == "special" {
		// Special fields like NetworkName are handled separately
		return "", nil
	}

	if value == "" && !strings.Contains(field.Flag, "%s") {
		return field.Flag, nil // Boolean flags like --internal
	}

	return fmt.Sprintf(field.Flag, value), nil
}

// GeneratePodmanNetworkCommand generates a `podman network create` command from the provided options.
func GeneratePodmanNetworkCommand(name string, options map[string][]string) (string, error) {
	parts := []string{}

	for key, values := range options {
		for _, value := range values {
			part, err := formatNetworkFlag(key, value)
			if err != nil {
				return "", fmt.Errorf("error formatting network flag '%s': %w", key, err)
			}
			parts = append(parts, part)
		}
	}

	return fmt.Sprintf(
		"podman network exists %s || podman network create %s %s",
		name,
		strings.Join(parts, " "),
		name,
	), nil
}
