package step

import (
	"errors"
	"testing"
	"time"

	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/bitrise-step-xcode-start-simulator/step/mocks"
	"github.com/stretchr/testify/require"
)

func Test_GivenBootOnlyConfig_WhenBoot_ThenSuccessfullyBoots(t *testing.T) {

	// Given
	const (
		dest = "dest"
		udid = "test-ID"
	)

	var (
		logger           = log.NewLogger()
		simulatorManager = new(mocks.SimulatorManager)
		s                = SimulatorStarter{
			logger:           logger,
			simulatorManager: simulatorManager,
		}
		config = Config{
			Input: Input{
				Destination:        dest,
				WaitForBootTimeout: 0,
			},
			SimulatorID: udid,
		}
	)

	simulatorManager.On("ResetLaunchServices").Once().Return(nil)
	simulatorManager.On("Boot", udid).Once().Return(nil)

	// When
	got, err := s.Run(config)

	// Then
	require.NoError(t, err)
	require.Equal(t, got, Result{
		SimulatorStatus: "booted",
		Destination:     dest,
	})
	simulatorManager.AssertExpectations(t)
}

func Test_GivenBootOnlyConfig_WhenSimulatorBootFails_ThenItReturnsError(t *testing.T) {
	// Given
	const (
		dest = "dest"
		udid = "test-ID"
	)

	var (
		logger           = log.NewLogger()
		simulatorManager = new(mocks.SimulatorManager)
		s                = SimulatorStarter{
			logger:           logger,
			simulatorManager: simulatorManager,
		}
		config = Config{
			Input: Input{
				Destination:        dest,
				WaitForBootTimeout: 0,
			},
			SimulatorID: udid,
		}
	)

	simulatorManager.On("ResetLaunchServices").Once().Return(nil)
	simulatorManager.On("Boot", udid).Once().Return(errors.New("boot error"))

	// When
	got, err := s.Run(config)

	// Then
	require.Error(t, err)
	require.Equal(t, got, Result{
		SimulatorStatus: "failed",
		Destination:     dest,
	})
	simulatorManager.AssertExpectations(t)
}

func Test_GivenWaitForBootConfig_WhenWaitForBootFails_ThenReturnsTimeoutError(t *testing.T) {
	// Given
	const (
		dest    = "dest"
		udid    = "test-ID"
		timeout = 1 * time.Second
	)

	var (
		logger           = log.NewLogger()
		simulatorManager = new(mocks.SimulatorManager)
		s                = SimulatorStarter{
			logger:           logger,
			simulatorManager: simulatorManager,
		}
		config = Config{
			Input: Input{
				Destination:        dest,
				WaitForBootTimeout: int(timeout.Seconds()),
			},
			SimulatorID: udid,
		}
	)

	simulatorManager.On("ResetLaunchServices").Once().Return(nil)
	simulatorManager.On("Boot", udid).Once().Return(nil)
	simulatorManager.On("WaitForBootFinished", udid, timeout).Once().Return(errors.New("timeout"))

	// When
	got, err := s.Run(config)

	// Then
	require.Error(t, err)
	require.Equal(t, got, Result{
		SimulatorStatus: "hanged",
		Destination:     dest,
	})
	simulatorManager.AssertExpectations(t)
}
