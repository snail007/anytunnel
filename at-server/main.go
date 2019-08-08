package main

import (
	utils "anytunnel/at-common"
	"fmt"
	"net"
	"time"

	logger "github.com/snail007/mini-logger"
)

const SERVER_CONN_IDLE_MINUTES = 24 * 60

var (
	control             utils.Author
	err                 error
	serverListeners     = NewListenerMap()
	serverListenerConns = NewServerUserConnPool()
	ipConnCounter       IPConnCounter
)

func init() {
	poster()
	initConfig()
	initLog()
	ipConnCounter = NewIPConnCounter(cfg.GetInt("max"), 3600)
}
func main() {
	defer func() {
		e := recover()
		if e != nil {
			log.Errorf("Exit ERR:%s", e)
		}
		logger.Flush()
	}()
	control, err = utils.NewAuthor(cfg.GetString("token"), cfg.GetString("host"), cfg.GetInt("port.control"), utils.CSTYPE_SERVER)
	if err != nil {
		log.Debugf("create author fail : %s", err)
		return
	}

	control.Channel.SetMsgErrorHandler(func(channel *utils.MessageChannel, msg interface{}, err error) {
		log.Warnf("message pre-process error , ERRR:%s", err)
	})
	control.Channel.RegMsg(utils.MSG_SERVER_STATUS_PORT, new(utils.MsgServerStatusPort), func(channel *utils.MessageChannel, msg interface{}) {
		msgServerStatusPort := msg.(*utils.MsgServerStatusPort)
		log.Debugf("MSG_SERVER_STATUS_PORT revecived , TunnleID:%d", msgServerStatusPort.TunnelID)
		var network, addr string
		var ok bool
		if msgServerStatusPort.Protocol == utils.TUNNEL_PROTOCOL_TCP {
			var listener *net.Listener
			listener, ok = serverListeners.Get(msgServerStatusPort.TunnelID)
			network = (*listener).Addr().Network()
			addr = (*listener).Addr().String()
		} else {
			var listener *net.UDPConn
			listener, ok = serverListeners.GetUDP(msgServerStatusPort.TunnelID)
			network = (*listener).LocalAddr().Network()
			addr = (*listener).LocalAddr().String()
		}
		rmsg := ""
		if !ok {
			rmsg = "listener not exists"
			utils.Response(channel, false, rmsg)
			log.Debugf("Tunnle %d status fail ,err:%s", msgServerStatusPort.TunnelID, rmsg)
			return
		}
		utils.Response(channel, true, fmt.Sprintf("%s:%s", network, addr))
		log.Debugf("Tunnle %d status success", msgServerStatusPort.TunnelID)
		return
	})

	control.Channel.RegMsg(utils.MSG_SERVER_CLOSE_PORT, new(utils.MsgServerClosePort), func(channel *utils.MessageChannel, msg interface{}) {
		msgServerClosePort := msg.(*utils.MsgServerClosePort)
		log.Debugf("MSG_SERVER_CLOSE_PORT revecived , TunnleID:%d", msgServerClosePort.TunnelID)
		var ok bool
		var err error
		if msgServerClosePort.Protocol == utils.TUNNEL_PROTOCOL_TCP {
			var listener *net.Listener
			listener, ok = serverListeners.Get(msgServerClosePort.TunnelID)
			if ok {
				err = (*listener).Close()
			}
		} else {
			var listener *net.UDPConn
			listener, ok = serverListeners.GetUDP(msgServerClosePort.TunnelID)
			if ok {
				err = (*listener).Close()
			}
		}
		rmsg := ""
		if !ok {
			rmsg = "listener not exists"
			utils.Response(channel, false, rmsg)
			log.Debugf("Tunnle %d closed fail ,err:%s", msgServerClosePort.TunnelID, rmsg)
			return
		}
		if err != nil {
			utils.Response(channel, false, err.Error())
			log.Debugf("Tunnle %d closed fail ,err:%s", msgServerClosePort.TunnelID, err)
			return
		}
		utils.Response(channel, true, "")
		log.Debugf("Tunnle %d closed success", msgServerClosePort.TunnelID)
		return
	})
	control.Channel.RegMsg(utils.MSG_SERVER_OPEN_PORT, new(utils.MsgServerOpenPort), func(channel *utils.MessageChannel, msg interface{}) {
		msgServerOpenPort := msg.(*utils.MsgServerOpenPort)
		//log.Debugf("MSG_SERVER_OPEN_PORT revecived [%s:%d], tunnleID : %d", msgServerOpenPort.BindIP, msgServerOpenPort.BindPort, msgServerOpenPort.TunnelID)
		var ok bool
		if msgServerOpenPort.Protocol == utils.TUNNEL_PROTOCOL_TCP {
			_, ok = serverListeners.Get(msgServerOpenPort.TunnelID)
		} else {
			_, ok = serverListeners.GetUDP(msgServerOpenPort.TunnelID)
		}
		if ok {
			serverListeners.Delete(msgServerOpenPort.TunnelID)
			time.Sleep(time.Millisecond * 300)
		}
		var sc utils.ServerChannel
		var err error
		var addr net.Addr
		if msgServerOpenPort.Protocol == utils.TUNNEL_PROTOCOL_TCP {
			sc, err = openPort(*msgServerOpenPort)
			addr = (*sc.Listener).Addr()
		} else {
			sc, err = openUDPPort(*msgServerOpenPort)
			addr = (*sc.UDPListener).LocalAddr()
		}
		if err != nil {
			log.Warnf("open %s port fail ,ERR:%s", msgServerOpenPort.ProtocolString(), err)
			utils.Response(channel, false, err.Error())
			return
		}
		utils.Response(channel, true, "")
		log.Infof("open %s port success , info %s", msgServerOpenPort.ProtocolString(), addr)
		return
	})
	control.Channel.DoServe(func(err error) {
		//log.Fatalf("offline , disconnected from cluster , %s", err)
		log.Fatalf("offline , disconnected from cluster")
	})
	err = control.DoControlAuth()
	if err != nil {
		log.Fatalf("login fail : %s", err)
		return
	}
	//log.Infof("login success %s - %s", control.Channel.LocalAddr(), control.Channel.RemoteAddr())
	log.Infof("login success")
	select {}
}
