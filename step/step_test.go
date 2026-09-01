package step

import (
	"errors"
	"testing"
	"time"

	"github.com/bitrise-io/go-utils/v2/command"
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
		UDID:            udid,
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
		UDID:            udid,
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
		UDID:            udid,
	})
}

func Test_GivenResetConfigAndShutdownDevice_WhenRun_ThenSkipsShutdown(t *testing.T) {
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
			Status:   "Shutdown",
		}
		config = Config{
			Input: Input{
				WaitForBootTimeout: 0,
				ShouldReset:        true,
			},
			Simulator: simulator,
		}
	)

	simulatorManager.On("ResetLaunchServices").Once().Return(nil)
	simulatorManager.On("Erase", udid).Once().Return(nil)
	simulatorManager.On("Boot", simulator).Once().Return(nil)

	// When
	got, err := s.Run(config)

	// Then
	require.NoError(t, err)
	require.Equal(t, got.SimulatorStatus, "booted")
	simulatorManager.AssertNotCalled(t, "Shutdown", udid)
}

func Test_GivenResetConfigAndBootedDevice_WhenRun_ThenShutsDownAndErases(t *testing.T) {
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
			Status:   "Booted",
		}
		config = Config{
			Input: Input{
				WaitForBootTimeout: 0,
				ShouldReset:        true,
			},
			Simulator: simulator,
		}
	)

	simulatorManager.On("ResetLaunchServices").Once().Return(nil)
	simulatorManager.On("Shutdown", udid).Once().Return(nil)
	simulatorManager.On("Erase", udid).Once().Return(nil)
	simulatorManager.On("Boot", simulator).Once().Return(nil)

	// When
	got, err := s.Run(config)

	// Then
	require.NoError(t, err)
	require.Equal(t, got.SimulatorStatus, "booted")
}

func Test_GivenDarkAppearanceConfig_WhenRun_ThenSetsAppearance(t *testing.T) {
	// Given
	const udid = "test-ID"

	var (
		logger           = log.NewLogger()
		simulatorManager = mocks.NewSimulatorManager(t)
		commandFactory   = mocks.NewCommandFactory(t)
		cmd              = mocks.NewCommand(t)
		s                = SimulatorStarter{
			logger:           logger,
			simulatorManager: simulatorManager,
			commandFactory:   commandFactory,
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
				Appearance:         "dark",
			},
			Simulator: simulator,
		}
	)

	simulatorManager.On("ResetLaunchServices").Once().Return(nil)
	simulatorManager.On("Boot", simulator).Once().Return(nil)
	commandFactory.On("Create", "xcrun", []string{"simctl", "ui", udid, "appearance", "dark"}, (*command.Opts)(nil)).Once().Return(cmd)
	cmd.On("PrintableCommandArgs").Once().Return("xcrun simctl ui test-ID appearance dark")
	cmd.On("RunAndReturnTrimmedCombinedOutput").Once().Return("", nil)

	// When
	got, err := s.Run(config)

	// Then
	require.NoError(t, err)
	require.Equal(t, got, Result{
		SimulatorStatus: "booted",
		Destination:     "platform=iOS Simulator,name=Bitrise iOS default,OS=11",
		UDID:            udid,
	})
}

func Test_GivenDarkAppearanceConfig_WhenSetAppearanceFails_ThenReturnsFailedStatus(t *testing.T) {
	// Given
	const udid = "test-ID"

	var (
		logger           = log.NewLogger()
		simulatorManager = mocks.NewSimulatorManager(t)
		commandFactory   = mocks.NewCommandFactory(t)
		cmd              = mocks.NewCommand(t)
		s                = SimulatorStarter{
			logger:           logger,
			simulatorManager: simulatorManager,
			commandFactory:   commandFactory,
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
				Appearance:         "dark",
			},
			Simulator: simulator,
		}
	)

	simulatorManager.On("ResetLaunchServices").Once().Return(nil)
	simulatorManager.On("Boot", simulator).Once().Return(nil)
	commandFactory.On("Create", "xcrun", []string{"simctl", "ui", udid, "appearance", "dark"}, (*command.Opts)(nil)).Once().Return(cmd)
	cmd.On("PrintableCommandArgs").Once().Return("xcrun simctl ui test-ID appearance dark")
	cmd.On("RunAndReturnTrimmedCombinedOutput").Once().Return("Invalid device state", errors.New("exit status 1"))

	// When
	got, err := s.Run(config)

	// Then
	require.Error(t, err)
	require.Equal(t, got, Result{
		SimulatorStatus: "failed",
		Destination:     "platform=iOS Simulator,name=Bitrise iOS default,OS=11",
		UDID:            udid,
	})
}
