package at_common

import (
	"crypto/tls"
	"fmt"
)

type Author struct {
	cstype      int
	token       string
	clusterHost string
	clusterPort int
	Conn        *tls.Conn
	Channel     MessageChannel
}

func NewAuthor(token, clusterHost string, clusterPort, cstype int) (author Author, err error) {
	conn, err := TlsConnect(clusterHost, clusterPort, 5000)
	if err != nil {
		return
	}
	author = Author{
		cstype:      cstype,
		token:       token,
		Conn:        &conn,
		clusterHost: clusterHost,
		clusterPort: clusterPort,
		Channel:     NewMessageChannelTls(&conn),
	}
	return
}
func (a *Author) DoControlAuth() (err error) {
	err = a.Channel.Write(MsgLogin{
		Msg:    Msg{MsgType: MSG_TYPE_LOGIN},
		Token:  a.token,
		CSType: a.cstype,
	})
	if err != nil {
		return
	}
	var resp MsgResponse
	err = a.Channel.ReadTimeout(&resp, 30000)
	if err != nil {
		return
	}
	if !resp.IsSuccess() {
		err = fmt.Errorf(resp.Message)
	}
	return
}

func (a *Author) IsServer() bool {
	return a.cstype == CSTYPE_SERVER
}
func (a *Author) IsClient() bool {
	return a.cstype == CSTYPE_CLIENT
}
