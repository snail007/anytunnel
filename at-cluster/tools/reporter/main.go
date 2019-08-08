package main

import (
	utils "anytunnel/at-common"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codeskyblue/kexec"
)

func main() {

	host := flag.String("h", "127.0.0.1", "api host")
	port := flag.Int("p", 37081, "api port")
	interval := flag.Int("i", 3, "every interval seconds to reporting")
	interfaceName := flag.String("e", "eth0", "interface name")
	failSleep := flag.Int("s", 5, "sleep seconds when fail to reporting")
	pname := flag.String("n", "at-cluster", "program name to counting connections")
	flag.Parse()
	sleep := time.Duration(*failSleep)
	lastbytes := 0
	for {
		output, err := kexec.CommandString("netstat -atn|grep -v \"LISTEN\" | wc -l").CombinedOutput()
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * sleep)
			continue
		}
		sysCount, err := strconv.Atoi(strings.Trim(string(output), "\r\n"))
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * sleep)
			continue
		}
		output, err = kexec.CommandString(`netstat -atnp 2>/dev/null |grep -v "LISTEN" |grep ` + *pname + `| wc -l`).CombinedOutput()
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * sleep)
			continue
		}
		programCount, err := strconv.Atoi(strings.Trim(string(output), "\r\n"))
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * sleep)
			continue
		}
		output, err = kexec.CommandString(`cat /proc/net/dev |grep ` + *interfaceName).CombinedOutput()
		if err != nil {
			fmt.Println("exec cat /proc/net/dev ERR:", err)
			time.Sleep(time.Second * sleep)
			continue
		}
		arr := strings.Fields(string(output))
		bandwidth := 0
		if len(arr) > 1 {
			bytesIn, err := strconv.Atoi(arr[1])
			if err != nil {
				fmt.Println(err)
				time.Sleep(time.Second * sleep)
				continue
			}
			bytesOut, err := strconv.Atoi(arr[2])
			if err != nil {
				fmt.Println(err)
				time.Sleep(time.Second * sleep)
				continue
			}
			if lastbytes == 0 {
				lastbytes = bytesOut + bytesIn
			}
			bandwidth = (bytesOut + bytesIn - lastbytes) / (*interval)
			lastbytes = bytesOut + bytesIn
		}

		url := fmt.Sprintf("https://%s:%d/cluster/report", *host, *port)
		values := map[string]string{
			"sys_conn_number":    strconv.Itoa(sysCount),
			"tunnel_conn_number": strconv.Itoa(programCount),
			"bandwidth":          fmt.Sprintf("%d", bandwidth),
		}
		body, code, err := utils.HttpPost(url, values, nil)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second * sleep)
			continue
		}
		if code != 204 {
			fmt.Println("ERR:", code, body)
			time.Sleep(time.Second * sleep)
			continue
		}
		fmt.Println("report ok ,", sysCount, programCount, bandwidth)
		time.Sleep(time.Second * time.Duration(*interval))
	}
}
