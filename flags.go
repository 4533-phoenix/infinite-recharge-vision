package main

import "github.com/spf13/pflag"

// configPath is the configuration file path flag value.
var configPath string

// logLevel is the log-level flag value.
var logLevel string

func init() {
	pflag.StringVarP(&configPath, "config", "", "", "path to config file")
	pflag.StringVarP(&logLevel, "log-level", "", "info", "logging level")
}
