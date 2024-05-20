package utils

import (
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type RunOpt struct {
	Name string
	Args []string
	Dir  string
	// nil runs with current process environment
	Env    []string
	Stdin  io.Reader
	Stdout io.Writer

	// If true, RunWithOpt will return error as-is without wrapping it
	ExpectsError bool
}

// RunWithParentEnv runs cmds in dir with additional environment variables supplied by env, in addition to the
// current process' environment (from os.Environ)
func RunWithParentEnv(cmds []string, dir string, env []string, options ...func(opt *RunOpt)) error {
	name := cmds[0]
	args := cmds[1:]
	combinedEnv := os.Environ()
	combinedEnv = append(combinedEnv, env...)

	opt := RunOpt{
		Name: name,
		Args: args,
		Dir:  dir,
		Env:  combinedEnv,
	}
	for _, option := range options {
		option(&opt)
	}
	return RunWithOpt(opt)
}

func Run(cmds []string, dir string, env []string) error {
	name := cmds[0]
	args := cmds[1:]

	return RunWithOpt(RunOpt{
		Name: name,
		Args: args,
		Dir:  dir,
		Env:  env,
	})
}

func RunWithOpt(opt RunOpt) error {
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
	if err := cmd.Wait(); err != nil {
		return errors.Wrapf(err, "Failed to run %s", opt.Name)
	}
	log.Info().Msgf("Completed %s ✅", opt.Name)

	return nil
}
