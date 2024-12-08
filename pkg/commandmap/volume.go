package commandmap

import (
	"fmt"
	"strings"
)

type VolumeField struct {
	Type        string // "default" or "special"
	Flag        string // Podman volume flag or empty for special handling
	Description string // Description of the field
}

var VolumeToPodman = map[string]VolumeField{
	// Default mappings
	"ContainersConfModule": {"default", "--module %s", ""},
	"Copy":                 {"default", "--opt copy", ""},
	"Device":               {"default", "--opt device=%s", ""},
	"Driver":               {"default", "--driver %s", ""},
	"GlobalArgs":           {"default", "--log-level %s", ""},
	"Group":                {"default", "--opt group=%s", ""},
	"Image":                {"default", "--opt image=%s", ""},
	"Label":                {"default", "--label %s", ""},
	"Options":              {"default", "--opt o=%s", ""},
	"PodmanArgs":           {"default", "%s", "Additional Podman arguments provided."},
	"Type":                 {"default", "--opt type=%s", ""},
	"User":                 {"default", "--opt uid=%s", ""},

	// Special mappings requiring custom handling
	"VolumeName": {"special", "", "Specifies the name of the volume."}, // Required positional argument
}

// formatVolumeFlag formats a single volume-related flag with its corresponding value.
func formatVolumeFlag(key, value string) (string, error) {
	field, ok := VolumeToPodman[key]
	if !ok {
		return "", fmt.Errorf("unknown key: %s", key)
	}

	if field.Type == "special" {
		// Special fields are handled separately
		return "", nil
	}

	if value == "" && !strings.Contains(field.Flag, "%s") {
		return field.Flag, nil // Boolean flags like --opt copy
	}

	return fmt.Sprintf(field.Flag, value), nil
}

// GeneratePodmanVolumeCommand generates a `podman volume create` command from the provided options.
func GeneratePodmanVolumeCommand(name string, options map[string][]string) (string, error) {
	parts := []string{}
	var volumeName string

	for key, values := range options {
		field, ok := VolumeToPodman[key]
		if !ok {
			return "", fmt.Errorf("unknown key: %s", key)
		}

		if field.Type == "special" {
			// Handle "VolumeName" as a required positional argument
			if key == "VolumeName" {
				if len(values) > 0 {
					volumeName = values[0] // Use the first value of VolumeName
				}
				continue
			}
			continue // Skip other special mappings for now
		}

		for _, value := range values {
			part, err := formatVolumeFlag(key, value)
			if err != nil {
				return "", err
			}
			parts = append(parts, part)
		}
	}

	//
	if volumeName == "" {
		volumeName = name + ".volume"
	}

	return fmt.Sprintf(
		"podman volume exists %s || podman volume create %s %s",
		volumeName,
		strings.Join(parts, " "),
		volumeName,
	), nil
}
