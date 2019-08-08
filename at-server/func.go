package main

import (
	utils "anytunnel/at-common"
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func openUDPPort(cmd utils.MsgServerOpenPort) (sc utils.ServerChannel, err error) {
	sc = utils.NewServerChannel(cmd.BindIP, cmd.BindPort)
	sc.SetErrAcceptHandler(func(err error) {
		addr := ""
		if cmd.Protocol == utils.TUNNEL_PROTOCOL_TCP {
			addr = (*sc.Listener).Addr().String()
		} else {
			addr = (*sc.UDPListener).LocalAddr().String()
		}
		log.Debugf("%s port %s closed , ERR:%s", cmd.ProtocolString(), addr, err)
		if s, ok := serverListeners.GetUDP(cmd.TunnelID); ok {
			if (*sc.UDPListener).LocalAddr().String() == (*s).LocalAddr().String() {
				serverListeners.Delete(cmd.TunnelID)
			}
		}
		go serverListenerConns.DeleteAll(cmd.TunnelID)
	})
	err = sc.ListenUDP(func(packet []byte, localAddr, srcAddr *net.UDPAddr) {
		openUDPConn(sc, packet, localAddr, srcAddr, cmd)
	})
	if err != nil {
		return
	}
	serverListeners.Put(cmd.TunnelID, sc.UDPListener)
	return
}
func openUDPConn(sc utils.ServerChannel, packet []byte, localAddr, srcAddr *net.UDPAddr, cmd utils.MsgServerOpenPort) {
	numLocal := crc32.ChecksumIEEE([]byte(localAddr.String()))
	numSrc := crc32.ChecksumIEEE([]byte(srcAddr.String()))
	connid := uint64((numLocal/10)*10 + numSrc%10)
	var clusterConn *tls.Conn
	conn, err := serverListenerConns.Get(cmd.TunnelID, connid)
	if err != nil {
		connid, _clusterConn, err := connectCluster(connid, cmd)
		if err != nil {
			log.Warnf("connect to cluster fail for udp,ERR:%s", err)
			return
		}
		log.Debugf("connection %d - %d created success", cmd.TunnelID, connid)
		clusterConn = &_clusterConn
		//防止并发时,同时对一个 cmd.TunnelID, connid 建立连接,设置之前,杀死旧的连接,如果存在的话.
		if _, err := serverListenerConns.Get(cmd.TunnelID, connid); err == nil {
			serverListenerConns.Delete(cmd.TunnelID, connid)
		}
		serverListenerConns.Put(cmd.TunnelID, connid, clusterConn)
		go func() {
			for {
				srcAddrFromCluster, body, err := utils.ReadUDPPacket(clusterConn)
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					log.Debugf("connection %d - %d released", cmd.TunnelID, connid)
					serverListenerConns.Delete(cmd.TunnelID, connid)
					break
				}
				//log.Debugf("udp packet revecived from cluster,local:%s", srcAddrFromCluster)
				_srcAddr := strings.Split(srcAddrFromCluster, ":")
				port, _ := strconv.Atoi(_srcAddr[1])
				dstAddr := &net.UDPAddr{IP: net.ParseIP(_srcAddr[0]), Port: port}
				_, err = sc.UDPListener.WriteToUDP(body, dstAddr)
				if err != nil {
					log.Warnf("udp response to local %s fail,ERR:%s", srcAddr, err)
					continue
				}
				//log.Debugf("udp response to local %s success", srcAddr)
			}
		}()
	} else {
		clusterConn = conn.Connection
		//log.Debugf("get conn to cluster for udp success, tunnelID:%d,connid:%d", cmd.TunnelID, connid)
	}
	if err != nil {
		return
	}
	writer := bufio.NewWriter(clusterConn)
	writer.Write(utils.UDPPacket(srcAddr.String(), packet))
	err = writer.Flush()
	if err != nil {
		log.Warnf("connect to cluster fail ,flush err:%s", err)
		return
	}
	//log.Debugf("write packet %v", packet)
	return
}
func openPort(cmd utils.MsgServerOpenPort) (sc utils.ServerChannel, err error) {
	sc = utils.NewServerChannel(cmd.BindIP, cmd.BindPort)
	sc.SetErrAcceptHandler(func(err error) {
		log.Debugf("%s port %s closed , ERR:%s", cmd.ProtocolString(), (*sc.Listener).Addr(), err)
		if s, ok := serverListeners.Get(cmd.TunnelID); ok {
			if (*sc.Listener).Addr().String() == (*s).Addr().String() {
				serverListeners.Delete(cmd.TunnelID)
			}
		}
		go serverListenerConns.DeleteAll(cmd.TunnelID)
	})
	err = sc.ListenTCP(func(mch utils.MessageChannel, conn net.Conn) {
		openConn(conn, cmd)
	})
	if err != nil {
		return
	}
	serverListeners.Put(cmd.TunnelID, sc.Listener)
	return
}

func openConn(conn net.Conn, cmd utils.MsgServerOpenPort) {
	addr := conn.RemoteAddr().String()
	ip := addr[0:strings.Index(addr, ":")]
	if !ipConnCounter.Check(ip) {
		conn.Close()
		log.Debugf("ip : %s conns reach max")
		return
	}
	connid, clusterConn, err := connectCluster(0, cmd)
	if err != nil {
		conn.Close()
		clusterConn.Close()
		return
	}
	utils.IoBind(conn, &clusterConn, func(err error) {
		conn.Close()
		clusterConn.Close()
		serverListenerConns.Delete(cmd.TunnelID, connid)
		log.Debugf("connection %d - %d released", cmd.TunnelID, connid)
	}, func(bytesCount int, isPositive bool) {}, 0)
	log.Debugf("connection %d - %d created success", cmd.TunnelID, connid)
	serverListenerConns.Put(cmd.TunnelID, connid, &clusterConn)
	return
}
func connectCluster(_connid uint64, cmd utils.MsgServerOpenPort) (connid uint64, clusterConn tls.Conn, err error) {
	tunnleID := cmd.TunnelID
	if _connid == 0 {
		var src = rand.NewSource(time.Now().UnixNano())
		s := fmt.Sprintf("%d", src.Int63())
		str := s[len(s)-5:len(s)-1] +
			fmt.Sprintf("%d", uint64(time.Now().UnixNano()))[7:]
		connid, _ = strconv.ParseUint(str, 10, 64)
	} else {
		connid = _connid
	}
	log.Debugf("new connection %d - %d", tunnleID, connid)
	clusterConn, err = utils.TlsConnect(cfg.GetString("host"), cfg.GetInt("port.conns"), 5000)
	if err != nil {
		log.Warnf("connect to cluster fail ,err:%s", err)
		return
	}
	writer := bufio.NewWriter(&clusterConn)
	pkg := new(bytes.Buffer)
	binary.Write(pkg, binary.LittleEndian, utils.CS_SERVER)
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
