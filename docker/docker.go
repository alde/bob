package docker

import (
	"fmt"
	"os"
	"strings"

	"github.com/alde/bob/config"
)

// some global variables (yuck)
var (
	homedir, _          = os.UserHomeDir()
	workingDirectory, _ = os.Getwd()
)

// Command transforms the target command into a docker run command
func Command(dc *config.ProjectConfig, target string) []string {
	envs := assembleEnvs(dc)
	volumes := assembleVolumes(dc)

	arguments := []string{
		"run", "--rm",
		"-w", "/workdir",
	}
	for key, value := range envs {
		arguments = append(arguments, "-e", fmt.Sprintf("%s=%s", key, value))
	}
	for from, to := range volumes {
		arguments = append(arguments, "-v", fmt.Sprintf("%s:%s", from, to))
	}

	arguments = append(arguments, dc.DockerImage)
	arguments = append(arguments, dc.Commands[target]...)

	return arguments
}

func modify(value string) string {
	value = strings.Replace(value, "@homeDir", homedir, -1)

	return value
}

// assemble the environment variables needed to run the command
func assembleEnvs(dc *config.ProjectConfig) map[string]string {
	envs := make(map[string]string)

	// Set some defaults
	envs["HOME"] = homedir

	// Add configured environments (overriding any defaults)
	// It will run through the modifier first in order to expand any '@<var>'s
	for key, value := range dc.Environment {
		modifiedValue := modify(value)
		envs[key] = modifiedValue
	}

	return envs
}

// assemble the volumes needed to run the command
func assembleVolumes(dc *config.ProjectConfig) map[string]string {
	volumes := make(map[string]string)
	// Set some defaults
	volumes[workingDirectory] = "/workdir"
	volumes[homedir] = homedir

	// Add configured volumes (overriding any defaults)
	// It will run through the modifier first in order to expand any '@<var>'s
	for key, value := range dc.Volumes {
		modifiedKey := modify(key)
		modifiedValue := modify(value)
		volumes[modifiedKey] = modifiedValue
	}

	return volumes
}
