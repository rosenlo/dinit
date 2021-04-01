package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"dinit/pkg/dinit/config"
	"dinit/pkg/dinit/process"
	"dinit/pkg/util"

	"github.com/rosenlo/toolkits/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version string
	MD5SUM  string
)

const (
	MainProcess string = "main"
)

func parseMainProcess(args []string, conf *config.Config) {
	service, exists := conf.Services[MainProcess]
	if !exists {
		log.Fatal("Main service not found")
	}
	execStart := viper.GetString("command")
	if execStart != "" {
		service.ExecStart = execStart
	}
	if len(args) > 1 {
		service.ExecStart = strings.Join(args, " ")
	}
	if service.ExecStart == "" {
		log.Fatal("Main service is emtpy")
	}
	return
}

func initLog() *os.File {
	LOG_FILE := viper.GetString("log_file")
	file, err := os.OpenFile(LOG_FILE, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}

	formatter := &logrus.JSONFormatter{
		TimestampFormat: TimeFormatFormat,
	}

	log.Init(viper.GetString("log_level"), formatter, file)
	hostname, _ := os.Hostname()
	log.SetField("hostname", hostname)
	for _, extraEnv := range strings.Split(viper.GetString("log_fields"), ",") {
		if extraEnv == "" {
			continue
		}
		val := os.Getenv(extraEnv)
		if val != "" {
			log.SetField(util.Str2Camel([]byte(extraEnv)), val)
		}
	}
	return file
}

func httpServer(sigs chan os.Signal) {
	http.HandleFunc("/init/reload", func(w http.ResponseWriter, r *http.Request) {
		sigs <- syscall.SIGHUP
		w.Write([]byte("success"))
	})
	address := viper.GetString("address")
	s := &http.Server{
		Addr:           address,
		MaxHeaderBytes: 1 << 20,
	}
	log.Infof("listening %s", address)
	s.ListenAndServe()
}

func Run(cmd *cobra.Command, args []string) {
	// 0.Init logging
	file := initLog()
	defer file.Close()

	log.WithField("version", Version).
		WithField("md5sum", MD5SUM).
		Infof("Main process: %v", args)
	// 1.Parse serviers to the struct of the Config
	conf, err := config.ParseConfig(viper.GetString("services_config"))
	if err != nil {
		log.Fatal(err)
	}
	// 2.Supply HTTP Server to receive signal
	sigs := make(chan os.Signal, 1)
	go httpServer(sigs)

RELOAD:
	// 3.Start these services separate
	stopCh := make(chan int, 1)
	ctx, cancel := context.WithCancel(context.Background())
	parseMainProcess(args, conf)
	serviceGroup := process.NewGroup(stopCh, conf)
	go serviceGroup.Run(ctx, cancel)

	// 4.Catch the signal and then shut down the processes
	log.Infof("pid: %d register signal notify", os.Getpid())
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		select {
		case exitCode := <-stopCh:

			if stop(exitCode) {
				os.Exit(exitCode)
			}

		case signal := <-sigs:
			log.Infof("The main process receive %s signal", signal)
			switch signal {

			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				cancel()

				if signal == syscall.SIGHUP {
					goto RELOAD
				}
			}
		}
	}
}

func stop(exitCode int) bool {
	exit := viper.GetBool("exit")
	if exitCode != 0 && !exit {
		log.Info("Main process hanging")
		return false
	}
	log.Infof("Main process quit with code %v, goodbye!", exitCode)
	return true
}
