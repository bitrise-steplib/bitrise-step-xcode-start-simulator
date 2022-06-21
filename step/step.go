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
	Erase       bool   `env:"erase,required"`
	WaitForBoot bool   `env:"wait_for_boot,required"`
	// Debugging
	DebugLog           bool `env:"verbose_log,required"`
	WaitForBootTimeout int  `env:"wait_for_boot_timeout,required"`
}

type Config struct {
	Input

	SimulatorID string
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
	s.logger.EnableDebugLog(input.DebugLog)

	sim, err := s.getSimulatorForDestination(input.Destination)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Input:       input,
		SimulatorID: sim.ID,
	}, nil
}

func (s SimulatorStarter) InstallDependencies() error {
	return nil
}

func (s SimulatorStarter) Run(config Config) (Result, error) {
	err := s.prepareSimulator(config)
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

func (s SimulatorStarter) getSimulatorForDestination(destinationSpecifier string) (destination.Device, error) {
	var device destination.Device

	simulatorDestination, err := destination.NewSimulator(destinationSpecifier)
	if err != nil || simulatorDestination == nil {
		return destination.Device{}, fmt.Errorf("invalid destination specifier (%s): %w", destinationSpecifier, err)
	}

	device, err = s.deviceFinder.FindDevice(*simulatorDestination)
	if err != nil {
		return destination.Device{}, fmt.Errorf("simulator UDID lookup failed: %w", err)
	}

	s.logger.Infof("Simulator info")
	s.logger.Printf("* simulator_name: %s, version: %s, UDID: %s, status: %s", device.Name, device.OS, device.ID, device.Status)

	return device, nil
}

func (s SimulatorStarter) prepareSimulator(config Config) error {
	err := s.simulatorManager.ResetLaunchServices()
	if err != nil {
		s.logger.Warnf("Failed to apply simulator boot workaround: %s", err)
	}

	if config.Erase {
		s.logger.Donef("Erasing simulator...")
		if err := s.simulatorManager.Shutdown(config.SimulatorID); err != nil {
			return err
		}
		if err := s.simulatorManager.Erase(config.SimulatorID); err != nil {
			return err
		}
	}

	s.logger.Println()
	s.logger.TDonef("Booting simulator...")
	if err := s.simulatorManager.Boot(config.SimulatorID); err != nil {
		return err
	}

	if config.WaitForBoot {
		s.logger.Println()
		s.logger.TDonef("Waiting for simulator to boot...")

		timeout := time.Duration(config.WaitForBootTimeout) * time.Second
		if err := s.simulatorManager.WaitForBootFinished(config.SimulatorID, timeout); err != nil {
			return err
		}

		s.logger.Println()
		s.logger.TDonef("Successfully started simulator.")

		s.logger.Infof("Enabling simulator verbose log for better diagnostics")
		if err := s.simulatorManager.EnableVerboseLog(config.SimulatorID); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}
