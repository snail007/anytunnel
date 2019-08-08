package main

import (
	utils "anytunnel/at-common"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	ui "gopkg.in/gizak/termui.v2"
)

type Msg struct {
	Code    int
	Message string
	Data    utils.TrafficStatistics
}

var lastStatistics utils.TrafficStatistics

func main() {

	host := flag.String("host", "127.0.0.1", "cluster's host")
	port := flag.Int("port", 37080, "cluster's api port")
	flag.Parse()
	var err error
	url := fmt.Sprintf("https://%s:%d/traffic/count", *host, *port)

	err = ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()
	// handle key q pressing
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		// press q to quit
		ui.StopLoop()
	})
	var first = true

	var sleep = 2
	go func() {
		for {
			statistics, err := getData(url)
			if first {
				first = false
				lastStatistics = statistics
			}
			if err != nil {
				fmt.Printf("ERR:%s\n", err)
				os.Exit(0)
			}
			speedUpload := humanize.Bytes((statistics.Total.UploadBytes-lastStatistics.Total.UploadBytes)/uint64(sleep)) + "/s"
			speedDownload := humanize.Bytes((statistics.Total.DownloadBytes-lastStatistics.Total.DownloadBytes)/uint64(sleep)) + "/s"
			speedTotal := humanize.Bytes((statistics.Total.UploadBytes+statistics.Total.DownloadBytes-lastStatistics.Total.DownloadBytes-lastStatistics.Total.UploadBytes)/uint64(sleep)) + "/s"

			ui.Clear()
			table2 := ui.NewTable()
			table2.FgColor = ui.ColorYellow
			table2.BgColor = ui.ColorBlack
			table2.BorderFg = ui.ColorYellow
			table2.TextAlign = ui.AlignCenter
			table2.Separator = true
			table2.Border = true
			table2.Rows = [][]string{
				[]string{"Tunnels", "Servers", "Clients", "Connections", "Upload", "Download", "Total"},
				[]string{strconv.Itoa(statistics.Total.Tunnels),
					strconv.Itoa(statistics.Total.Servers),
					strconv.Itoa(statistics.Total.Clients),
					strconv.Itoa(statistics.Total.Connections),
					strconv.FormatUint(statistics.Total.UploadBytes, 10) + " / " + humanize.Bytes(statistics.Total.UploadBytes) + " (" + speedUpload + ")",
					strconv.FormatUint(statistics.Total.DownloadBytes, 10) + " / " + humanize.Bytes(statistics.Total.DownloadBytes) + " (" + speedDownload + ")",
					strconv.FormatUint(statistics.Total.TotalBytes, 10) + " / " + humanize.Bytes(statistics.Total.TotalBytes) + " (" + speedTotal + ")",
				},
			}
			table2.Analysis()
			table2.SetSize()

			table3 := ui.NewTable()
			table3.FgColor = ui.ColorWhite
			table3.BgColor = ui.ColorDefault
			table3.TextAlign = ui.AlignCenter
			table3.Separator = true
			table3.Border = true
			table3.Rows = [][]string{
				[]string{"TunnelID", "Connections", "Upload", "Download", "Total", "Limit"},
			}
			for _, total := range sortTunnels(statistics.Traffic, 20) {
				TunnelID := strconv.FormatUint(total["TunnelID"], 10)
				speedUpload := ""
				speedDownload := ""
				speedTotal := ""
				var connections uint64
				if v, ok := lastStatistics.Traffic[TunnelID]; ok {
					speedUpload = humanize.Bytes((total["positive"]-v["positive"])/uint64(sleep)) + "/s"
					speedDownload = humanize.Bytes((total["negative"]-v["negative"])/uint64(sleep)) + "/s"
					_speedTotal := (total["positive"] + total["negative"] - v["negative"] - v["positive"]) / uint64(sleep)
					speedTotal = humanize.Bytes(_speedTotal) + "/s"
					if total["limit"] > 0 {
						connections = _speedTotal / total["limit"]
					}
					if _speedTotal > 0 && connections == 0 {
						connections = 1
					}

				} else {
					speedUpload = ""
					speedDownload = ""
					speedTotal = ""
				}
				table3.Rows = append(table3.Rows, []string{TunnelID,
					strconv.FormatUint(connections, 10),
					strconv.FormatUint(total["positive"], 10) + " / " + humanize.Bytes(total["positive"]) + " (" + speedUpload + ")",
					strconv.FormatUint(total["negative"], 10) + " / " + humanize.Bytes(total["negative"]) + " (" + speedDownload + ")",
					strconv.FormatUint(total["positive"]+total["negative"], 10) + " / " + humanize.Bytes(total["positive"]+total["negative"]) + " (" + speedTotal + ")",
					humanize.Bytes(total["limit"]) + "/s",
				})
			}
			table3.Analysis()
			table3.SetSize()

			ui.Body.Rows = []*ui.Row{}
			ui.Body.AddRows(
				ui.NewRow(
					ui.NewCol(12, 0, table2),
				),
				ui.NewRow(
					ui.NewCol(12, 0, table3),
				),
			)
			ui.Body.Align()
			ui.Render(ui.Body)
			lastStatistics = statistics
			time.Sleep(time.Second * time.Duration(sleep))
		}
	}()

	ui.Loop()
}
func getData(url string) (statistics utils.TrafficStatistics, err error) {
	body, code, err := utils.HttpGet(url)
	if err != nil {
		err = fmt.Errorf("ERR:%s", err)
		return
	}
	if code != 200 {
		err = fmt.Errorf("ERR: http code %d,body : %s", code, string(body))
		return
	}

	var msg Msg
	err = json.Unmarshal(body, &msg)
	if err != nil {
		err = fmt.Errorf("ERR:%s", err)
		return
	}
	if msg.Code != 1 {
		err = fmt.Errorf("ERR:%s", msg.Message)
		return
	}
	statistics = msg.Data
	return
}

type DirRange []uint64

func (a DirRange) Len() int           { return len(a) }
func (a DirRange) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DirRange) Less(i, j int) bool { return a[i] > a[j] }

func sortTunnels(data map[string]map[string]uint64, limit int) []map[string]uint64 {
	sortMap := map[string]uint64{}
	sortUint := DirRange{}
	for k, v := range data {
		var p uint64 = 0
		if ov, ok := lastStatistics.Traffic[k]; ok {
			p = ov["positive"] + ov["negative"]
		}
		total := (v["positive"] + v["negative"]) - p
		sortMap[k] = total
		sortUint = append(sortUint, total)
	}
	sortUintUniqueMap := map[uint64]bool{}
	for _, v := range sortUint {
		sortUintUniqueMap[v] = true
	}
	uniqueSortUint := DirRange{}
	for k := range sortUintUniqueMap {
		uniqueSortUint = append(uniqueSortUint, k)
	}
	sort.Sort(uniqueSortUint)
	reslut := []map[string]uint64{}
	i := 0
	for _, v := range uniqueSortUint {
		for k, v1 := range sortMap {
			if v == v1 {
				id, _ := strconv.ParseUint(k, 10, 64)
				newItem := map[string]uint64{
					"positive": data[k]["positive"],
					"negative": data[k]["negative"],
					"TunnelID": id,
					"limit":    data[k]["limit"],
				}
				reslut = append(reslut, newItem)
			}
		}
		i++
		if i > limit {
			break
		}
	}
	return reslut
}
