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

const (
	simulatorResultStatusSuccess = "booted"
	simulatorResultStatusFailed  = "failed"
	simulatorResultStatusHanged  = "hanged"
)

var errTimeout = errors.New("simulator boot timed out")

type Input struct {
	Destination string `env:"destination,required"`
	// Debugging
	IsVerboseLog       bool `env:"verbose_log,opt[yes,no]"`
	WaitForBootTimeout int  `env:"wait_for_boot_timeout,required"`
	ShouldReset        bool `env:"reset,opt[yes,no]"`
}

type Config struct {
	Input

	Simulator destination.Device
}

type Result struct {
	SimulatorStatus string
	Destination     string
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
	s.logger.EnableDebugLog(input.IsVerboseLog)

	sim, err := s.getSimulatorForDestination(input.Destination)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Input:     input,
		Simulator: sim,
	}, nil
}

func (s SimulatorStarter) InstallDependencies() error {
	return nil
}

func (s SimulatorStarter) Run(config Config) (Result, error) {
	err := s.prepareSimulator(config.Simulator, config.WaitForBootTimeout, config.ShouldReset)

	simulatorStatus := simulatorResultStatusSuccess
	if err != nil {
		if errors.Is(err, errTimeout) {
			simulatorStatus = simulatorResultStatusHanged
		} else {
			simulatorStatus = simulatorResultStatusFailed
		}
	}

	// Omitting `arch` option on purpose, to enable xcode-test Step pick up prebooted Rosetta Simulator
	exportedDestination := fmt.Sprintf("platform=%s,name=%s,OS=%s", config.Simulator.Platform, config.Simulator.Name, config.Simulator.OS)
	return Result{
		SimulatorStatus: simulatorStatus,
		Destination:     exportedDestination,
	}, err
}

func (s SimulatorStarter) ExportOutputs(result Result) error {
	const (
		simulatorStatusKey = "BITRISE_SIMULATOR_STATUS"
		destinationKey     = "BITRISE_XCODE_DESTINATION"
	)

	s.logger.Println()
	s.logger.Donef("Exporting ouputs")

	s.logger.Infof("Output %s = %s", simulatorStatusKey, result.SimulatorStatus)
	if err := s.stepenvRepository.Set(simulatorStatusKey, result.SimulatorStatus); err != nil {
		return err
	}

	s.logger.Infof("Output %s = %s", destinationKey, result.Destination)
	if err := s.stepenvRepository.Set(destinationKey, result.Destination); err != nil {
		return err
	}

	return nil
}

func (s SimulatorStarter) prepareSimulator(simulator destination.Device, waitForBootTimeout int, shouldReset bool) error {
	err := s.simulatorManager.ResetLaunchServices()
	if err != nil {
		s.logger.Warnf("Failed to apply simulator boot workaround: %s", err)
	}

	UDID := simulator.ID
	if shouldReset {
		s.logger.Println()
		s.logger.Donef("Erasing simulator...")
		if err := s.simulatorManager.Shutdown(UDID); err != nil {
			return err
		}
		if err := s.simulatorManager.Erase(UDID); err != nil {
			return err
		}
	}

	s.logger.Println()
	s.logger.TDonef("Booting simulator...")
	if err := s.simulatorManager.Boot(simulator); err != nil {
		return err
	}

	if waitForBootTimeout > 0 {
		s.logger.Println()
		s.logger.TDonef("Waiting for the simulator to finish booting...")

		if simulator.Platform == string(destination.VisionOSSimulator) {
			// Xcode 15.1 Beta 3: `xcrun simctl launch booted com.apple.Preferences` does not work on visionOS.
			// Error:
			// The system shell probably crashed.
			// Underlying error (domain=FBSOpenApplicationServiceErrorDomain, code=5):
			// The request to open "com.apple.Preferences" failed. The system shell probably crashed.
			s.logger.Warnf("Detecting boot completion is not working with visionOS yet, skipping...")
		} else {
			timeout := time.Duration(waitForBootTimeout) * time.Second
			if err := s.simulatorManager.WaitForBootFinished(UDID, timeout); err != nil {
				s.logger.Errorf("%s", err)
				return errTimeout
			}
		}

		s.logger.Println()
		s.logger.TDonef("Successfully started simulator.")
	} else {
		s.logger.Printf("Not waiting for the simulator to finish booting (timeout not set).")
	}

	return nil
}
