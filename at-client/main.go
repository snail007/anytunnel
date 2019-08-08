package main

import (
	utils "anytunnel/at-common"
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	logger "github.com/snail007/mini-logger"
)

var (
	control utils.Author
	err     error
)

func init() {
	poster()
	initConfig()
	initLog()
}
func main() {
	defer func() {
		e := recover()
		if e != nil {
			log.Errorf("Exit ERR:%s", e)
		}
		logger.Flush()
	}()
	control, err = utils.NewAuthor(cfg.GetString("token"), cfg.GetString("host"), cfg.GetInt("port.control"), utils.CSTYPE_CLIENT)
	if err != nil {
		log.Debugf("create author fail : %s", err)
		return
	}
	control.Channel.SetMsgErrorHandler(func(channel *utils.MessageChannel, msg interface{}, err error) {
		log.Warnf("message pre-process error , ERRR:%s", err)
	})
	control.Channel.RegMsg(utils.MSG_TYPE_PONG, new(utils.MsgPong), func(channel *utils.MessageChannel, msg interface{}) {
		msgPong := msg.(*utils.MsgPong)
		log.Infof("pong revecived , id : %s", msgPong.ID)
		return
	})
	control.Channel.RegMsg(utils.MSG_CLIENT_OPEN_CONNECTION, new(utils.MsgClientOpenConnection), func(channel *utils.MessageChannel, msg interface{}) {
		msgClientOpenConnection := msg.(*utils.MsgClientOpenConnection)
		if msgClientOpenConnection.Protocol == utils.TUNNEL_PROTOCOL_TCP {
			openConnection(channel, msgClientOpenConnection)
		} else {
			openUDPConnection(channel, msgClientOpenConnection)
		}
	})
	control.Channel.DoServe(func(err error) {
		//log.Fatalf("offline , disconnected from cluster , %s", err)
		log.Fatalf("offline , disconnected from cluster .")
	})
	err = control.DoControlAuth()
	if err != nil {
		log.Fatalf("login fail : %s", err)
		return
	}
	//	log.Infof("login success %s - %s", control.Channel.LocalAddr(), control.Channel.RemoteAddr())

	log.Infof("login success")
	select {}
}

var clientClusterConnPool = utils.NewConcurrentMap()

func openUDPConnection(controlChannel *utils.MessageChannel, msg *utils.MsgClientOpenConnection) {
	tunnleID := msg.TunnelID
	connid := msg.ConnectinID
	connidStr := fmt.Sprintf("%d", (*msg).ConnectinID)
	var clusterConn *tls.Conn
	_, ok := clientClusterConnPool.Get(connidStr)
	if !ok {
		_clusterConn, err := connectCluster(*msg)
		if err != nil {
			return
		}
		clusterConn = &_clusterConn
		clientClusterConnPool.Set(connidStr, clusterConn)
		log.Debugf("connection %d - %d created success", tunnleID, connid)
		go func() {
			for {
				srcAddr, body, err := utils.ReadUDPPacket(clusterConn)
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					log.Debugf("connection %d - %d released", tunnleID, connid)
					clientClusterConnPool.Remove(connidStr)
					break
				}
				func() {
					//log.Debugf("udp packet revecived:%s,%v", srcAddr, body)
					dstAddr := &net.UDPAddr{IP: net.ParseIP(msg.LocalHost), Port: msg.LocalPort}
					clientSrcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
					conn, err := net.DialUDP("udp", clientSrcAddr, dstAddr)
					if err != nil {
						log.Warnf("connect to udp %s fail,ERR:%s", dstAddr.String(), err)
						return
					}
					conn.SetDeadline(time.Now().Add(time.Second * time.Duration(cfg.GetInt("udp.timeout"))))
					_, err = conn.Write(body)
					if err != nil {
						log.Warnf("send udp packet to %s fail,ERR:%s", dstAddr.String(), err)
						return
					}
					//log.Debugf("send udp packet to %s success", dstAddr.String())
					buf := make([]byte, 512)
					len, _, err := conn.ReadFromUDP(buf)
					if err != nil {
						log.Warnf("read udp response from %s fail ,ERR:%s", dstAddr.String(), err)
						return
					}
					respBody := buf[0:len]
					//log.Debugf("revecived udp packet from %s , %v", dstAddr.String(), respBody)
					_, err = clusterConn.Write(utils.UDPPacket(srcAddr, respBody))
					if err != nil {
						log.Warnf("send udp response to cluster fail ,ERR:%s", err)
						return
					}
					//log.Debugf("send udp response to cluster success ,from:%s", dstAddr.String())
				}()

			}
		}()
	}
	return
}
func openConnection(controlChannel *utils.MessageChannel, msg *utils.MsgClientOpenConnection) {
	tunnleID := msg.TunnelID
	connid := msg.ConnectinID
	clusterConn, err := connectCluster(*msg)
	if err != nil {
		return
	}
	localConn, err := utils.Connect(msg.LocalHost, msg.LocalPort, 5000)
	if err != nil {
		log.Debugf("connection %d - %d created fail ,err:%s", tunnleID, connid, err)
		clusterConn.Close()
		return
	}
	utils.IoBind(localConn, &clusterConn, func(err error) {
		localConn.Close()
		clusterConn.Close()
		log.Debugf("connection %d - %d released", tunnleID, connid)
	}, func(bytesCount int, isPositive bool) {}, 0)
	log.Debugf("connection %d - %d created success", tunnleID, connid)
	return
}
func connectCluster(cmd utils.MsgClientOpenConnection) (clusterConn tls.Conn, err error) {
	tunnleID := cmd.TunnelID
	connid := cmd.ConnectinID
	log.Debugf("new connection %d - %d", tunnleID, connid)
	clusterConn, err = utils.TlsConnect(cfg.GetString("host"), cfg.GetInt("port.conns"), 3000)
	if err != nil {
		log.Warnf("connect to cluster fail ,err:%s", err)
		return
	}
	writer := bufio.NewWriter(&clusterConn)
	pkg := new(bytes.Buffer)
	binary.Write(pkg, binary.LittleEndian, utils.CS_CLIENT)
	binary.Write(pkg, binary.LittleEndian, tunnleID)
	binary.Write(pkg, binary.LittleEndian, connid)
	writer.Write(pkg.Bytes())
	err = writer.Flush()
	if err != nil {
		log.Warnf("connect to cluster fail ,flush err:%s", err)
		return
	}
	return
}
