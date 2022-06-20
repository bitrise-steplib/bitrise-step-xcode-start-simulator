package step

import (
	"fmt"
	"os"
	"time"

	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-xcode/v2/destination"
)

func (s SimulatorStarter) getSimulatorForDestination(destinationSpecifier string) (destination.Device, error) {
	var device destination.Device

	simulatorDestination, err := destination.NewSimulator(destinationSpecifier)
	if err != nil {
		return destination.Device{}, fmt.Errorf("invalid destination specifier (%s): %w", destinationSpecifier, err)
	}

	device, err = s.deviceFinder.GetSimulator(*simulatorDestination)
	if err != nil {
		return destination.Device{}, fmt.Errorf("simulator UDID lookup failed: %w", err)
	}

	s.logger.Infof("Simulator infos")
	s.logger.Printf("* simulator_name: %s, version: %s, UDID: %s, status: %s", device.Name, device.OS, device.ID, device.Status)

	return device, nil
}

func (s SimulatorStarter) WaitForSimulatorBoot(id string) error {
	const timeout = time.Second * 60

	timer := time.NewTimer(timeout)
	defer func() {
		timer.Stop()
	}()

	launchDoneCh := make(chan error, 1)
	doWait := func() {
		waitCmd := s.commandFactory.Create("xcrun", []string{"simctl", "launch", id, "com.apple.Preferences"}, &command.Opts{
			Stdout: os.Stderr,
			Stderr: os.Stderr,
		})

		s.logger.Println()
		s.logger.TDonef("$ %s", waitCmd.PrintableCommandArgs())
		launchDoneCh <- waitCmd.Run()
	}

	go doWait()

	for {
		select {
		case err := <-launchDoneCh:
			{
				if err != nil {
					return fmt.Errorf("failed to wait for simulator boot: %w", err)
				}
				return nil // launch succeeded
			}
		case <-timer.C:
			return fmt.Errorf("failed to boot Simulator in %s", timeout)
		}
	}
}
