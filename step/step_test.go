package step

import (
	"errors"
	"testing"
	"time"

	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/bitrise-step-xcode-start-simulator/step/mocks"
	"github.com/stretchr/testify/require"
)

func TestSimulatorStarter_Run_WhenBootOnly_DoesBoot(t *testing.T) {
	const (
		dest = "dest"
		udid = "test-ID"
	)

	logger := log.NewLogger()
	simulatorManager := new(mocks.SimulatorManager)
	simulatorManager.On("ResetLaunchServices").Once().Return(nil)
	simulatorManager.On("Boot", udid).Once().Return(nil)

	s := SimulatorStarter{
		logger:           logger,
		simulatorManager: simulatorManager,
	}
	config := Config{
		Input: Input{
			Destination: dest,
			WaitForBoot: false,
		},
		SimulatorID: udid,
	}

	got, err := s.Run(config)

	require.NoError(t, err)
	require.Equal(t, got, Result{
		IsSimulatorTimeout: false,
		Destination:        dest,
	})
	simulatorManager.AssertExpectations(t)
}

func TestSimulatorStarter_Run_WhenBootFails_ReturnsError(t *testing.T) {
	const (
		dest = "dest"
		udid = "test-ID"
	)

	logger := log.NewLogger()
	simulatorManager := new(mocks.SimulatorManager)
	simulatorManager.On("ResetLaunchServices").Once().Return(nil)
	simulatorManager.On("Boot", udid).Once().Return(errors.New("boot error"))

	s := SimulatorStarter{
		logger:           logger,
		simulatorManager: simulatorManager,
	}
	config := Config{
		Input: Input{
			Destination: dest,
			WaitForBoot: false,
		},
		SimulatorID: udid,
	}

	got, err := s.Run(config)

	require.Error(t, err)
	require.Equal(t, got, Result{
		IsSimulatorTimeout: false,
		Destination:        dest,
	})
	simulatorManager.AssertExpectations(t)
}

func TestSimulatorStarter_Run_WhenWaitForBootFails_ReturnsTimeoutError(t *testing.T) {
	const (
		dest    = "dest"
		udid    = "test-ID"
		timeout = 1 * time.Second
	)

	logger := log.NewLogger()
	simulatorManager := new(mocks.SimulatorManager)
	simulatorManager.On("ResetLaunchServices").Once().Return(nil)
	simulatorManager.On("Boot", udid).Once().Return(nil)
	simulatorManager.On("WaitForBootFinished", udid, timeout).Once().Return(errors.New("timeout"))

	s := SimulatorStarter{
		logger:           logger,
		simulatorManager: simulatorManager,
	}
	config := Config{
		Input: Input{
			Destination:        dest,
			WaitForBoot:        true,
			WaitForBootTimeout: int(timeout.Seconds()),
		},
		SimulatorID: udid,
	}

	got, err := s.Run(config)

	require.Error(t, err)
	require.Equal(t, got, Result{
		IsSimulatorTimeout: true,
		Destination:        dest,
	})
	simulatorManager.AssertExpectations(t)
}
