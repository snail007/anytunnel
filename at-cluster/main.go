package main

import (
	utils "anytunnel/at-common"
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/snail007/mini-logger"
)

var (
	poolTunnel               ClusterTunnelPool
	poolServerControlChannel ClusterServerControlChannelPool
	poolClientControlChannel ClusterClientControlChannelPool
	serverConns              *ConnMap
	csStatus                 *CSStatus
	trafficCounter           TunnelTrafficCounter
	serverTunnelPool         ServerTunnelPool
	dataRecovery             PortData
)

const SERVER_CONN_IDLE_SECONDS = 5

func init() {
	poster()
	err := initConfig()
	if err != nil {
		fmt.Printf("init config fail ,ERR:%s", err)
		return
	}
	initLog()
	initHttp(func(err error) {
		if err != nil {
			fmt.Printf("init http api fail ,ERR:%s", err)
			os.Exit(100)
		}
	})
	//fill
	poolTunnel = NewClusterTunnelPool()
	poolServerControlChannel = NewClusterServerControlChannelPool()
	poolClientControlChannel = NewClusterClientControlChannelPool()
	serverConns = NewConnMap()
	csStatus = NewCSStatus()
	trafficCounter = NewTunnelTrafficCounter()
	serverTunnelPool = NewServerTunnelPool()
	dataRecovery = NewPortData()
	initTrafficReporter()
}

func main() {
	defer func() {
		e := recover()
		if e != nil {
			log.Errorf("Exit ERR:%s", e)
		}
		logger.Flush()
	}()
	var err error
	for _, ip := range cfg.GetStringSlice("port.ip-data") {
		lnData := utils.NewServerChannel(ip, cfg.GetInt("port.conns"))
		err = lnData.ListenTls(dataConnCallback)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Infof("listening on %s for connections", (*lnData.Listener).Addr().String())
	}
	for _, ip := range cfg.GetStringSlice("port.ip-control") {
		lnControl := utils.NewServerChannel(ip, cfg.GetInt("port.control"))
		err = lnControl.ListenTls(controlConnCallback)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Infof("listening on %s for control", (*lnControl.Listener).Addr().String())
	}
	select {}
}

func dataConnCallback(mch utils.MessageChannel, conn net.Conn) {
	log.Debugf("new conn from %s", conn.RemoteAddr())
	reader := bufio.NewReader(conn)
	var timeout = 300
	var typ uint8
	var tunnleid uint64
	var connid uint64
	var err error
	okchn := make(chan bool, 1)
	go func() {
		err = binary.Read(reader, binary.LittleEndian, &typ)
		if err != nil {
			return
		}
		err = binary.Read(reader, binary.LittleEndian, &tunnleid)
		if err != nil {
			return
		}
		err = binary.Read(reader, binary.LittleEndian, &connid)
		if err != nil {
			return
		}
		okchn <- true
	}()
	select {
	case <-okchn:
	case <-time.After(time.Millisecond * time.Duration(timeout)):
		conn.SetDeadline(time.Now().Add(time.Millisecond))
		conn.Close()
		log.Debugf("read connection info timeout, %s, from:%s", err, conn.RemoteAddr())
		return
	}
	tunnel, err := poolTunnel.Get(tunnleid)
	if err != nil {
		conn.SetDeadline(time.Now().Add(time.Millisecond))
		conn.Close()
		log.Debugf("server connection %d - %d created fail ,err:%s", tunnleid, connid, err)
		return
	}
	if typ == utils.CS_SERVER {
		clientControl, err := poolClientControlChannel.Get(tunnel.ClientToken)
		if err != nil {
			conn.SetDeadline(time.Now().Add(time.Millisecond))
			conn.Close()
			log.Debugf("server connection %d - %d created fail ,err:%s", tunnleid, connid, err)
			return
		}
		serverConns.Put(tunnel.ServerToken, tunnel.TunnelID, connid, &conn)
		msgClientOpenConnection := utils.MsgClientOpenConnection{
			Msg:         utils.Msg{MsgType: utils.MSG_CLIENT_OPEN_CONNECTION},
			TunnelID:    tunnleid,
			ConnectinID: connid,
			LocalHost:   tunnel.ClientLocalHost,
			LocalPort:   tunnel.ClientLocalPort,
			Protocol:    tunnel.Protocol,
		}
		err = clientControl.ClientMessageChannel.Write(msgClientOpenConnection)
		if err != nil {
			log.Warnf("MsgClientOpenConnection write fail,%s", err)
		}
		log.Debugf("server connection %d - %d created success", tunnleid, connid)
	}
	if typ == utils.CS_CLIENT {
		connServer, ok := serverConns.Get(connid)
		if !ok {
			conn.SetDeadline(time.Now().Add(time.Millisecond))
			conn.Close()
			log.Debugf("client connection %d - %d created fail , err : server connection not exists", tunnleid, connid)
			return
		}
		utils.IoBind(*connServer, conn, func(err error) {
			conn.SetDeadline(time.Now().Add(time.Millisecond))
			(*connServer).SetDeadline(time.Now().Add(time.Millisecond))
			(*connServer).Close()
			conn.Close()
			serverConns.Delete(connid)
			log.Debugf("connection %d - %d released", tunnleid, connid)
		}, func(bytesCount int, isPositive bool) {
			if isPositive {
				trafficCounter.AddPositive(tunnleid, uint64(bytesCount))
			} else {
				trafficCounter.AddNegative(tunnleid, uint64(bytesCount))
			}
		}, tunnel.BytesPerSec)
		serverConns.ClearTimeout(connid)
		log.Debugf("connection %d - %d created success", tunnleid, connid)
	}
	return
}

func controlConnCallback(mch utils.MessageChannel, conn net.Conn) {
	var loginMsg utils.MsgLogin
	mch.DoServe(func(err error) {
		if loginMsg.IsServer() {
			g, e := poolServerControlChannel.Get(loginMsg.Token)
			if e == nil {
				if mch.RemoteAddr().String() == g.ServerMessageChannel.RemoteAddr().String() {
					poolServerControlChannel.Delete(loginMsg.Token)
				}
			}
			csStatus.ServerOffline(loginMsg.Token, mch.RemoteAddr())
			serverTunnelPool.DeleteServer(loginMsg.Token)
		}
		if loginMsg.IsClient() {
			g, e := poolClientControlChannel.Get(loginMsg.Token)
			if e == nil {
				if mch.RemoteAddr().String() == g.ClientMessageChannel.RemoteAddr().String() {
					poolClientControlChannel.Delete(loginMsg.Token)
				}
			}
			csStatus.ClientOffline(loginMsg.Token, mch.RemoteAddr())
		}
		log.Debugf("%s offline , %s ,token:%s", utils.GetCSTypeString(loginMsg.CSType), err, loginMsg.Token)
	})
	err := mch.ReadTimeout(&loginMsg, 600)
	if err != nil {
		log.Debugf("read login message error ,ERR:%s", err)
		mch.CloseConn()
		return
	}
	if e := login(&mch, loginMsg); e != nil {
		utils.Response(&mch, false, e.Error())
		time.AfterFunc(time.Second*3, func() {
			mch.CloseConn()
		})
		return
	}
	err = utils.Response(&mch, true, "")
	if err != nil {
		log.Warnf("write login response fail , ERR:%s", err)
		mch.CloseConn()
		return
	}
	//登录成功
	if loginMsg.IsClient() {
		//一个token只能登录一个client,之前的会被挤下线
		poolClientControlChannel.Delete(loginMsg.Token)
		poolClientControlChannel.Set(ClusterClientControlChannel{
			ClientToken:          loginMsg.Token,
			ClientMessageChannel: &mch,
		})
		//状态上报
		csStatus.ClientOnline(loginMsg.Token, mch.RemoteAddr())
	}
	if loginMsg.IsServer() {
		//一个token只能登录一个server,之前的会被挤下线
		poolServerControlChannel.Delete(loginMsg.Token)
		poolServerControlChannel.Set(ClusterServerControlChannel{
			ServerToken:          loginMsg.Token,
			ServerMessageChannel: &mch,
		})
		//状态上报
		csStatus.ServerOnline(loginMsg.Token, mch.RemoteAddr())
		//server之前打开的端口恢复
		go func() {
			time.Sleep(time.Second * 3)
			dataRecovery.Rcovery(loginMsg.Token)
		}()
	}
	log.Infof("%s online %s, token:%s", utils.GetCSTypeString(loginMsg.CSType), conn.RemoteAddr(), loginMsg.Token)

}
func login(channel *utils.MessageChannel, msg utils.MsgLogin) (err error) {
	log.Infof("%s login check FROM:%s", msg.CSTypeString(), channel.RemoteAddr())
	//check form auth file
	ok := false
	if msg.IsServer() {
		_, ok = tokenServerMap[msg.Token]
	} else {
		_, ok = tokenClientMap[msg.Token]
	}
	if ok {
		log.Debugf("%s auth success from auth file,%s", msg.CSTypeString(), channel.RemoteAddr())
		return
	}
	log.Infof("%s auth fail from auth file,%s", msg.CSTypeString(), channel.RemoteAddr())
	url := cfg.GetString("url.auth")
	if url == "" {
		err = fmt.Errorf("token error")
		return
	}
	if strings.Contains(url, "?") {
		url += "&"
	} else {
		url += "?"
	}
	typ := "server"
	if msg.IsClient() {
		typ = "client"
	}
	addr := channel.RemoteAddr().String()
	ip := addr[0:strings.Index(addr, ":")]
	url += fmt.Sprintf("token=%s&type=%s&ip=%s", msg.Token, typ, ip)
	var code int
	var tryCount = 0
	var body []byte
	for tryCount <= cfg.GetInt("url.fail-retry") {
		tryCount++
		body, code, err = HttpGet(url)
		if err == nil && code == cfg.GetInt("url.success-code") {
			break
		} else if err != nil {
			log.Infof("%s auth fail from auth url %s,resonse err:%s , %s", msg.CSTypeString(), url, err, channel.RemoteAddr())
			err = fmt.Errorf("auth fail from api")
		} else {
			if len(body) > 0 {
				err = fmt.Errorf(string(body[0:100]))
			} else {
				err = fmt.Errorf("token error")
			}
			log.Infof("%s auth fail from auth url %s,resonse code: %d, except: %d , %s , %s", msg.CSTypeString(), url, code, cfg.GetInt("url.success-code"), string(body), channel.RemoteAddr())
		}
		if err != nil && tryCount <= cfg.GetInt("url.fail-retry") {
			time.Sleep(time.Second * time.Duration(cfg.GetInt("url.fail-wait")))
		}
	}
	if err != nil {
		return
	}
	log.Infof("%s auth success from auth url, %s", msg.CSTypeString(), channel.RemoteAddr())
	return
}
