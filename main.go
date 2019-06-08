package main

import (
	"github.com/sirupsen/logrus"

	"github.com/alde/bob/cmd"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	cmd.Execute()
}
