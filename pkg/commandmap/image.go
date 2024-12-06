package commandmap

import (
	"fmt"
	"strings"
)

type ImageField struct {
	Type        string // "default" or "special"
	Flag        string // Podman image flag or empty for special handling
	Description string // Description of the field
}

var ImageToPodman = map[string]ImageField{
	// Default mappings
	"AllTags":              {"default", "--all-tags", ""},
	"Arch":                 {"default", "--arch %s", ""},
	"AuthFile":             {"default", "--authfile %s", ""},
	"CertDir":              {"default", "--cert-dir %s", ""},
	"ContainersConfModule": {"default", "--module %s", ""},
	"Creds":                {"default", "--creds %s", ""},
	"DecryptionKey":        {"default", "--decryption-key %s", ""},
	"GlobalArgs":           {"default", "--log-level %s", ""},
	"OS":                   {"default", "--os %s", ""},
	"TLSVerify":            {"default", "--tls-verify=%s", ""},
	"Variant":              {"default", "--variant %s", ""},

	// Special mappings requiring custom handling
	"Image":      {"special", "", "Specifies the image to pull."},      // Required positional argument
	"ImageTag":   {"special", "", "Used to resolve image references."}, // May affect other parts of the command
	"PodmanArgs": {"special", "", "Direct arguments to pass to Podman."},
}

// formatImageFlag formats a single image-related flag with its corresponding value.
func formatImageFlag(key, value string) (string, error) {
	field, ok := ImageToPodman[key]
	if !ok {
		return "", fmt.Errorf("unknown key: %s", key)
	}

	if field.Type == "special" {
		// Handle special cases without returning an error
		return "", nil
	}

	if value == "" && !strings.Contains(field.Flag, "%s") {
		return field.Flag, nil // Boolean flags like --all-tags
	}

	return fmt.Sprintf(field.Flag, value), nil
}

// GeneratePodmanImageCommand generates a `podman image pull` command from the provided options.
func GeneratePodmanImageCommand(options map[string][]string) (string, error) {
	parts := []string{}
	var image string

	for key, values := range options {
		field, ok := ImageToPodman[key]
		if !ok {
			return "", fmt.Errorf("unknown key: %s", key)
		}

		if field.Type == "special" {
			// Handle "Image" as a required positional argument
			if key == "Image" {
				if len(values) > 0 {
					image = values[0] // Use the first value for the image
				}
				continue
			}
			// Skip other special mappings (like ImageTag) if not directly needed
			continue
		}

		for _, value := range values {
			part, err := formatImageFlag(key, value)
			if err != nil {
				return "", fmt.Errorf("error formatting image flag '%s': %w", key, err)
			}
			parts = append(parts, part)
		}
	}

	if image == "" {
		return "", fmt.Errorf("missing required field: Image\n")
	}

	return fmt.Sprintf(
		"podman image pull %s %s",
		strings.Join(parts, " "),
		image,
	), nil
}
