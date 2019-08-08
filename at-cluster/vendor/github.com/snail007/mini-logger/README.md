# mini-logger
mini but flexible and powerful logger for go
# Notice
1.Do not call runtime.Goexit() in main , it will be blocking logger.Flush().   
# Features
<ul>
<li>Console</li>
<li>File,supoort rotate base size</li>
</ul>
# Usage Example
<pre>
package main

import (
	"github.com/snail007/mini-logger"
	"github.com/snail007/mini-logger/writers/console"
	"github.com/snail007/mini-logger/writers/files"
)

var log logger.MiniLogger
var accessLog logger.MiniLogger

//initLog
func initLog() {
	var level uint8
	switch cfg.GetString("log.console-level") {
	case "debug":
		level = logger.AllLevels
	case "info":
		level = logger.InfoLevel | logger.WarnLevel | logger.ErrorLevel | logger.FatalLevel
	case "warn":
		level = logger.WarnLevel | logger.ErrorLevel | logger.FatalLevel
	case "error":
		level = logger.ErrorLevel | logger.FatalLevel
	case "fatal":
		level = logger.FatalLevel
	default:
		level = 0
	}
	log = logger.New(false, nil)
	log.AddWriter(console.NewDefault(), level)
	cfgF := files.GetDefaultFileConfig()
	cfgF.LogPath = cfg.GetString("log.dir")
	cfgF.MaxBytes = cfg.GetInt64("log.FileMaxSize")
	cfgF.MaxCount = cfg.GetInt("log.MaxCount")
	cfgLevels := cfg.GetStringSlice("log.level")
	if ok, _ := inArray("debug", cfgLevels); ok {
		cfgF.FileNameSet["debug"] = logger.AllLevels
	}
	if ok, _ := inArray("info", cfgLevels); ok {
		cfgF.FileNameSet["info"] = logger.InfoLevel
	}
	if ok, _ := inArray("error", cfgLevels); ok {
		cfgF.FileNameSet["error"] = logger.WarnLevel | logger.ErrorLevel | logger.FatalLevel
	}
	log.AddWriter(files.New(cfgF), logger.AllLevels)

	accessLog = logger.New(false, nil)
	//accessLog.AddWriter(console.NewDefault(), logger.AllLevels)
	if cfg.GetBool("log.access") {
		accessCfg := files.GetDefaultFileConfig()
		accessCfg.LogPath = cfg.GetString("log.dir")
		accessCfg.MaxBytes = cfg.GetInt64("log.FileMaxSize")
		accessCfg.MaxCount = cfg.GetInt("log.MaxCount")
		accessCfg.FileNameSet = map[string]uint8{"access": logger.InfoLevel}
		accessLog.AddWriter(files.New(accessCfg), logger.InfoLevel)
	}

    log.With(logger.Fields{"func": "getMqConnection", "call": "pools.Get"}).Errorf("fail,%s", err)
}

//MiniLogger is a interface below:
type MiniLogger interface {
	Debug(v ...interface{}) MiniLogger
	Info(v ...interface{}) MiniLogger
	Warn(v ...interface{}) MiniLogger
	Error(v ...interface{}) MiniLogger
	Fatal(v ...interface{})
	Debugf(format string, v ...interface{}) MiniLogger
	Infof(format string, v ...interface{}) MiniLogger
	Warnf(format string, v ...interface{}) MiniLogger
	Errorf(format string, v ...interface{}) MiniLogger
	Fatalf(format string, v ...interface{})
	Debugln(v ...interface{}) MiniLogger
	Infoln(v ...interface{}) MiniLogger
	Warnln(v ...interface{}) MiniLogger
	Errorln(v ...interface{}) MiniLogger
	Fatalln(v ...interface{})
	AddWriter(w Writer, levels uint8) MiniLogger
	Safe() MiniLogger
	Unsafe() MiniLogger
	With(fields Fields) MiniLogger
}
</pre>
