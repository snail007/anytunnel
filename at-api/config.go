package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const APP_VERSION = "1.0.0"

var (
	cfg = viper.New()
)

//init config
func initConfig() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	configFile := pflag.String("config", "", "config file path")
	pflag.StringSlice("bind-ip-intra", []string{"0.0.0.0"}, "ip to binding for intra api")
	pflag.Int("port-intra", 37081, "intra api port")

	pflag.StringSlice("bind-ip-extra", []string{"0.0.0.0"}, "ip to binding for extra api")
	pflag.Int("port-extra", 29531, "extra api port")

	pflag.Bool("log-open", true, "if true: store log files, false: no log files")
	pflag.String("level", "debug", "console log level,should be one of debug,info,warn,error")
	pflag.String("log-dir", "log", "the directory which store log files")
	pflag.Bool("log-access", true, "access log on or off")
	pflag.Bool("log-post", false, "log post data on or off")
	pflag.Int64("log-max-size", 102400000, "log file max size(bytes) for rotate")
	pflag.Int("log-max-count", 3, "log file max count for rotate to remain")
	pflag.StringSlice("log-level", []string{"info", "error", "debug"}, "log to file level,multiple splitted by comma(,)")
	pflag.Parse()

	cfg.BindPFlag("port.ip-intra", pflag.Lookup("bind-ip-intra"))
	cfg.BindPFlag("port.ip-extra", pflag.Lookup("bind-ip-extra"))
	cfg.BindPFlag("port.intra", pflag.Lookup("port-intra"))
	cfg.BindPFlag("port.extra", pflag.Lookup("port-extra"))

	cfg.BindPFlag("log.open", pflag.Lookup("log-open"))
	cfg.BindPFlag("log.dir", pflag.Lookup("log-dir"))
	cfg.BindPFlag("log.level", pflag.Lookup("log-level"))
	cfg.BindPFlag("log.access", pflag.Lookup("log-access"))
	cfg.BindPFlag("log.post", pflag.Lookup("log-post"))
	cfg.BindPFlag("log.console-level", pflag.Lookup("level"))
	cfg.BindPFlag("log.fileMaxSize", pflag.Lookup("log-max-size"))
	cfg.BindPFlag("log.maxCount", pflag.Lookup("log-max-count"))
	fmt.Printf("%s", *configFile)
	if *configFile != "" {
		cfg.SetConfigFile(*configFile)
	} else {
		cfg.SetConfigName("api")
		cfg.AddConfigPath("/etc/anytunnel/")
		cfg.AddConfigPath("$HOME/.anytunnel")
		cfg.AddConfigPath(".anytunnel")
		cfg.AddConfigPath(".")
	}
	err := cfg.ReadInConfig()
	file := cfg.ConfigFileUsed()
	if err != nil && !strings.Contains(err.Error(), "Not") {
		fmt.Printf("%s\n", err)
	} else if file != "" {
		fmt.Printf("use config file : %s\n", file)
	}
	err = nil
	return
}

// poster version
func poster() {
	fmt.Printf(`
	╔═╗┌┐┌┬ ┬┌┬┐┬ ┬┌┐┌┌┐┌┌─┐┬   ╔═╗┌─┐┬
	╠═╣│││└┬┘ │ │ │││││││├┤ │───╠═╣├─┘│
	╩ ╩┘└┘ ┴  ┴ └─┘┘└┘┘└┘└─┘┴─┘ ╩ ╩┴  ┴ v%s`+"\n\n", APP_VERSION)
}
