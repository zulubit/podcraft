package commandmap

import (
	"fmt"
	"strings"
)

type BuildField struct {
	Type        string // "default" or "special"
	Flag        string // Podman build flag or empty for special handling
	Description string // Description of the field
}

var BuildToPodman = map[string]BuildField{
	"Annotation":           {"default", "--annotation %s", ""},
	"Arch":                 {"default", "--arch %s", ""},
	"AuthFile":             {"default", "--authfile %s", ""},
	"ContainersConfModule": {"default", "--module %s", ""},
	"DNS":                  {"default", "--dns %s", ""},
	"DNSOption":            {"default", "--dns-option %s", ""},
	"DNSSearch":            {"default", "--dns-search %s", ""},
	"Environment":          {"default", "--env %s", ""},
	"File":                 {"default", "--file %s", ""},
	"ForceRM":              {"default", "--force-rm=%s", ""},
	"GlobalArgs":           {"default", "--log-level %s", ""},
	"GroupAdd":             {"default", "--group-add %s", ""},
	"ImageTag":             {"default", "--tag %s", ""},
	"Label":                {"default", "--label %s", ""},
	"Network":              {"default", "--network %s", ""},
	"PodmanArgs":           {"default", "%s", "Additional Podman arguments provided."},
	"Pull":                 {"default", "--pull=%s", ""},
	"Secret":               {"default", "--secret %s", ""},
	"Target":               {"default", "--target %s", ""},
	"TLSVerify":            {"default", "--tls-verify=%s", ""},
	"Variant":              {"default", "--variant %s", ""},
	"Volume":               {"default", "--volume %s", ""},

	// Special mappings requiring custom handling
	"SetWorkingDirectory": {"special", "", "Sets the working directory of the systemd unit file."}, // Affects systemd unit, not Podman
}

// formatBuildFlag formats a single build flag with its corresponding value.
func formatBuildFlag(key, value string) (string, error) {
	field, ok := BuildToPodman[key]
	if !ok {
		return "", fmt.Errorf("unknown key: %s", key)
	}

	if field.Type == "special" {
		// Handle special cases here, if needed
		return "", nil
	}

	return fmt.Sprintf(field.Flag, value), nil
}

// GeneratePodmanBuildCommand generates a `podman build` command from the provided options.
func GeneratePodmanBuildCommand(options map[string][]string) (string, error) {
	parts := []string{}

	for key, values := range options {
		field, ok := BuildToPodman[key]
		if !ok {
			return "", fmt.Errorf("unknown key: %s", key)
		}

		if field.Type == "special" {
			// Skip special mappings for this implementation
			continue
		}

		for _, value := range values {
			part, err := formatBuildFlag(key, value)
			if err != nil {
				return "", fmt.Errorf("error formatting build flag '%s': %w", key, err)
			}
			parts = append(parts, part)
		}
	}

	return fmt.Sprintf("podman build %s", strings.Join(parts, " ")), nil
}
