package at_common

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
)

const (
	nONE = iota
	TUNNEL_PROTOCOL_TCP
	TUNNEL_PROTOCOL_UDP
	CSTYPE_SERVER
	CSTYPE_CLIENT
	MSG_CLIENT_OPEN_CONNECTION
	MSG_SERVER_OPEN_PORT
	MSG_SERVER_CLOSE_PORT
	MSG_SERVER_STATUS_PORT
	MSG_RESPONSE
	STATUS_SUCCESS
	STATUS_FAIL
	MSG_TYPE_LOGIN
	MSG_TYPE_PING
	MSG_TYPE_PONG
	CS_CLIENT uint8 = 1
	CS_SERVER uint8 = 2
)

type MsgPing struct {
	Msg
	ID string
}
type MsgPong struct {
	Msg
	ID string
}

type MsgClientOpenConnection struct {
	Msg
	TunnelID    uint64
	ConnectinID uint64
	LocalHost   string
	LocalPort   int
	Protocol    int
}

type MsgServerOpenPort struct {
	Msg
	TunnelID uint64
	BindPort int
	BindIP   string
	Protocol int
}

func (m *MsgServerOpenPort) ProtocolString() string {
	return GetProtocolString(m.Protocol)
}

type MsgServerClosePort struct {
	Msg
	TunnelID uint64
	Protocol int
}
type MsgServerStatusPort struct {
	Msg
	TunnelID uint64
	Protocol int
}
type MsgLogin struct {
	Msg
	Token  string
	CSType int
}

func (m *MsgLogin) IsServer() bool {
	return m.CSType == CSTYPE_SERVER
}
func (m *MsgLogin) IsClient() bool {
	return m.CSType == CSTYPE_CLIENT
}
func (m *MsgLogin) CSTypeString() string {
	return GetCSTypeString(m.CSType)
}

type MsgResponse struct {
	Msg
	Status  int
	Message string
}

func (m *MsgResponse) IsSuccess() bool {
	return m.Status == STATUS_SUCCESS
}

func Encode(data interface{}) (msg []byte, err error) {
	var message []byte
	message, err = json.Marshal(data)
	if err != nil {
		return
	}
	// 读取消息的长度
	var length int32 = int32(len(message))
	var pkg *bytes.Buffer = new(bytes.Buffer)
	// 写入消息头
	err = binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return
	}
	// 写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, message)
	if err != nil {
		return
	}
	return pkg.Bytes(), nil
}

func Decode(reader *bufio.Reader) (msg message, err error) {
	// 读取消息的长度
	lengthByte, _ := reader.Peek(4)
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err = binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return
	}
	if int32(reader.Buffered()) < length+4 {
		err = errors.New("data length error")
		return
	}
	// 读取消息真正的内容
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return
	}
	var _msg message
	err = json.Unmarshal(pack[4:], &_msg)
	if err != nil {
		return
	}
	return _msg, nil
}
func GetCSTypeString(cstype int) string {
	switch cstype {
	case CSTYPE_CLIENT:
		return "CLIENT"
	case CSTYPE_SERVER:
		return "SERVER"
	}
	return "UNKONWN_CSTYPE"
}
func GetProtocolString(protocol int) string {
	switch protocol {
	case TUNNEL_PROTOCOL_UDP:
		return "udp"
	case TUNNEL_PROTOCOL_TCP:
		return "tcp"
	}
	return "UNKONWN_PROTOCOL"
}
