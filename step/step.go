package step

import (
	"errors"
	"fmt"
	"time"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-xcode/v2/destination"
	"github.com/bitrise-io/go-xcode/v2/simulator"
)

var errTimeout = errors.New("simulator boot timed out")

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
	Destination        string
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

	return Result{
		IsSimulatorTimeout: errors.Is(err, errTimeout),
		Destination:        config.Destination,
	}, err
}

func (s SimulatorStarter) ExportOutputs(result Result) error {
	const (
		isSimErrorKey  = "BITRISE_IS_SIMULATOR_TIMEOUT"
		destinationKey = "BITRISE_XCODE_DESTINATION"
	)

	s.logger.Println()
	s.logger.Donef("Exporting ouputs")

	isTimeout := fmt.Sprintf("%t", result.IsSimulatorTimeout)
	s.logger.Infof("Output %s = %s", isSimErrorKey, isTimeout)
	if err := s.stepenvRepository.Set(isSimErrorKey, isTimeout); err != nil {
		return err
	}

	s.logger.Infof("Output %s = %s", destinationKey, result.Destination)
	if err := s.stepenvRepository.Set(destinationKey, result.Destination); err != nil {
		return err
	}

	return nil
}

func (s SimulatorStarter) prepareSimulator(config Config) error {
	err := s.simulatorManager.ResetLaunchServices()
	if err != nil {
		s.logger.Warnf("Failed to apply simulator boot workaround: %s", err)
	}

	if config.Erase {
		s.logger.Println()
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
			s.logger.Errorf("%s", err)
			return errTimeout
		}

		s.logger.Println()
		s.logger.TDonef("Successfully started simulator.")
	}

	return nil
}
