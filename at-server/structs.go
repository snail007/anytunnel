package main

import (
	utils "anytunnel/at-common"
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
	"time"
)

//###############ServerUserConn##################

type ServerUserConn struct {
	TunnelID     uint64
	ConnectionID uint64
	Connection   *tls.Conn
}

type ServerUserConnPool struct {
	pool utils.ConcurrentMap
}

func NewServerUserConnPool() ServerUserConnPool {
	return ServerUserConnPool{
		pool: utils.NewConcurrentMap(),
	}
}
func (m *ServerUserConnPool) IntToString(id uint64) string {
	return strconv.FormatUint(id, 10)
}
func (t *ServerUserConnPool) Put(tunnelID, connectionID uint64, conn *tls.Conn) {
	tid := t.IntToString(tunnelID)
	cid := t.IntToString(connectionID)
	if _, e := t.pool.Get(tid); !e {
		mp := utils.NewConcurrentMap()
		t.pool.Set(tid, &mp)
	}
	m, _ := t.pool.Get(tid)
	m.(*utils.ConcurrentMap).Set(cid, ServerUserConn{
		Connection:   conn,
		TunnelID:     tunnelID,
		ConnectionID: tunnelID,
	})
}
func (t *ServerUserConnPool) Get(tunnelID, connectionID uint64) (conn ServerUserConn, err error) {
	tid := t.IntToString(tunnelID)
	cid := t.IntToString(connectionID)
	_, e := t.pool.Get(tid)
	if !e {
		err = fmt.Errorf("tunnel %d not exists", tunnelID)
		return
	}
	m, _ := t.pool.Get(tid)

	v1, e1 := m.(*utils.ConcurrentMap).Get(cid)
	if !e1 {
		err = fmt.Errorf("server user connection %d not exists", connectionID)
		return
	}
	return v1.(ServerUserConn), err
}
func (t *ServerUserConnPool) Exists(tunnelID uint64, connectionID uint64) bool {
	tid := t.IntToString(tunnelID)
	cid := t.IntToString(connectionID)
	_, e := t.pool.Get(tid)
	if e {
		m, _ := t.pool.Get(tid)
		_, e = m.(*utils.ConcurrentMap).Get(cid)
	}
	return e
}
func (t *ServerUserConnPool) DeleteAll(tunnelID uint64) {
	tid := t.IntToString(tunnelID)
	_, e := t.pool.Get(tid)
	if !e {
		return
	}
	m, _ := t.pool.Get(tid)
	for _, v := range m.(*utils.ConcurrentMap).Keys() {
		el, _ := m.(*utils.ConcurrentMap).Get(v)
		(*(el.(ServerUserConn)).Connection).Close()
	}
	t.pool.Remove(tid)
}
func (t *ServerUserConnPool) Delete(tunnelID uint64, connectionID uint64) {
	tid := t.IntToString(tunnelID)
	cid := t.IntToString(connectionID)
	v1, e := t.pool.Get(tid)
	if !e {
		return
	}
	c, e := v1.(*utils.ConcurrentMap).Get(cid)
	if !e {
		return
	}
	(*(c.(ServerUserConn)).Connection).Close()
	v1.(*utils.ConcurrentMap).Remove(cid)
	if v1.(*utils.ConcurrentMap).IsEmpty() {
		t.pool.Remove(tid)
	}
}

//###############ListenerMap##################
type ListenerMap struct {
	data utils.ConcurrentMap
}
type ListenerTimeout struct {
	Listener *net.Listener
	LastTime int64
}
type UDPListenerTimeout struct {
	Listener *net.UDPConn
	LastTime int64
}

func NewListenerMap() ListenerMap {
	cm := ListenerMap{
		data: utils.NewConcurrentMap(),
	}
	return cm
}
func (m *ListenerMap) IntToString(id uint64) string {
	return strconv.FormatUint(id, 10)
}
func (m *ListenerMap) Get(id uint64) (conn *net.Listener, exists bool) {
	mid := m.IntToString(id)
	_conn, exists := m.data.Get(mid)
	if exists {
		conn = (_conn).(ListenerTimeout).Listener
	}
	return
}
func (m *ListenerMap) GetUDP(id uint64) (conn *net.UDPConn, exists bool) {
	mid := m.IntToString(id)
	_conn, exists := m.data.Get(mid)
	if exists {
		conn = (_conn).(UDPListenerTimeout).Listener
	}
	return
}
func (m *ListenerMap) Put(id uint64, Listener interface{}) {
	mid := m.IntToString(id)
	if v, ok := Listener.(*net.Listener); ok {
		m.data.Set(mid, ListenerTimeout{
			Listener: v,
			LastTime: time.Now().Unix(),
		})
	} else {
		m.data.Set(mid, UDPListenerTimeout{
			Listener: Listener.(*net.UDPConn),
			LastTime: time.Now().Unix(),
		})
	}

	return
}
func (m *ListenerMap) Delete(id uint64) {
	mid := m.IntToString(id)
	if _conn, exists := m.data.Get(mid); exists {
		if _, ok := _conn.(ListenerTimeout); ok {
			(*(_conn.(ListenerTimeout)).Listener).Close()
		} else {
			(*(_conn.(UDPListenerTimeout)).Listener).Close()
		}

	}
	m.data.Remove(mid)
}

//#################ip conn counter################
type IPConnCounter struct {
	data         utils.ConcurrentMap
	blockSeconds int
	maxConnCount int
}
type IPConnCounterItem struct {
	LastTime int64
	Count    int
}

func NewIPConnCounter(maxConnCount, blockSeconds int) (ic IPConnCounter) {
	ic = IPConnCounter{
		data:         utils.NewConcurrentMap(),
		blockSeconds: blockSeconds,
		maxConnCount: maxConnCount - 1,
	}
	ic.gc()
	return
}
func (ic *IPConnCounter) Check(ip string) bool {
	_item, ok := ic.data.Get(ip)
	if ok {
		item := _item.(IPConnCounterItem)
		if item.Count > ic.maxConnCount {
			if time.Now().Unix()-item.LastTime > int64(ic.blockSeconds) {
				ic.data.Remove(ip)
				return true
			}
			return false
		}
	}
	ic.data.Upsert(ip, nil, func(exist bool, valueInMap interface{}, newValue interface{}) (res interface{}) {
		var _valueInMap IPConnCounterItem
		if exist {
			_valueInMap = valueInMap.(IPConnCounterItem)
			_valueInMap.Count = _valueInMap.Count + 1
			_valueInMap.LastTime = time.Now().Unix()
		} else {
			_valueInMap = IPConnCounterItem{
				Count:    1,
				LastTime: time.Now().Unix(),
			}
		}
		return _valueInMap
	})
	return true
}
func (ic *IPConnCounter) gc() {
	go func() {
		for {
			for k, v := range ic.data.Items() {
				item := v.(IPConnCounterItem)
				if time.Now().Unix()-item.LastTime > int64(ic.blockSeconds) {
					ic.data.Remove(k)
				}
			}
			time.Sleep(time.Second * 300)
		}
	}()
}
