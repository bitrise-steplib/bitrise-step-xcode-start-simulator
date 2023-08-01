package step

import (
	"errors"
	"testing"
	"time"

	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-xcode/v2/destination"
	"github.com/bitrise-steplib/bitrise-step-xcode-start-simulator/step/mocks"
	"github.com/stretchr/testify/require"
)

func Test_GivenBootOnlyConfig_WhenBoot_ThenSuccessfullyBoots(t *testing.T) {
	// Given
	const udid = "test-ID"

	var (
		logger           = log.NewLogger()
		simulatorManager = mocks.NewSimulatorManager(t)
		s                = SimulatorStarter{
			logger:           logger,
			simulatorManager: simulatorManager,
		}
		simulator = destination.Device{
			ID:       udid,
			Platform: "iOS Simulator",
			Name:     "Bitrise iOS default",
			OS:       "11",
		}
		config = Config{
			Input: Input{
				WaitForBootTimeout: 0,
			},
			Simulator: simulator,
		}
	)

	simulatorManager.On("ResetLaunchServices").Once().Return(nil)
	simulatorManager.On("Boot", simulator).Once().Return(nil)

	// When
	got, err := s.Run(config)

	// Then
	require.NoError(t, err)
	require.Equal(t, got, Result{
		SimulatorStatus: "booted",
		Destination:     "platform=iOS Simulator,name=Bitrise iOS default,OS=11",
	})
}

func Test_GivenBootOnlyConfig_WhenSimulatorBootFails_ThenItReturnsError(t *testing.T) {
	// Given
	const udid = "test-ID"

	var (
		logger           = log.NewLogger()
		simulatorManager = mocks.NewSimulatorManager(t)
		s                = SimulatorStarter{
			logger:           logger,
			simulatorManager: simulatorManager,
		}
		simulator = destination.Device{
			ID:       udid,
			Platform: "iOS Simulator",
			Name:     "Bitrise iOS default",
			OS:       "11",
		}
		config = Config{
			Input: Input{
				WaitForBootTimeout: 0,
			},
			Simulator: simulator,
		}
	)

	simulatorManager.On("ResetLaunchServices").Once().Return(nil)
	simulatorManager.On("Boot", simulator).Once().Return(errors.New("boot error"))

	// When
	got, err := s.Run(config)

	// Then
	require.Error(t, err)
	require.Equal(t, got, Result{
		SimulatorStatus: "failed",
		Destination:     "platform=iOS Simulator,name=Bitrise iOS default,OS=11",
	})
}

func Test_GivenWaitForBootConfig_WhenWaitForBootFails_ThenReturnsTimeoutError(t *testing.T) {
	// Given
	const (
		udid    = "test-ID"
		timeout = 1 * time.Second
	)

	var (
		logger           = log.NewLogger()
		simulatorManager = mocks.NewSimulatorManager(t)
		s                = SimulatorStarter{
			logger:           logger,
			simulatorManager: simulatorManager,
		}
		simulator = destination.Device{
			ID:       udid,
			Platform: "iOS Simulator",
			Name:     "Bitrise iOS default",
			OS:       "11",
		}
		config = Config{
			Input: Input{
				WaitForBootTimeout: int(timeout.Seconds()),
			},
			Simulator: simulator,
		}
	)

	simulatorManager.On("ResetLaunchServices").Once().Return(nil)
	simulatorManager.On("Boot", simulator).Once().Return(nil)
	simulatorManager.On("WaitForBootFinished", udid, timeout).Once().Return(errors.New("timeout"))

	// When
	got, err := s.Run(config)

	// Then
	require.Error(t, err)
	require.Equal(t, got, Result{
		SimulatorStatus: "hanged",
		Destination:     "platform=iOS Simulator,name=Bitrise iOS default,OS=11",
	})
}
