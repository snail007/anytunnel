package main

import (
	utils "anytunnel/at-common"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const APP_VERSION = "1.0.0"

var (
	cfg = viper.New()
)

func initConfig() (err error) {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	configFile := pflag.String("config", "", "config file path")

	//发布时开启,开发时注释掉
	// cfg.SetDefault("url", "https://atcloud.host900.com:29531/cluster/get")
	cfg.SetDefault("url", "https://127.0.0.1:29531/cluster/get")
	cfg.SetDefault("host", "127.0.0.1")
	cfg.SetDefault("port.conns", 37501)
	cfg.SetDefault("port.control", 37601)
	//结束

	//开发时开启,发布时注释掉
	// pflag.String("url", "", "dispatch url of cluster to connecting")
	// cfg.BindPFlag("url", pflag.Lookup("url"))

	// pflag.String("host", "127.0.0.1", "host of cluster to connecting")
	// cfg.BindPFlag("host", pflag.Lookup("host"))

	// pflag.Int("conns-port", 37501, "port of cluster to connecting")
	// cfg.BindPFlag("port.conns", pflag.Lookup("conns-port"))

	// pflag.Int("control-port", 37601, "port of cluster to control")
	// cfg.BindPFlag("port.control", pflag.Lookup("control-port"))
	//结束

	pflag.String("token", "", "token of client connect to cluster")
	cfg.BindPFlag("token", pflag.Lookup("token"))

	pflag.Int("udp-timeout", 2, "seconds of waiting for remote udp server reponse")
	cfg.BindPFlag("udp.timeout", pflag.Lookup("udp-timeout"))

	pflag.Bool("log-open", false, "if true: store log files, false: no log files")
	pflag.String("level", "debug", "console log level,should be one of debug,info,warn,error")
	pflag.String("log-dir", "log", "the directory which store log files")
	pflag.Bool("log-post", false, "log post data on or off")
	pflag.Int64("log-max-size", 102400000, "log file max size(bytes) for rotate")
	pflag.Int("log-max-count", 3, "log file max count for rotate to remain")
	pflag.StringSlice("log-level", []string{"info", "error", "debug"}, "log to file level,multiple splitted by comma(,)")
	pflag.Parse()

	cfg.BindPFlag("log.dir", pflag.Lookup("log-dir"))
	cfg.BindPFlag("log.level", pflag.Lookup("log-level"))
	cfg.BindPFlag("log.open", pflag.Lookup("log-open"))
	cfg.BindPFlag("log.console-level", pflag.Lookup("level"))
	cfg.BindPFlag("log.fileMaxSize", pflag.Lookup("log-max-size"))
	cfg.BindPFlag("log.maxCount", pflag.Lookup("log-max-count"))
	addr, err := utils.GetClusterHost(cfg.GetString("url"), cfg.GetString("token"), "client")
	if err == nil {
		cfg.Set("host", addr)
	} else if cfg.GetString("url") != "" {
		fmt.Printf("%s , please try again\n", err)
		os.Exit(0)
	}

	if *configFile != "" {
		cfg.SetConfigFile(*configFile)
	} else {
		cfg.SetConfigName("client")
		cfg.AddConfigPath("/etc/anytunnel/")
		cfg.AddConfigPath("$HOME/.anytunnel")
		cfg.AddConfigPath(".anytunnel")
		cfg.AddConfigPath(".")
	}
	err = cfg.ReadInConfig()
	file := cfg.ConfigFileUsed()
	if err != nil && !strings.Contains(err.Error(), "Not") {
		fmt.Printf("%s", err)
	} else if file != "" {
		fmt.Printf("use config file : %s\n", file)
	}
	err = nil
	return
}

func poster() {
	fmt.Printf(`
	╔═╗┌┐┌┬ ┬╔╦╗┬ ┬┌┐┌┌┐┌┌─┐┬   ╔═╗┬  ┬┌─┐┌┐┌┌┬┐
	╠═╣│││└┬┘ ║ │ │││││││├┤ │───║  │  │├┤ │││ │ 
	╩ ╩┘└┘ ┴  ╩ └─┘┘└┘┘└┘└─┘┴─┘ ╚═╝┴─┘┴└─┘┘└┘ ┴  v%s`+"\n\n", APP_VERSION)
}
