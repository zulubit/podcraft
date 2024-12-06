package commandmap

import (
	"fmt"
	"strings"
)

type QuadletField struct {
	Type        string // "default" or "special"
	Flag        string // Podman run flag or empty for special handling
	Description string // Description of the field
	Section     string // "start", "middle", or "end"
}

var QuadletToPodman = map[string]QuadletField{
	"AddCapability":         {"default", "--cap-add %s", "", "middle"},
	"AddDevice":             {"default", "--device %s", "", "middle"},
	"AddHost":               {"default", "--add-host %s", "", "middle"},
	"Annotation":            {"default", "--annotation %s", "", "middle"},
	"AutoUpdate":            {"special", "--label io.containers.autoupdate=%s", "Adds the io.containers.autoupdate label", "middle"},
	"CgroupsMode":           {"default", "--cgroups %s", "", "middle"},
	"ContainerName":         {"default", "--name %s", "", "middle"},
	"ContainersConfModule":  {"default", "--module %s", "", "middle"},
	"DNS":                   {"default", "--dns %s", "", "middle"},
	"DNSOption":             {"default", "--dns-option %s", "", "middle"},
	"DNSSearch":             {"default", "--dns-search %s", "", "middle"},
	"DropCapability":        {"default", "--cap-drop %s", "", "middle"},
	"Entrypoint":            {"default", "--entrypoint %s", "", "middle"},
	"Environment":           {"default", "--env %s", "", "middle"},
	"EnvironmentFile":       {"default", "--env-file %s", "", "middle"},
	"EnvironmentHost":       {"default", "--env-host", "", "middle"},
	"ExposeHostPort":        {"default", "--expose %s", "", "middle"},
	"GIDMap":                {"default", "--gidmap %s", "", "middle"},
	"Group":                 {"default", "--user UID:%s", "", "middle"},
	"GroupAdd":              {"default", "--group-add %s", "", "middle"},
	"HealthCmd":             {"default", "--health-cmd %s", "", "middle"},
	"HealthInterval":        {"default", "--health-interval %s", "", "middle"},
	"HealthLogDestination":  {"default", "--health-log-destination %s", "", "middle"},
	"HealthMaxLogCount":     {"default", "--health-max-log-count %s", "", "middle"},
	"HealthMaxLogSize":      {"default", "--health-max-log-size %s", "", "middle"},
	"HealthOnFailure":       {"default", "--health-on-failure %s", "", "middle"},
	"HealthRetries":         {"default", "--health-retries %s", "", "middle"},
	"HealthStartPeriod":     {"default", "--health-start-period=period=%s", "", "middle"},
	"HealthStartupCmd":      {"default", "--health-startup-cmd %s", "", "middle"},
	"HealthStartupInterval": {"default", "--health-startup-interval %s", "", "middle"},
	"HealthStartupRetries":  {"default", "--health-startup-retries %s", "", "middle"},
	"HealthStartupSuccess":  {"default", "--health-startup-success %s", "", "middle"},
	"HealthStartupTimeout":  {"default", "--health-startup-timeout %s", "", "middle"},
	"HealthTimeout":         {"default", "--health-timeout %s", "", "middle"},
	"HostName":              {"default", "--hostname %s", "", "middle"},
	"IP":                    {"default", "--ip %s", "", "middle"},
	"IP6":                   {"default", "--ip6 %s", "", "middle"},
	"Label":                 {"default", "--label %s", "", "middle"},
	"LogDriver":             {"default", "--log-driver %s", "", "middle"},
	"LogOpt":                {"default", "--log-opt path=%s", "", "middle"},
	"Mask":                  {"default", "--security-opt mask=%s", "", "middle"},
	"Mount":                 {"default", "--mount %s", "", "middle"},
	"Network":               {"default", "--network %s", "", "middle"},
	"NetworkAlias":          {"default", "--network-alias %s", "", "middle"},
	"NoNewPrivileges":       {"default", "--security-opt no-new-privileges", "", "middle"},
	"Notify":                {"default", "--sdnotify %s", "", "middle"},
	"PidsLimit":             {"default", "--pids-limit %s", "", "middle"},
	"Pod":                   {"special", "--pod=%s", "Associates the container with a pod.", "start"},
	"PodmanArgs":            {"special", "%s", "Additional Podman arguments provided.", "middle"},
	"PublishPort":           {"default", "--publish %s", "", "middle"},
	"Pull":                  {"default", "--pull %s", "", "middle"},
	"ReadOnly":              {"default", "--read-only", "", "middle"},
	"ReadOnlyTmpfs":         {"default", "--read-only-tmpfs", "", "middle"},
	"Rootfs":                {"default", "--rootfs %s", "", "middle"},
	"RunInit":               {"default", "--init", "", "middle"},
	"SeccompProfile":        {"default", "--security-opt seccomp=%s", "", "middle"},
	"Secret":                {"default", "--secret %s", "", "middle"},
	"SecurityLabelDisable":  {"default", "--security-opt label=disable", "", "middle"},
	"SecurityLabelFileType": {"default", "--security-opt label=filetype:%s", "", "middle"},
	"SecurityLabelLevel":    {"default", "--security-opt label=level:%s", "", "middle"},
	"SecurityLabelNested":   {"default", "--security-opt label=nested", "", "middle"},
	"SecurityLabelType":     {"default", "--security-opt label=type:%s", "", "middle"},
	"ShmSize":               {"default", "--shm-size %s", "", "middle"},
	"StopSignal":            {"default", "--stop-signal %s", "", "middle"},
	"StopTimeout":           {"default", "--stop-timeout %s", "", "middle"},
	"SubGIDMap":             {"default", "--subgidname %s", "", "middle"},
	"SubUIDMap":             {"default", "--subuidname %s", "", "middle"},
	"Sysctl":                {"default", "--sysctl %s", "", "middle"},
	"Timezone":              {"default", "--tz %s", "", "middle"},
	"Tmpfs":                 {"default", "--tmpfs %s", "", "middle"},
	"UIDMap":                {"default", "--uidmap %s", "", "middle"},
	"Ulimit":                {"default", "--ulimit %s", "", "middle"},
	"Unmask":                {"default", "--security-opt unmask=%s", "", "middle"},
	"User":                  {"default", "--user %s", "", "middle"},
	"UserNS":                {"default", "--userns %s", "", "middle"},
	"Volume":                {"default", "--volume %s", "", "middle"},
	"WorkingDir":            {"default", "--workdir %s", "", "middle"},
	"Image":                 {"special", "%s", "Image specification", "end"},
	"Exec":                  {"special", "%s", "Command to run after the image", "end"},
}

// formatContainerFlag formats a single flag with its corresponding value.
func formatContainerFlag(key, value string) (string, error) {
	field, ok := QuadletToPodman[key]
	if !ok {
		return "", fmt.Errorf("unknown key: %s", key)
	}

	return fmt.Sprintf(field.Flag, value), nil
}

// GeneratePodmanContainerCommand generates a `podman run` command from the provided options.
func GeneratePodmanContainerCommand(options map[string][]string) (string, error) {
	startParts := []string{}
	middleParts := []string{}
	endParts := []string{}

	for key, values := range options {
		field, ok := QuadletToPodman[key]
		if !ok {
			return "", fmt.Errorf("unknown key: %s", key)
		}

		for _, value := range values {
			part := fmt.Sprintf(field.Flag, value)
			switch field.Section {
			case "start":
				startParts = append(startParts, part)
			case "middle":
				middleParts = append(middleParts, part)
			case "end":
				endParts = append(endParts, part)
			default:
				return "", fmt.Errorf("invalid section for key %s: %s", key, field.Section)
			}
		}
	}

	return fmt.Sprintf(
		"podman create %s %s %s",
		strings.Join(startParts, " "),
		strings.Join(middleParts, " "),
		strings.Join(endParts, " "),
	), nil
}
