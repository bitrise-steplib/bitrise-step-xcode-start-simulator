package main

import (
	"os"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/bitrise-step-xcode-start-simulator/step"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := log.NewLogger()
	buildStep := createStep(logger)

	_, err := buildStep.ProcessConfig()
	if err != nil {
		logger.Errorf("Process config: %s", err)
		return 1
	}

	return 0
}

func createStep(logger log.Logger) step.SimulatorStarter {
	envRepository := env.NewRepository()
	inputParser := stepconf.NewInputParser(envRepository)

	return step.NewStep(logger, inputParser)
}
