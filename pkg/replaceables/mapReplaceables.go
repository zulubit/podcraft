package replaceables

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"regexp"

	"github.com/zulubit/podcraft/pkg/configfile"
)

type ReplaceablesMap map[string]string

func mapReplaceables(config configfile.Config, isProd bool) (*ReplaceablesMap, error) {
	replaceables := make(ReplaceablesMap)

	for _, r := range config.Replaceables {
		if r.Id == "" {
			return nil, errors.New("Replaceable with empty id is not allowed")
		}

		re := regexp.MustCompile(`^[a-z]+$`)
		match := re.MatchString(r.Id)
		if !match {
			return nil, errors.New("Replaceable Id can only be lowercase letters with no spaces")
		}

		var val string
		if isProd {
			val = r.Prod
		} else {
			val = r.Dev
		}

		replaceables[fmt.Sprint(r.Id)] = val
	}

	return &replaceables, nil
}

func ReplaceReplaceables(rawToml string, config configfile.Config, isProd bool) (*string, error) {

	replaceables, err := mapReplaceables(config, isProd)
	if err != nil {
		return nil, err
	}

	// Create a new template with custom delimiters
	tmpl := template.New("example").Delims("<<", ">>")

	// Parse the template string
	tmpl, err = tmpl.Parse(rawToml)
	if err != nil {
		return nil, err
	}

	var replacedTemplate bytes.Buffer

	err = tmpl.Execute(&replacedTemplate, replaceables)
	if err != nil {
		return nil, err
	}

	newConfig := replacedTemplate.String()

	return &newConfig, nil
}
