package main

import (
	"anytunnel/at-common/qqwry"
	"fmt"
	"os"

	logger "github.com/snail007/mini-logger"
)

func init() {
	poster()
	initConfig()
	initLog()
	initHttp(func(err error) {
		if err != nil {
			fmt.Printf("init http api fail ,ERR:%s", err)
			os.Exit(100)
		}
	})
	initDB()
	qqwry.LoadData("data/qqwry.dat")
}

func main() {
	defer func() {
		e := recover()
		if e != nil {
			log.Errorf("Exit ERR:%v", e)
		}
		logger.Flush()
	}()
	select {}
}
