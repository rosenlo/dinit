package util

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rosenlo/toolkits/log"
)

type Config struct {
	AppID    string `mapstructure:"app_id"`
	Address  string `mapstructure:"address"`
	LogLevel string `mapstructure:"log_level"`
}

type Server struct {
	Config *Config
}

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v *Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func NewServer(c *Config) *Server {
	log.Init(c.LogLevel, nil, nil)

	server := &Server{
		Config: c,
	}

	return server
}

func (s *Server) ListenSignal(exitfunc ...func()) {
	StartSignal(os.Getpid(), exitfunc...)
}

func StartSignal(pid int, exitfunc ...func()) {
	sigs := make(chan os.Signal, 1)
	log.Infof("pid: %d register signal notify", pid)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		s := <-sigs
		log.Infof("recv %s", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			if len(exitfunc) != 0 {
				for _, f := range exitfunc {
					f()
				}
			}
			log.Info("graceful shut down")
			log.Infof("main pid: %d exit", pid)
			os.Exit(0)
		}
	}
}
