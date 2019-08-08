package at_common

type TrafficTotal struct {
	Tunnels       int    `json="tunnels"`
	Servers       int    `json="servers"`
	Clients       int    `json="clinets"`
	UploadBytes   uint64 `json="uploadBytes"`
	DownloadBytes uint64 `json="downloadBytes"`
	TotalBytes    uint64 `json="totalBytes"`
	Connections   int    `json="connections"`
}
type TrafficStatistics struct {
	Total   TrafficTotal                 `json="total"`
	Traffic map[string]map[string]uint64 `json="traffic"`
}
