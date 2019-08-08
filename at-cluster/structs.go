package main

import (
	utils "anytunnel/at-common"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	CSSTATUS_TYPE_SERVER    = "server"
	CSSTATUS_TYPE_CLIENT    = "client"
	CSSTATUS_ACTION_ONLINE  = "online"
	CSSTATUS_ACTION_OFFLINE = "offline"
)

type CSStatus struct {
	url    string
	bufchn chan CSStatusAction
}
type CSStatusAction struct {
	Type   string
	Token  string
	Action string
	IP     string
}

func NewCSStatus() *CSStatus {
	s := CSStatus{
		url:    cfg.GetString("url.status"),
		bufchn: make(chan CSStatusAction, 500),
	}
	s.init()
	return &s
}
func (cs *CSStatus) init() {
	go func() {
		for {
			action := <-cs.bufchn
			var tryCount = 0
			for tryCount <= cfg.GetInt("url.fail-retry") {
				tryCount++
				_url := utils.UrlArgs(cs.url, fmt.Sprintf("type=%s&token=%s&action=%s&ip=%s", action.Type, action.Token, action.Action, action.IP))
				body, code, err := HttpGet(_url)
				if err == nil && code == cfg.GetInt("url.success-code") {
					break
				} else if err != nil {
					log.Warnf("report online/offline status fail to url %s, err: %s", cs.url, err)
				} else {
					err = fmt.Errorf("token error")
					log.Warnf("report online/offline status fail to url %s, code: %d, except: %d ,body:%s", cs.url, code, cfg.GetInt("url.success-code"), string(body))
				}
				if err != nil && tryCount <= cfg.GetInt("url.fail-retry") {
					time.Sleep(time.Second * time.Duration(cfg.GetInt("url.fail-wait")))
				}
			}
		}
	}()
}
func (cs *CSStatus) ClientOnline(token string, addr net.Addr) {
	cs.pushStatus(CSSTATUS_TYPE_CLIENT, token, CSSTATUS_ACTION_ONLINE, addr)
}
func (cs *CSStatus) ClientOffline(token string, addr net.Addr) {
	cs.pushStatus(CSSTATUS_TYPE_CLIENT, token, CSSTATUS_ACTION_OFFLINE, addr)
}
func (cs *CSStatus) ServerOnline(token string, addr net.Addr) {
	cs.pushStatus(CSSTATUS_TYPE_SERVER, token, CSSTATUS_ACTION_ONLINE, addr)
}
func (cs *CSStatus) ServerOffline(token string, addr net.Addr) {
	cs.pushStatus(CSSTATUS_TYPE_SERVER, token, CSSTATUS_ACTION_OFFLINE, addr)
}
func (cs *CSStatus) pushStatus(typ, token, action string, addr net.Addr) {
	if cs.url == "" {
		return
	}
	_addr := addr.String()
	ip := _addr[0:strings.Index(_addr, ":")]
	cs.bufchn <- CSStatusAction{
		Action: action,
		Type:   typ,
		Token:  token,
		IP:     ip,
	}
}

//###################ConnMap######################
type ConnMap struct {
	data utils.ConcurrentMap
}
type ConnItem struct {
	ServerToken string
	TunnelID    uint64
	Count       uint64
}
type ConnTimeout struct {
	Conn        *net.Conn
	ServerToken string
	TunnelID    uint64
	LastTime    int64
}

func NewConnMap() *ConnMap {
	cm := ConnMap{
		data: utils.NewConcurrentMap(),
	}
	cm.gc()
	return &cm
}
func (m *ConnMap) IntToString(id uint64) string {
	return strconv.FormatUint(id, 10)
}
func (m *ConnMap) Get(id uint64) (conn *net.Conn, exists bool) {
	mid := m.IntToString(id)
	_conn, exists := m.data.Get(mid)
	if exists {
		conn = (_conn).(ConnTimeout).Conn
	}
	return
}
func (m *ConnMap) Put(serverToken string, tunnleID, connid uint64, conn *net.Conn) {
	mid := m.IntToString(connid)
	m.data.Set(mid, ConnTimeout{
		Conn:        conn,
		LastTime:    time.Now().Unix(),
		ServerToken: serverToken,
		TunnelID:    tunnleID,
	})
	//b, _ := m.data.MarshalJSON()
	//fmt.Println("put:", serverToken, tunnleID, connid, m.GetTunnelConnCountMap(), string(b))
	return
}

func (m *ConnMap) Delete(id uint64) {
	mid := m.IntToString(id)
	_conn, exists := m.data.Get(mid)
	if exists {
		(*(_conn).(ConnTimeout).Conn).Close()
	}
	m.data.Remove(mid)
}
func (m *ConnMap) ClearTimeout(id uint64) {
	mid := m.IntToString(id)
	_conn, exists := m.data.Get(mid)
	if exists {
		_c := _conn.(ConnTimeout)
		_c.LastTime = 0
		m.data.Set(mid, _c)
	}
}
func (m *ConnMap) GetTunnelConnCountMap() map[uint64]ConnItem {
	countMap := map[uint64]ConnItem{}
	for _, k := range m.data.Keys() {
		v, _ := m.data.Get(k)
		c := v.(ConnTimeout)
		item, ok := countMap[c.TunnelID]
		if !ok {
			countMap[c.TunnelID] = ConnItem{
				Count:       1,
				ServerToken: c.ServerToken,
				TunnelID:    c.TunnelID,
			}
		} else {
			item.Count = item.Count + 1
			countMap[c.TunnelID] = item
		}
	}
	return countMap
}
func (m *ConnMap) gc() {
	go func() {
		for {
			for _, k := range m.data.Keys() {
				v, _ := m.data.Get(k)
				c := v.(ConnTimeout)
				if c.LastTime > 0 && time.Now().Unix()-c.LastTime > SERVER_CONN_IDLE_SECONDS {
					(*c.Conn).Close()
					m.data.Remove(k)
				}
			}
			time.Sleep(60 * time.Second)
		}
	}()
}

//###############ClusterTunnelPool##################
type ClusterTunnel struct {
	TunnelID         uint64
	ServerToken      string
	ServerBindIP     string
	ServerListenPort int
	ClientToken      string
	ClientLocalHost  string
	ClientLocalPort  int
	Protocol         int
	BytesPerSec      float64
}
type ClusterTunnelPool struct {
	pool utils.ConcurrentMap
}

func NewClusterTunnelPool() ClusterTunnelPool {
	return ClusterTunnelPool{
		pool: utils.NewConcurrentMap(),
	}
}
func (m *ClusterTunnelPool) IntToString(id uint64) string {
	return strconv.FormatUint(id, 10)
}
func (t *ClusterTunnelPool) Set(tunnel ClusterTunnel) {
	mid := t.IntToString(tunnel.TunnelID)
	t.pool.Set(mid, tunnel)
}
func (t *ClusterTunnelPool) Get(tunnelID uint64) (tunnel ClusterTunnel, err error) {
	mid := t.IntToString(tunnelID)
	v, e := t.pool.Get(mid)
	if !e {
		err = fmt.Errorf("tunnel %d not exists", tunnelID)
		return
	}
	return v.(ClusterTunnel), err
}
func (t *ClusterTunnelPool) Exists(tunnelID uint64) bool {
	mid := t.IntToString(tunnelID)
	return t.pool.Has(mid)
}
func (t *ClusterTunnelPool) Delete(tunnelID uint64) {
	mid := t.IntToString(tunnelID)
	t.pool.Remove(mid)
}

//###############ClusterServerControlChannelPool##################
type ClusterServerControlChannel struct {
	ServerToken          string
	ServerMessageChannel *utils.MessageChannel
}
type ClusterServerControlChannelPool struct {
	pool utils.ConcurrentMap
}

func NewClusterServerControlChannelPool() ClusterServerControlChannelPool {
	return ClusterServerControlChannelPool{
		pool: utils.NewConcurrentMap(),
	}
}
func (t *ClusterServerControlChannelPool) Set(serverChannel ClusterServerControlChannel) {
	t.pool.Set(serverChannel.ServerToken, serverChannel)
}
func (t *ClusterServerControlChannelPool) Get(serverToken string) (serverChannel ClusterServerControlChannel, err error) {
	v, e := t.pool.Get(serverToken)
	if !e {
		err = fmt.Errorf("ClusterServerControlChannel %s not exists", serverToken)
		return
	}
	return v.(ClusterServerControlChannel), err
}
func (t *ClusterServerControlChannelPool) Exists(serverToken string) bool {
	return t.pool.Has(serverToken)
}
func (t *ClusterServerControlChannelPool) Delete(serverToken string) {
	v, e := t.pool.Get(serverToken)
	if e {
		v.(ClusterServerControlChannel).ServerMessageChannel.CloseConn()
	}
	t.pool.Remove(serverToken)
}

//###############ClusterClientControlChannelPool##################
type ClusterClientControlChannel struct {
	ClientToken          string
	ClientMessageChannel *utils.MessageChannel
}
type ClusterClientControlChannelPool struct {
	pool utils.ConcurrentMap
}

func NewClusterClientControlChannelPool() ClusterClientControlChannelPool {
	return ClusterClientControlChannelPool{
		pool: utils.NewConcurrentMap(),
	}
}
func (t *ClusterClientControlChannelPool) Set(clientChannel ClusterClientControlChannel) {
	t.pool.Set(clientChannel.ClientToken, clientChannel)
}
func (t *ClusterClientControlChannelPool) Get(clientToken string) (clientChannel ClusterClientControlChannel, err error) {
	v, e := t.pool.Get(clientToken)
	if !e {
		err = fmt.Errorf("ClusterClientControlChannel %s not exists", clientToken)
		return
	}
	return v.(ClusterClientControlChannel), err
}
func (t *ClusterClientControlChannelPool) Exists(clientToken string) bool {
	return t.pool.Has(clientToken)
}
func (t *ClusterClientControlChannelPool) Delete(clientToken string) {
	v, e := t.pool.Get(clientToken)
	if e {
		v.(ClusterClientControlChannel).ClientMessageChannel.CloseConn()
	}
	t.pool.Remove(clientToken)
}

//
//###############counter##################
type TunnelTrafficCounter struct {
	counter utils.ConcurrentMap
	bufchn  chan TunnelTrafficEntry
}
type TunnelTrafficEntry struct {
	IsPositive bool
	TunnelID   uint64
	CountBytes uint64
}

//map[string]uint64{}
func NewTunnelTrafficCounter() TunnelTrafficCounter {
	tc := TunnelTrafficCounter{
		counter: utils.NewConcurrentMap(),
		bufchn:  make(chan TunnelTrafficEntry, 500000),
	}
	tc.init()
	return tc
}
func (m *TunnelTrafficCounter) IntToString(id uint64) string {
	return strconv.FormatUint(id, 10)
}
func (tc *TunnelTrafficCounter) init() {
	go func() {
		for {
			entry := <-tc.bufchn
			mid := tc.IntToString(entry.TunnelID)
			tc.counter.Upsert(mid, nil, func(exists bool, vold interface{}, vnew interface{}) (res interface{}) {
				if vold != nil {
					if entry.IsPositive {
						vold.(map[string]uint64)["positive"] = vold.(map[string]uint64)["positive"] + entry.CountBytes
					} else {
						vold.(map[string]uint64)["negative"] = vold.(map[string]uint64)["negative"] + entry.CountBytes
					}
					res = vold
				}
				return
			})
		}
	}()
}
func (tc *TunnelTrafficCounter) InitTunnel(tunnelID uint64) {
	mid := tc.IntToString(tunnelID)
	tc.counter.Set(mid, map[string]uint64{
		"positive": 0,
		"negative": 0,
	})
}
func (tc *TunnelTrafficCounter) DeleteTunnel(tunnelID uint64) {
	mid := tc.IntToString(tunnelID)
	tc.counter.Remove(mid)
}
func (tc *TunnelTrafficCounter) AddPositive(tunnelID, bytesCount uint64) {
	tc.bufchn <- TunnelTrafficEntry{
		TunnelID:   tunnelID,
		IsPositive: true,
		CountBytes: bytesCount,
	}
}
func (tc *TunnelTrafficCounter) AddNegative(tunnelID, bytesCount uint64) {
	tc.bufchn <- TunnelTrafficEntry{
		TunnelID:   tunnelID,
		IsPositive: false,
		CountBytes: bytesCount,
	}
}
func (tc *TunnelTrafficCounter) AllDataJSON() ([]byte, error) {
	JSON := utils.TrafficStatistics{
		Total: utils.TrafficTotal{
			Tunnels:       poolTunnel.pool.Count(),
			Servers:       poolServerControlChannel.pool.Count(),
			Clients:       poolClientControlChannel.pool.Count(),
			UploadBytes:   0,
			DownloadBytes: 0,
			TotalBytes:    0,
			Connections:   serverConns.data.Count(),
		},
	}
	tunnelCount := serverConns.GetTunnelConnCountMap()

	data := map[string]map[string]uint64{}
	for nk, nv := range tc.counter.Items() {
		tid, _ := strconv.ParseUint(nk, 10, 64)
		var count uint64
		if v, ok := tunnelCount[tid]; ok {
			count = v.Count
		}
		uploadBytes := nv.(map[string]uint64)["positive"]
		downloadBytes := nv.(map[string]uint64)["negative"]

		tunnel, _ := poolTunnel.Get(tid)
		data[nk] = map[string]uint64{
			"positive":    uploadBytes,
			"negative":    downloadBytes,
			"limit":       uint64(tunnel.BytesPerSec),
			"connections": count,
		}
		JSON.Total.UploadBytes += uploadBytes
		JSON.Total.DownloadBytes += downloadBytes
	}
	JSON.Total.TotalBytes = JSON.Total.UploadBytes + JSON.Total.DownloadBytes
	JSON.Traffic = data
	return json.Marshal(JSON)
}
func (tc *TunnelTrafficCounter) TunnelDataJSON(tunnelID uint64) ([]byte, error) {
	mid := tc.IntToString(tunnelID)
	if v, e := tc.counter.Get(mid); e {
		return json.Marshal(v)
	}
	return []byte("{}"), nil
}
func (tc *TunnelTrafficCounter) AllData() map[string]map[string]uint64 {
	data := map[string]map[string]uint64{}
	tc.counter.IterCb(func(nk string, nv interface{}) {
		if nv != nil {
			data[nk] = map[string]uint64{
				"positive": nv.(map[string]uint64)["positive"],
				"negative": nv.(map[string]uint64)["negative"],
			}
		}
	})
	return data
}

//###############ServerTunnelPool##################
type ServerTunnelPool struct {
	pool utils.ConcurrentMap
}

func NewServerTunnelPool() ServerTunnelPool {
	return ServerTunnelPool{
		pool: utils.NewConcurrentMap(),
	}
}
func (s *ServerTunnelPool) AddTunnel(ServerToken string, TunnelID uint64) {
	v, ok := s.pool.Get(ServerToken)
	if !ok {
		v = map[uint64]bool{}
	}
	_v := v.(map[uint64]bool)
	_v[TunnelID] = true
	s.pool.Set(ServerToken, _v)
}
func (s *ServerTunnelPool) DeleteServer(ServerToken string) {
	if v, ok := s.pool.Get(ServerToken); ok {
		_v := v.(map[uint64]bool)
		for TunnelID := range _v {
			poolTunnel.Delete(TunnelID)
			trafficCounter.DeleteTunnel(TunnelID)
		}
		s.pool.Remove(ServerToken)
	}
}

//###############ServerTunnelPool##################

type PortData struct {
}

func NewPortData() PortData {
	return PortData{}
}
func (d *PortData) Rcovery(ServerToken string) (err error) {
	defer func() {
		if err != nil {
			log.Warnf("recovery return an error , ERR: %s", err)
		}
	}()
	path, err := d.getPath(ServerToken)
	if err != nil {
		return
	}
	var oldContets []byte
	oldContets, err = d.fileGetContents(path)
	if err != nil {
		return
	}
	oldContetsLines := strings.Split(string(oldContets), "\n")
	// log.Debugf("maybe %d tunnels to recovery,ServerToken:%s", len(oldContetsLines), ServerToken)
	for _, line := range oldContetsLines {
		lineFields := strings.Fields(line)
		if len(lineFields) != 2 {
			continue
		}
		ip := cfg.GetStringSlice("port.ip-api")[0]
		if ip == "0.0.0.0" || ip == "" {
			ip = "127.0.0.1"
		}
		url := fmt.Sprintf("https://%s:%d%s", ip, 37080, lineFields[1])
		utils.HttpGet(url)
		// _, code, err := utils.HttpGet(url)
		// if err != nil {
		// 	log.Debugf("recovery err : %d %s %s", code, err, url)

		// } else {
		// 	log.Debugf("recovery result : %d %s", code, url)
		// }
	}
	return
}
func (d *PortData) Delete(token string, tunnelID uint64) (err error) {
	return d._StoreOrDelete(ClusterTunnel{
		TunnelID:    tunnelID,
		ServerToken: token,
	}, false)
}
func (d *PortData) DeleteServer(token string) (err error) {
	path, err := d.getPath(token)
	if err != nil {
		return
	}
	return os.Remove(path)
}
func (d *PortData) Store(t ClusterTunnel) (err error) {
	return d._StoreOrDelete(t, true)
}
func (d *PortData) _StoreOrDelete(t ClusterTunnel, isStore bool) (err error) {
	// d.l.Lock()
	// defer d.l.Unlock()
	path, err := d.getPath(t.ServerToken)
	if err != nil {
		return
	}
	cmd := ""
	if isStore {
		cmd = fmt.Sprintf("/port/open/%d/%s/%s/%d/%s/%s/%d/%d/%.0f", t.TunnelID, t.ServerToken, t.ServerBindIP, t.ServerListenPort,
			t.ClientToken, t.ClientLocalHost, t.ClientLocalPort, t.Protocol, t.BytesPerSec)
	}
	content := fmt.Sprintf("%d %s\n", t.TunnelID, cmd)
	if d.exists(path) {
		var oldContets []byte
		oldContets, err = d.fileGetContents(path)
		if err != nil {
			return
		}
		oldContetsLines := strings.Split(string(oldContets), "\n")
		newContests := []byte{}
		writer := bytes.NewBuffer(newContests)
		found := false
		for _, line := range oldContetsLines {
			lineFields := strings.Fields(line)
			if len(lineFields) != 2 {
				continue
			}
			id, _ := strconv.Atoi(lineFields[0])
			if uint64(id) == t.TunnelID {
				found = true
				if isStore {
					writer.WriteString(content)
				}
				continue
			} else {
				writer.WriteString(line + "\n")
			}
		}
		if !found && isStore {
			writer.WriteString(content)
		}
		err = d.filePutContents(path, writer.Bytes())
		if err != nil {
			return
		}
	} else if isStore {
		err = d.mkdirs(path)
		if err != nil {
			return
		}
		err = d.filePutContents(path, []byte(content))
		if err != nil {
			return
		}
	}
	return
}
func (d *PortData) getPath(token string) (path string, err error) {
	return filepath.Abs(fmt.Sprintf("%s/%d/%s", cfg.GetString("data.dir"), crc32.ChecksumIEEE([]byte(token))%10, token))
}
func (d *PortData) mkdirs(path string) (err error) {
	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	return
}
func (d *PortData) exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}
func (d *PortData) fileGetContents(filename string) ([]byte, error) {
	fp, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	reader := bufio.NewReader(fp)
	contents, _ := ioutil.ReadAll(reader)
	return contents, nil
}

func (d *PortData) filePutContents(filename string, content []byte) error {
	mode := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	fp, err := os.OpenFile(filename, mode, os.ModePerm)
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.Write(content)
	return err
}
