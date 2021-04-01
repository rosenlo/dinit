package app

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/rosenlo/toolkits/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

const (
	TimeFormatFormat string = "2006-01-02 15:04:05.000 Z0700"
)

// dinit represents the base command when called without any subcommands
var dinit = &cobra.Command{
	Use:                "dinit",
	Short:              "Container daemon process",
	Long:               `Container daemon process`,
	DisableFlagParsing: true,
	Run:                Run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return dinit.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		viper.SetConfigName("dinit")
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/dinit")
	}
	viper.SetEnvPrefix("DINIT")
	// basic
	viper.SetDefault("address", ":8888")
	viper.SetDefault("command", "")
	viper.SetDefault("log_level", "debug")
	viper.SetDefault("log_file", "/data/logs/dinit/dinit.log")
	viper.SetDefault("log_fields", "APP_NAME,POD_IP,NODE_NAME")
	viper.SetDefault("services_config", "/etc/dinit/dinit.yaml")
	viper.SetDefault("http_retry", 3)
	viper.SetDefault("http_retry_time", 3)

	// process
	viper.SetDefault("graceful_timeout", 30)
	viper.SetDefault("exit", true)
	viper.SetDefault("executor", "/bin/sh")
	viper.SetDefault("executor_arg", "-ec")

	// flowcontrol
	viper.SetDefault("openresty_addr", "")
	viper.SetDefault("service_prefix_key", "nginx/upstreams/")
	viper.SetDefault("health_check_interval", 10)
	viper.SetDefault("consul_username", "")
	viper.SetDefault("consul_password", "")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// test config
	if viper.GetString("env") == "test" {
	}
}
