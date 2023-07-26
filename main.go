package main

import (
	"os"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-steputils/v2/stepenv"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-xcode/v2/destination"
	"github.com/bitrise-io/go-xcode/v2/simulator"
	"github.com/bitrise-io/go-xcode/v2/xcodeversion"
	"github.com/bitrise-steplib/bitrise-step-xcode-start-simulator/step"
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
	exportErr := simulatorStarter.ExportOutputs(result)

	exitCode := 0
	if runErr != nil {
		logger.Errorf("Run: %s", runErr)
		exitCode = 1
	}

	if exportErr != nil {
		logger.Errorf("Export outputs: %s", exportErr)
		exitCode = 1
	}

	return exitCode
}

func createStep(logger log.Logger) step.SimulatorStarter {
	envRepository := env.NewRepository()
	inputParser := stepconf.NewInputParser(envRepository)
	commandFactory := command.NewFactory(envRepository)
	xcodebuildVersionProvider := xcodeversion.NewXcodeVersionProvider(commandFactory)
	xcodeVersion, err := xcodebuildVersionProvider.GetVersion()
	if err != nil { // not fatal error, continuing with empty version
		logger.Errorf("failed to read Xcode version: %s", err)
	}
	deviceFinder := destination.NewDeviceFinder(logger, commandFactory, xcodeVersion)
	simulatorManager := simulator.NewManager(logger, commandFactory)
	stepenvRepository := stepenv.NewRepository(envRepository)

	return step.NewStep(logger, inputParser, stepenvRepository, commandFactory, deviceFinder, simulatorManager)
}
