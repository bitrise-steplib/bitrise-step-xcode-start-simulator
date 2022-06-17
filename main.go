package main

import (
	"os"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/bitrise-step-xcode-start-simulator/step"
	"github.com/bitrise-steplib/steps-xcode-test/simulator"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := log.NewLogger()
	simulatorStarter := createStep(logger)

	config, err := simulatorStarter.ProcessConfig()
	if err != nil {
		logger.Errorf("Process config: %s", err)
		return 1
	}

	result, runErr := simulatorStarter.Run(config)
	exportErr := simulatorStarter.ExportOtputs(result)

	if runErr != nil {
		logger.Errorf("Run: %s", err)
		return 1
	}

	if exportErr != nil {
		logger.Errorf("Export outputs: %s", err)
		return 1
	}

	return 0
}

func createStep(logger log.Logger) step.SimulatorStarter {
	envRepository := env.NewRepository()
	inputParser := stepconf.NewInputParser(envRepository)
	commandFactory := command.NewFactory(envRepository)
	simulatorManager := simulator.NewManager(commandFactory)

	return step.NewStep(logger, inputParser, commandFactory, simulatorManager)
}
