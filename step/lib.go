package step

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bitrise-io/go-utils/errorutil"
	"github.com/bitrise-io/go-utils/retry"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-xcode/v2/destination"
	"github.com/bitrise-steplib/steps-xcode-test/simulator"
)

func (s SimulatorStarter) getSimulatorForDestination(destinationSpecifier string) (simulator.Simulator, error) {
	var sim simulator.Simulator
	var osVersion string

	simulatorDestination, err := destination.NewSimulator(destinationSpecifier)
	if err != nil {
		return simulator.Simulator{}, fmt.Errorf("invalid destination specifier (%s): %w", destinationSpecifier, err)
	}

	platform := strings.TrimSuffix(simulatorDestination.Platform, " Simulator")
	// Retry gathering device information since xcrun simctl list can fail to show the complete device list
	if err := retry.Times(3).Wait(10 * time.Second).Try(func(attempt uint) error {
		var errGetSimulator error
		if simulatorDestination.OS == "latest" {
			simulatorDevice := simulatorDestination.Name
			if simulatorDevice == "iPad" {
				s.logger.Warnf("Given device (%s) is deprecated, using iPad Air (3rd generation)...", simulatorDevice)
				simulatorDevice = "iPad Air (3rd generation)"
			}

			sim, osVersion, errGetSimulator = s.simulatorManager.GetLatestSimulatorAndVersion(platform, simulatorDevice)
		} else {
			normalizedOsVersion := simulatorDestination.OS
			osVersionSplit := strings.Split(normalizedOsVersion, ".")
			if len(osVersionSplit) > 2 {
				normalizedOsVersion = strings.Join(osVersionSplit[0:2], ".")
			}
			osVersion = fmt.Sprintf("%s %s", platform, normalizedOsVersion)

			sim, errGetSimulator = s.simulatorManager.GetSimulator(osVersion, simulatorDestination.Name)
		}

		if errGetSimulator != nil {
			s.logger.Warnf("attempt %d to get simulator UDID failed with error: %s", attempt, errGetSimulator)
		}

		return errGetSimulator
	}); err != nil {
		return simulator.Simulator{}, fmt.Errorf("simulator UDID lookup failed: %w", err)
	}

	s.logger.Infof("Simulator infos")
	s.logger.Printf("* simulator_name: %s, version: %s, UDID: %s, status: %s", sim.Name, osVersion, sim.ID, sim.Status)

	return sim, nil
}

func (s SimulatorStarter) WaitForSimulatorBoot(id string) error {
	const timeout = time.Second * 180

	timer := time.NewTimer(timeout)
	defer func() {
		timer.Stop()
	}()

	waitCmd := s.commandFactory.Create("xcrun", []string{"simctl", "launch", id, "com.apple.mobileslideshow"}, &command.Opts{
		Stdout: os.Stderr,
		Stderr: os.Stderr,
	})
	launchDoneCh := make(chan error, 1)

	doWait := func() {
		s.logger.Println()
		s.logger.TDonef("$ %s", waitCmd.PrintableCommandArgs())
		launchDoneCh <- waitCmd.Run()
	}
	go doWait()

	for {
		select {
		case err := <-launchDoneCh:
			{
				if err != nil { // error or timeout
					if errorutil.IsExitStatusError(err) {
						s.logger.Warnf("launch failed, restarting")
						go doWait() // restart wait
						continue
					}
					return err
				}
				return nil // launch succeeded
			}
		case <-timer.C:
			return fmt.Errorf("failed to boot Simulator in %s", timeout)
		}
	}
}
