package commandmap

import (
	"fmt"
	"strings"
)

type PodField struct {
	Type        string // "default" or "special"
	Flag        string // Podman pod flag or empty for special handling
	Description string // Description of the field
}

var PodToPodman = map[string]PodField{
	// Default mappings
	"AddHost":              {"default", "--add-host %s", ""},
	"ContainersConfModule": {"default", "--module %s", ""},
	"DNS":                  {"default", "--dns %s", ""},
	"DNSOption":            {"default", "--dns-option %s", ""},
	"DNSSearch":            {"default", "--dns-search %s", ""},
	"GIDMap":               {"default", "--gidmap %s", ""},
	"GlobalArgs":           {"default", "--log-level %s", ""},
	"IP":                   {"default", "--ip %s", ""},
	"IP6":                  {"default", "--ip6 %s", ""},
	"Network":              {"default", "--network %s", ""},
	"NetworkAlias":         {"default", "--network-alias %s", ""},
	"PublishPort":          {"default", "--publish %s", ""},
	"SubGIDMap":            {"default", "--subgidname %s", ""},
	"SubUIDMap":            {"default", "--subuidname %s", ""},
	"UIDMap":               {"default", "--uidmap %s", ""},
	"UserNS":               {"default", "--userns %s", ""},
	"Volume":               {"default", "--volume %s", ""},

	// Special mappings requiring custom handling
	"PodName":     {"special", "--name %s", "Overrides the default pod name."},
	"PodmanArgs":  {"special", "%s", "Additional Podman arguments provided."},
	"ServiceName": {"special", "", "Name the systemd unit."}, // Systemd-specific handling
}

// formatPodFlag formats a single pod-related flag with its corresponding value.
func formatPodFlag(key, value string) (string, error) {
	field, ok := PodToPodman[key]
	if !ok {
		return "", fmt.Errorf("unknown key: %s", key)
	}

	if field.Type == "special" {
		// Special fields are handled separately
		return "", nil
	}

	if value == "" && !strings.Contains(field.Flag, "%s") {
		return field.Flag, nil // Boolean flags like --internal
	}

	return fmt.Sprintf(field.Flag, value), nil
}

// GeneratePodmanPodCommand generates a `podman pod create` command from the provided options.
func GeneratePodmanPodCommand(name string, options map[string][]string) (string, error) {
	parts := []string{}

	for key, values := range options {
		for _, value := range values {
			part, err := formatPodFlag(key, value)
			if err != nil {
				return "", fmt.Errorf("error formatting pod flag '%s': %w", key, err)
			}
			parts = append(parts, part)
		}
	}

	return fmt.Sprintf(
		"podman pod exists %s && podman pod rm -f %s; podman pod create %s %s",
		name,
		name,
		strings.Join(parts, " "),
		name,
	), nil
}
