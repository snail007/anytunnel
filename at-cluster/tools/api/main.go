package main

import (
	utils "anytunnel/at-common"
	"flag"
	"fmt"
)

func main() {
	host := flag.String("host", "127.0.0.1", "cluster's host")
	port := flag.Int("port", 37080, "cluster's api port")
	path := flag.String("path", "/traffic/count", "request path to api")
	list := flag.Bool("list", false, "show all api path")
	flag.Parse()
	if *list {
		fmt.Println(`1.open a port
	/port/open/:TunnelID/:ServerToken/:ServerBindIP/:ServerListenPort/:ClientToken/:ClientLocalHost/:ClientLocalPort/:Protocol/:BytesPerSec
	notice : 
	BytesPerSec : bytes/seconds , zero means no limit
	Protocol : 1 is tcp , 2 is udp
2.close a port
	/port/close/:TunnelID
3.port status
	/port/close/:TunnelID
4.strip a server offline
	/server/offline/:ServerToken
5.strip a clinet offline
	/client/offline/:ClientToken
6.get all tunnel traffic statistics
	/traffic/count
	result notice:
	key : TunnelID
	negative : download bytes
	positive : upload bytes
7.get a tunnel traffic statistics
	/traffic/count/:TunnelID
		`)
		return
	}
	url := fmt.Sprintf("https://%s:%d%s", *host, *port, *path)
	body, code, err := utils.HttpGet(url)
	if err != nil {
		fmt.Printf("ERR:%s", err)
	} else {
		b := string(body)
		fmt.Println(code)
		fmt.Println(b)
	}
}
