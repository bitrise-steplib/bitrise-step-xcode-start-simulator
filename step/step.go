package step

import (
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/steps-xcode-test/simulator"
)

type Input struct {
	Destination string `env:"destination,required"`
}

type Config struct {
	SimulatorID       string
	IsSimulatorBooted bool
}

type Result struct {
	IsSimulatorTimeout bool
}

type SimulatorStarter struct {
	logger         log.Logger
	inputParser    stepconf.InputParser
	commandFactory command.Factory
	manager        simulator.Manager
}

func NewStep(logger log.Logger, inputParser stepconf.InputParser, commandFactory command.Factory, simualatorManager simulator.Manager) SimulatorStarter {
	return SimulatorStarter{
		logger:         logger,
		inputParser:    inputParser,
		commandFactory: commandFactory,
		manager:        simualatorManager,
	}
}

func (s SimulatorStarter) ProcessConfig() (Config, error) {
	var input Input
	if err := s.inputParser.Parse(&input); err != nil {
		return Config{}, err
	}

	stepconf.Print(input)
	s.logger.Println()

	return Config{
		SimulatorID: "",
	}, nil
}

func (s SimulatorStarter) InstallDependencies() error {
	return nil
}

func (s SimulatorStarter) Run(config Config) (Result, error) {
	return Result{}, nil
}

func (s SimulatorStarter) ExportOtputs(result Result) error {
	return nil
}
