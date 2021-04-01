package process

import (
	"context"
	"sync"
	"time"

	"dinit/pkg/dinit/config"

	"github.com/rosenlo/toolkits/log"
	"github.com/spf13/viper"
)

type Group struct {
	Processes map[string]Manager
	conf      *config.Config
	sync.RWMutex
	stopCh chan int
}

func NewGroup(stopCh chan int, conf *config.Config) *Group {
	group := &Group{
		Processes: make(map[string]Manager),
		conf:      conf,
		stopCh:    stopCh,
	}
	for name, service := range conf.Services {
		group.AddProcess(name, service)
	}
	return group
}

func (c *Group) Get(name string) (Manager, bool) {
	c.RLock()
	defer c.RUnlock()
	manager, exists := c.Processes[name]
	return manager, exists
}

func (c *Group) GetAll() map[string]Manager {
	c.RLock()
	defer c.RUnlock()
	return c.Processes
}

// Watch watch child processes, cancel context as child processes exits
func (c *Group) Watch(ctx context.Context, cancel context.CancelFunc) {

	select {
	case <-ctx.Done():
		log.Infof("Main %v, cancel other processes", ctx.Err())
		c.Stop(cancel)
	}
}

// Wait wait for dependencies to start first
func (c *Group) Wait(ctx context.Context, cancel context.CancelFunc, name string, depends []string) {
	log := log.WithField("name", name)
	for _, depend := range depends {
		dependSvc, exists := c.Processes[depend]
		if !exists {
			log.WithField("depend", depend).Fatal("Dependent service do not exists.")
		}
		for !dependSvc.Ready() {
			log.WithField("depend", depend).Info("Waiting dependencies")
			if dependSvc.Terminated() {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
	go c.Processes[name].Start(ctx, cancel)
}

func (c *Group) Run(ctx context.Context, cancel context.CancelFunc) {

	groupCtx, groupCancel := context.WithCancel(context.Background())
	go c.Watch(ctx, groupCancel)

	for name, svc := range c.conf.Services {
		childCtx, childCancel := context.WithCancel(groupCtx)
		if len(svc.DependsOn) != 0 {
			go c.Wait(childCtx, childCancel, name, svc.DependsOn)
		} else {
			go c.Processes[name].Start(childCtx, childCancel)
		}

		go func(n string) {
			select {
			case <-childCtx.Done():
				log.WithField("name", n).Debugf("Child %v, notify the main context", childCtx.Err())
				cancel()
			}
		}(name)
	}
}

func (c *Group) Stop(cancel context.CancelFunc) (exitCode int) {
	gracefulTimeout := viper.GetInt("graceful_timeout")

	if cancel != nil {
		log.Info("Child cancel")
		cancel()
	}

	for name, svc := range c.GetAll() {
		log := log.WithField("name", name)

		for gracefulTimeout > 0 {
			gracefulTimeout -= 1

			if svc.Terminated() {
				log.WithField("status", svc.Status()).Debug("The process has stopped")
				break
			} else {
				log.Info("Waiting stop")
				time.Sleep(1 * time.Second)
			}
		}

		if svc.ExitCode() != 0 {
			exitCode = svc.ExitCode()
		}
		c.DelProcess(name)
	}

	c.stopCh <- exitCode

	return
}

func (c *Group) AddProcess(name string, svc *config.Service) {
	if svc.ExecStart != "" {
		gracefulTimeout := viper.GetInt("graceful_timeout")
		c.Processes[name] = NewManager(name, gracefulTimeout, svc)
	} else {
		log.WithField("name", name).Warn("Failed to create a command manager")
	}
}

func (c *Group) DelProcess(name string) {
	c.Lock()
	defer c.Unlock()
	delete(c.Processes, name)
}
