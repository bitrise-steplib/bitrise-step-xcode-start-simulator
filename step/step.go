package step

import (
	"fmt"
	"time"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-xcode/v2/destination"
	"github.com/bitrise-io/go-xcode/v2/simulator"
)

type Input struct {
	Destination string `env:"destination,required"`
}

type Config struct {
	SimulatorID       string
	IsSimulatorBooted bool
}

type Result struct {
	IsSimulatorTimeout bool
}

type SimulatorStarter struct {
	logger            log.Logger
	inputParser       stepconf.InputParser
	stepenvRepository env.Repository
	commandFactory    command.Factory
	deviceFinder      destination.DeviceFinder
	simulatorManager  simulator.Manager
}

func NewStep(logger log.Logger, inputParser stepconf.InputParser, stepenvRepository env.Repository, commandFactory command.Factory, deviceFinder destination.DeviceFinder, simualatorManager simulator.Manager) SimulatorStarter {
	return SimulatorStarter{
		logger:            logger,
		inputParser:       inputParser,
		commandFactory:    commandFactory,
		deviceFinder:      deviceFinder,
		simulatorManager:  simualatorManager,
		stepenvRepository: stepenvRepository,
	}
}

func (s SimulatorStarter) ProcessConfig() (Config, error) {
	var input Input
	if err := s.inputParser.Parse(&input); err != nil {
		return Config{}, err
	}

	stepconf.Print(input)
	s.logger.Println()

	sim, err := s.getSimulatorForDestination(input.Destination)
	if err != nil {
		return Config{}, err
	}

	return Config{
		SimulatorID: sim.ID,
	}, nil
}

func (s SimulatorStarter) InstallDependencies() error {
	return nil
}

func (s SimulatorStarter) Run(config Config) (Result, error) {
	err := s.prepareSimulator(true, config.SimulatorID)
	if err != nil {
		return Result{
			IsSimulatorTimeout: true,
		}, err
	}

	return Result{}, nil
}

func (s SimulatorStarter) ExportOutputs(result Result) error {
	err := s.stepenvRepository.Set("BITRISE_IS_SIMULATOR_ERROR", fmt.Sprintf("%t", result.IsSimulatorTimeout))
	if err != nil {
		return err
	}

	return nil
}

func (s SimulatorStarter) prepareSimulator(enableSimulatorVerboseLog bool, simulatorID string) error {
	err := s.simulatorManager.ResetLaunchServices()
	if err != nil {
		s.logger.Warnf("Failed to apply simulator boot workaround: %s", err)
	}

	if err := s.simulatorManager.Shutdown(simulatorID); err != nil {
		return err
	}
	if err := s.simulatorManager.Erase(simulatorID); err != nil {
		return err
	}

	s.logger.Println()
	s.logger.TDonef("Booting Simulator...")
	if err := s.simulatorManager.Boot(simulatorID); err != nil {
		return err
	}

	s.logger.Println()
	s.logger.TDonef("Waiting for simulator to boot...")
	const timeout = time.Second * 60
	if err := s.simulatorManager.WaitForBootFinished(simulatorID, timeout); err != nil {
		return err
	}

	s.logger.Println()
	s.logger.TDonef("Successfully started Simulator.")

	if enableSimulatorVerboseLog {
		s.logger.Infof("Enabling Simulator verbose log for better diagnostics")
		if err := s.simulatorManager.EnableVerboseLog(simulatorID); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}
