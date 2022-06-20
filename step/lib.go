package step

import (
	"fmt"

	"github.com/bitrise-io/go-xcode/v2/destination"
)

func (s SimulatorStarter) getSimulatorForDestination(destinationSpecifier string) (destination.Device, error) {
	var device destination.Device

	simulatorDestination, err := destination.NewSimulator(destinationSpecifier)
	if err != nil {
		return destination.Device{}, fmt.Errorf("invalid destination specifier (%s): %w", destinationSpecifier, err)
	}

	device, err = s.deviceFinder.FindDevice(*simulatorDestination)
	if err != nil {
		return destination.Device{}, fmt.Errorf("simulator UDID lookup failed: %w", err)
	}

	s.logger.Infof("Simulator infos")
	s.logger.Printf("* simulator_name: %s, version: %s, UDID: %s, status: %s", device.Name, device.OS, device.ID, device.Status)

	return device, nil
}
