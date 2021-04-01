package process

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"dinit/pkg/dinit/config"

	"github.com/rosenlo/toolkits/log"
	"github.com/spf13/viper"
)

type ManagerState int

const (
	StatePending ManagerState = iota
	StateRunning
	StateStopping
	StateSucceeded
	StateFailed
	StateTerminated
)

var StateName = map[ManagerState]string{
	StatePending:   "pending",
	StateRunning:   "running",
	StateStopping:  "stopping",
	StateSucceeded: "succeeded",
	StateFailed:    "failed",
}

type Manager interface {
	Status() string
	Ready() bool
	Terminated() bool
	ExitCode() int
	PreStart() error
	PostStart() error
	Start(ctx context.Context, cancel context.CancelFunc)
	PreStop() error
	Stop()
	PostStop() error
}

type Managers struct {
	service         *config.Service
	path            string
	arg             string
	cmd             *exec.Cmd
	name            string
	fullName        string
	gracefulTimeout int
	status          ManagerState
	exitCode        int
	stderr          *bytes.Buffer
	finished        bool
}

func NewManager(name string, gracefulTimeout int, svc *config.Service) *Managers {
	var stderr bytes.Buffer
	executor := viper.GetString("executor")
	executorArg := viper.GetString("executor_arg")

	cmd := NewCommand(executor, executorArg, svc.ExecStart)

	if svc != nil && len(svc.Environment) > 0 {
		env := make([]string, 0)
		for k, v := range svc.Environment {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		cmd.Env = append(os.Environ(), env...)
	}
	cmd.Stderr = &stderr

	return &Managers{
		path:            executor,
		arg:             executorArg,
		service:         svc,
		name:            name,
		fullName:        svc.ExecStart,
		gracefulTimeout: gracefulTimeout,
		cmd:             cmd,
		stderr:          &stderr,
	}
}

func (c *Managers) Ready() bool {
	return c.status == StateRunning
}

func (c *Managers) Status() string {
	return StateName[c.status]
}

func (c *Managers) Terminated() bool {
	return c.cmd.ProcessState != nil && c.finished
}

func (c *Managers) ExitCode() int {
	return c.exitCode
}

func (c *Managers) PreStart() error {
	if c.service.Lifecycle == nil || c.service.Lifecycle.PreStart == nil {
		return nil
	}
	preStart := c.service.Lifecycle.PreStart
	if preStart.Exec.Command != "" {
		log.WithField("name", c.name).
			WithField("preStart", preStart.Exec.Command).
			Debug("pre_start tigger new command")

		var stderr bytes.Buffer
		cmd := NewCommand(c.path, c.arg, preStart.Exec.Command)
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return errors.New(stderr.String())
		}
	}
	return nil
}

func (c *Managers) Start(ctx context.Context, cancel context.CancelFunc) {
	log := log.WithField("name", c.name)

	if err := c.PreStart(); err != nil {
		log.Errorf("pre_start exec failed, due to %v", err)
	}
	if err := c.cmd.Start(); err != nil {
		log.Errorf("start with err: %v", err)
		c.Stop()
		return
	}

	if err := c.PostStart(); err != nil {
		log.Errorf("post_start exec failed, due to %v", err)
	}

	log.WithField("pid", c.cmd.Process.Pid).Debug("Process started")
	c.status = StateRunning
	go c.wait(ctx, cancel)
}

func (c *Managers) wait(ctx context.Context, cancel context.CancelFunc) {
	waitDone := make(chan error, 1)
	go func() {
		waitDone <- c.cmd.Wait()
	}()

	log := log.WithField("name", c.name)

	select {
	case err := <-waitDone:
		c.status = StateStopping
		if err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				status := exiterr.Sys().(syscall.WaitStatus)
				if !status.Signaled() {
					log.WithField("stderr", c.stderr.String()).Error(exiterr.Error())
					c.exitCode = exiterr.ExitCode()
					c.status = StateFailed
				}
			} else {
				log.Errorf("Process exit with err: %s", err)
			}
		}
		log.Infof("wait with %v", err)
		cancel()
	case <-ctx.Done():
		log.Debug("childCtx Done")
		c.Stop()
	}
}

func (c *Managers) PostStart() error {
	if c.service.Lifecycle == nil || c.service.Lifecycle.PostStart == nil {
		return nil
	}
	postStart := c.service.Lifecycle.PostStart
	if postStart.Exec.Command != "" {
		log.WithField("name", c.name).
			WithField("cmd", postStart.Exec.Command).
			Debug("post_start tigger new command")

		var stderr bytes.Buffer
		cmd := NewCommand(c.path, c.arg, postStart.Exec.Command)
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return errors.New(stderr.String())
		}
	}
	return nil
}

func (c *Managers) PreStop() error {
	if c.service.Lifecycle == nil || c.service.Lifecycle.PreStop == nil {
		return nil
	}
	preStop := c.service.Lifecycle.PreStop
	if preStop.Exec.Command != "" {
		log.WithField("name", c.name).
			WithField("cmd", preStop.Exec.Command).
			Debug("pre_stop tigger new command")

		var stderr bytes.Buffer
		cmd := NewCommand(c.path, c.arg, preStop.Exec.Command)
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return errors.New(stderr.String())
		}
	}
	return nil
}

func (c *Managers) Stop() {
	log := log.WithField("name", c.name).WithField("pid", c.cmd.Process.Pid)

	log.Debug("The child process was receive shutdown signal")
	if err := c.PreStop(); err != nil {
		log.Errorf("pre_stop exec failed, due to %v", err)
	}

	if c.status != StateStopping {
		if err := syscall.Kill(-c.cmd.Process.Pid, syscall.SIGTERM); err != nil {
			c.status = StateFailed
			log.Warnf("Sigterm with err: %v", err)

			time.Sleep(time.Duration(c.gracefulTimeout) * time.Second)

			log.Warn("Sigkill")
			if err := syscall.Kill(-c.cmd.Process.Pid, syscall.SIGKILL); err != nil {
				log.Warnf("Kill with err: %v", err)
			}
			log.Info("Process killed as graceful shutdown timeout reached")
		} else {
			c.status = StateSucceeded
		}
	}

	if err := c.PostStop(); err != nil {
		log.Errorf("postStop exec failed, due to %v", err)
	}

	c.finished = true
	log.WithField("status", c.Status()).Debug("Process finished")
	return
}

func (c *Managers) PostStop() error {
	if c.service.Lifecycle == nil || c.service.Lifecycle.PostStop == nil {
		return nil
	}
	postStop := c.service.Lifecycle.PostStop
	if postStop.Exec.Command != "" {
		log.WithField("name", c.name).
			WithField("cmd", postStop.Exec.Command).
			Info("post_stop tigger new command")

		var stderr bytes.Buffer
		cmd := NewCommand(c.path, c.arg, postStop.Exec.Command)
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return errors.New(stderr.String())
		}
	}
	return nil
}
