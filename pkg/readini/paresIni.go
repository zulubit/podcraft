package readini

import (
	"fmt"

	"gopkg.in/ini.v1"
)

// HARD GPT fyi
// ReadDataFromIni parses INI-style data into a map[string]map[string][]string.
// Handles sections, including the default section, and supports multiple values for the same key using ShadowLoad.
func ReadDataFromIni(data string) (map[string]map[string][]string, error) {
	// Create a new map to hold the parsed data with sections
	parsedData := make(map[string]map[string][]string)

	// Load the data using ShadowLoad to preserve multiple values for the same key
	cfg, err := ini.ShadowLoad([]byte(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse INI data: %w", err)
	}

	// Iterate through all sections in the file
	for _, section := range cfg.Sections() {
		sectionName := section.Name()
		if sectionName == ini.DefaultSection {
			sectionName = "default" // Rename default section for clarity
		}

		// Create a map for the current section if it doesn't exist
		sectionData := make(map[string][]string)

		// Iterate over all keys in the section
		for _, key := range section.Keys() {
			// Use ValueWithShadows to retrieve all values for the key
			sectionData[key.Name()] = key.ValueWithShadows()
		}

		// Add the section data to the parsedData map
		parsedData[sectionName] = sectionData
	}

	return parsedData, nil
}
