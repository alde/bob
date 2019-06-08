package cmd

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/alde/bob/config"
	"github.com/alde/bob/docker"
)

var rootCmd = &cobra.Command{
	Use:   "bob <task>",
	Short: "bob is a tool to run the recommended tools based on the project you're working on",
	Long: `bob will run the recommended tool based
	 on the project you're in, and it will do this using docker.

	 Examples:
		 bob test -> finds pom.xml -> docker run --rm -v $(pwd):/workdir -w /workdir -v ${HOME}:${HOME} mvn clean verify`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must provide a task")
		}
		task := args[0]
		cfg, err := config.New()
		if err != nil {
			return err
		}

		projectConfig, err := cfg.GetProjectConfig()
		if _, ok := projectConfig.Commands[task]; !ok {
			return fmt.Errorf("unfortunately bob doesn't know what to do about %s for projects of type %s", task, projectConfig.ProjectType)
		}
		a := docker.Command(projectConfig, task)
		command := exec.Command("docker", a...)
		runCmd(command)
		return nil
	},
}

// Execute the root command
func Execute() {
	rootCmd.Execute()
}

func runCmd(cmd *exec.Cmd) {
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	logrus.Infof("Running %s\n", strings.Join(cmd.Args, " "))
	go print(stdout)
	go print(stderr)
	cmd.Run()
}

func print(r io.Reader) {
	buf := make([]byte, 80)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			fmt.Print(string(buf[0:n]))
		}
		if err != nil {
			break
		}
	}
}
