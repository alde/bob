package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// ProjectConfig holds the structure of a command that bob will run
type ProjectConfig struct {
	ProjectType string              `yaml:"projectType"`
	Identifier  string              `yaml:"identifier"`
	DockerImage string              `yaml:"dockerImage"`
	Environment map[string]string   `yaml:"environment"`
	Volumes     map[string]string   `yaml:"volumes"`
	Commands    map[string][]string `yaml:"commands"`
}

// Config object holding the configuration for bob
type Config struct {
	Version int `yaml:"version"`

	Projects []ProjectConfig `yaml:"projects"`
}

var defaultConfig = Config{
	Version: 2,
	Projects: []ProjectConfig{
		ProjectConfig{
			ProjectType: "maven",
			Identifier:  "pom.xml",
			DockerImage: "docker.io/library/maven",
			Environment: map[string]string{
				"_JAVA_OPTIONS": "-Duser.home=@homeDir",
			},
			Volumes: map[string]string{},
			Commands: map[string][]string{
				"test":       []string{"mvn", "clean", "verify"},
				"checkstyle": []string{"mvn", "checkstyle:check"},
			},
		},
		ProjectConfig{
			ProjectType: "node",
			Identifier:  "package.json",
			DockerImage: "docker.io/library/node",
			Environment: map[string]string{},
			Volumes:     map[string]string{},
			Commands: map[string][]string{
				"test":       []string{"yarn", "test"},
				"checkstyle": []string{"yarn", "lint"},
			},
		},
	},
}

// New creates a new config object
func New() (*Config, error) {
	cfg, err := loadDefaultConfig()
	if err != nil {
		return nil, err
	}

	err = loadGlobalConfig(cfg)
	if err != nil {
		return nil, err
	}

	err = loadLocalConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, err
}

// GetProjectConfig returns the config specific to the project of the current dir
func (c *Config) GetProjectConfig() (*ProjectConfig, error) {
	dir, _ := os.Getwd()
	for _, config := range c.Projects {
		if exists(dir, config.Identifier) {
			logrus.Infof("project identified as %s due to precense of %s", config.ProjectType, config.Identifier)
			return &config, nil
		}
	}

	return nil, errors.New("no configuration found for the current project")
}
func loadGlobalConfig(config *Config) error {
	homedir, _ := os.UserHomeDir()
	configFile := path.Join(homedir, ".config", "bob.yaml")
	if _, err := os.Stat(configFile); err != nil {
		return nil
	}
	logrus.Debug("found global config file, using it instead of default config")

	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	return nil
}

func loadLocalConfig(config *Config) error {
	workdir, _ := os.Getwd()
	configFile := path.Join(workdir, ".bob.yaml")
	if _, err := os.Stat(configFile); err != nil {
		return nil
	}
	logrus.Debug("found local config file, using it instead of default config")

	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	return nil
}

func loadDefaultConfig() (*Config, error) {
	homedir, _ := os.UserHomeDir()
	configFile := path.Join(homedir, ".config", "bob_default.yaml")
	if _, err := os.Stat(configFile); err != nil {
		writeDefaultConfig()
	}

	var cfg Config
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	content, _ := ioutil.ReadAll(file)
	yaml.Unmarshal(content, &cfg)
	// Version mismatch of default config, so re-write it and read it again
	if cfg.Version < defaultConfig.Version {
		logrus.Debug("version mismatch in default config, updating it")
		writeDefaultConfig()
		return loadDefaultConfig()
	}

	return &cfg, nil

}
func writeDefaultConfig() {
	homedir, _ := os.UserHomeDir()
	configFile := path.Join(homedir, ".config", "bob_default.yaml")

	dir, _ := filepath.Split(configFile)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		logrus.Fatalf("unable to create directory\n%+v", err)
	}
	f, err := os.Create(configFile)
	if err != nil {
		logrus.Fatalf("unable to create file %s\n%+v", configFile, err)
	}
	defer f.Close()
	e := yaml.NewEncoder(f)
	if err = e.Encode(defaultConfig); err != nil {
		logrus.Fatalf("unable to write configuration\n%+v", err)
	}
	logrus.WithFields(logrus.Fields{
		"config": configFile,
	}).Debug("wrote default config")
}

func exists(pwd, file string) bool {
	fullpath := path.Join(pwd, file)
	if _, err := os.Stat(fullpath); err == nil {
		return true
	}
	return false
}
