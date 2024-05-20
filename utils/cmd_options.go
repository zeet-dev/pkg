package utils

import (
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	zoptions "github.com/zeet-dev/pkg/utils/options"
)

type CommandRunner interface {
	RunWithOpt(options ...zoptions.Option[RunOpt]) error
}

type execRunner struct{}

// TODO: refactor RunWithOptV2 to this
func NewExecRunner() execRunner {
	return execRunner{}
}

func (e execRunner) RunWithOpt(options ...zoptions.Option[RunOpt]) error {
	return RunWithOptV2(options...)
}

// RunWithParentEnv runs cmds in dir with additional environment variables supplied by env, in addition to the
// current process' environment (from os.Environ)
func RunWithParentEnvOption(opt *RunOpt) error {
	opt.Env = append(os.Environ(), opt.Env...)
	return nil
}

func RunWithCommandOption(cmds []string, dir string, env []string) zoptions.Option[RunOpt] {
	return func(opt *RunOpt) error {
		opt.Name = cmds[0]
		opt.Args = cmds[1:]
		opt.Dir = dir
		opt.Env = env
		return nil
	}
}

func RunWithStdoutOption(stdout io.Writer) zoptions.Option[RunOpt] {
	return func(opt *RunOpt) error {
		opt.Stdout = stdout
		return nil
	}
}

// TODO: refactor RunWithOpt to this
func RunWithOptV2(options ...zoptions.Option[RunOpt]) error {
	opt, err := zoptions.New(options...)
	if err != nil {
		return err
	}

	cmd := exec.Command(opt.Name, opt.Args...)
	cmd.Dir = opt.Dir
	cmd.Env = opt.Env
	cmd.Stdin = os.Stdin
	if opt.Stdin != nil {
		cmd.Stdin = opt.Stdin
	}
	cmd.Stdout = os.Stdout
	if opt.Stdout != nil {
		cmd.Stdout = opt.Stdout
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return errors.Wrapf(err, "Unable to start %s", opt.Name)
	}
	log.Info().Msgf("Running %s...", opt.Name)
	exitErr := cmd.Wait()
	if exitErr != nil {
		if !opt.ExpectsError {
			return errors.Wrapf(exitErr, "Failed to run %s", opt.Name)
		}
	}
	log.Info().Msgf("Completed %s ✅", opt.Name)
	return exitErr
}
