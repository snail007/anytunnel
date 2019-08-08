package main

import (
	utils "anytunnel/at-common"
	"flag"
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const APP_VERSION = "1.0.0"

var (
	cfg            = viper.New()
	tokenServerMap = map[string]bool{}
	tokenClientMap = map[string]bool{}
)

func initConfig() (err error) {

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	configFile := pflag.String("config", "", "config file path")

	pflag.String("data-dir", "./data", "data dir to store opened tunnels , and recovered at next start")

	pflag.Bool("url-is-internal", true, "http api port")
	pflag.String("auth-url", "", "check token when c/s login in,this url return http status 204 auth success,other fail.when access this url token will attached ,such as url?token=xxx")
	pflag.String("status-url", "", "c/s online or offline , cluster will access this url")
	pflag.String("auth-file", "", "server or client token auth,every line one token")
	pflag.StringSlice("bind-ip-data", []string{"0.0.0.0"}, "ip to bind when data listener listening")
	pflag.StringSlice("bind-ip-control", []string{"0.0.0.0"}, "ip to bind when control listener listening")
	pflag.StringSlice("bind-ip-api", []string{"0.0.0.0"}, "ip to bind when api listening")
	pflag.Int("api-port", 37080, "http api port")
	pflag.Int("conns-port", 37501, "port to listen and wait for server/client to connect")
	pflag.Int("control-port", 37601, "port to listen and control server/client")
	pflag.String("traffic-url", "", "post traffic statist to this url , if url not empty.")
	pflag.Int("traffic-interval", 30, "post traffic to traffic-url , every traffic-interval seconds")
	pflag.Int("url-fail-retry", 1, "retry count when access all url fail")
	pflag.Int("url-fail-wait", 3, "wait seconds when access all url fail")
	pflag.Int("url-success-code", 204, "this decide access url if success or not,check compare this code and response http code after access all url")
	pflag.Bool("log-open", true, "if true: store log files, false: no log files")
	pflag.String("level", "debug", "console log level,should be one of debug,info,warn,error")
	pflag.String("log-dir", "log", "the directory which store log files")
	pflag.Bool("log-access", true, "access log on or off")
	pflag.Bool("log-post", false, "log post data on or off")
	pflag.Int64("log-max-size", 102400000, "log file max size(bytes) for rotate")
	pflag.Int("log-max-count", 3, "log file max count for rotate to remain")
	pflag.StringSlice("log-level", []string{"info", "error", "debug"}, "log to file level,multiple splitted by comma(,)")
	pflag.Parse()

	cfg.BindPFlag("url.is-internal", pflag.Lookup("url-is-internal"))
	cfg.BindPFlag("url.fail-retry", pflag.Lookup("url-fail-retry"))
	cfg.BindPFlag("url.fail-wait", pflag.Lookup("url-fail-wait"))
	cfg.BindPFlag("url.success-code", pflag.Lookup("url-success-code"))
	cfg.BindPFlag("url.auth", pflag.Lookup("auth-url"))
	cfg.BindPFlag("url.status", pflag.Lookup("status-url"))
	cfg.BindPFlag("url.traffic", pflag.Lookup("traffic-url"))
	cfg.BindPFlag("url.traffic-interval", pflag.Lookup("traffic-interval"))
	cfg.BindPFlag("auth.file", pflag.Lookup("auth-file"))
	cfg.BindPFlag("port.ip-data", pflag.Lookup("bind-ip-data"))
	cfg.BindPFlag("port.ip-control", pflag.Lookup("bind-ip-control"))
	cfg.BindPFlag("port.ip-api", pflag.Lookup("bind-ip-api"))
	cfg.BindPFlag("port.api", pflag.Lookup("api-port"))
	cfg.BindPFlag("port.conns", pflag.Lookup("conns-port"))
	cfg.BindPFlag("port.control", pflag.Lookup("control-port"))
	cfg.BindPFlag("data.dir", pflag.Lookup("data-dir"))

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
		cfg.SetConfigName("cluster")
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
	loadAuthFile(cfg.GetString("auth.file"))
	return
}

func poster() {
	fmt.Printf(`
	╔═╗┌┐┌┬ ┬╔╦╗┬ ┬┌┐┌┌┐┌┌─┐┬   ╔═╗┬  ┬ ┬┌─┐┌┬┐┌─┐┬─┐
	╠═╣│││└┬┘ ║ │ │││││││├┤ │───║  │  │ │└─┐ │ ├┤ ├┬┘
	╩ ╩┘└┘ ┴  ╩ └─┘┘└┘┘└┘└─┘┴─┘ ╚═╝┴─┘└─┘└─┘ ┴ └─┘┴└─ v%s`+"\n\n", APP_VERSION)
}
func loadAuthFile(filePath string) {
	if filePath == "" {
		return
	}
	content, err := utils.FileGetContents(filePath)
	if err != nil {
		fmt.Printf("access %s fail , err:%s\n", filePath, err)
	}

	lines := strings.Split(strings.Trim(content, "\r\n"), "\n")
	i := 0
	for _, line := range lines {
		if len(line) > 0 && line[0:1] == "#" {
			break
		}
		if len(strings.TrimSpace(line)) > 7 && strings.Trim(line, "\r\n") != "" {
			if strings.Index(line, "server:") == 0 {
				tokenServerMap[line[7:]] = true
				i++
			}
			if strings.Index(line, "client:") == 0 {
				tokenClientMap[line[7:]] = true
				i++
			}
		}
	}
	fmt.Printf("auth file loaded,tokens:%d,%s\n", i, filePath)
	return
}
