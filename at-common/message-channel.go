package at_common

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"reflect"
	"sync"
	"time"
)

type Msg struct {
	MsgType int
}
type message struct {
	Type int
	Data interface{}
}

type msgErrorHandler func(channel *MessageChannel, rawMsg interface{}, err error)
type msgCallback func(channel *MessageChannel, msg interface{})
type closeCallback func(channel *MessageChannel, isPeerClose bool)

type MessageChannel struct {
	reader          *bufio.Reader
	writer          *bufio.Writer
	msgHandler      map[int][]msgCallback
	msgErrorHandler msgErrorHandler
	msgTypeMap      map[int]interface{}
	readlock        *sync.Mutex
	writelock       *sync.Mutex
	readHandler     func(msg message)
	serveStared     bool
	ConnectionID    uint64
	Conn            *net.Conn
	remoteAddr      net.Addr
	localAddr       net.Addr
}

func NewMessageChannel(conn *net.Conn) MessageChannel {
	ch := MessageChannel{
		Conn:            conn,
		remoteAddr:      (*conn).RemoteAddr(),
		localAddr:       (*conn).LocalAddr(),
		reader:          bufio.NewReader(*conn),
		writer:          bufio.NewWriter(*conn),
		msgTypeMap:      map[int]interface{}{},
		readlock:        &sync.Mutex{},
		writelock:       &sync.Mutex{},
		ConnectionID:    rand.Uint64(),
		msgErrorHandler: func(channel *MessageChannel, rawMsg interface{}, err error) {},
		msgHandler:      map[int][]msgCallback{},
	}
	ch.RegMsg(MSG_TYPE_PING, new(MsgPing), func(channel *MessageChannel, msg interface{}) {
		msgPing := msg.(*MsgPing)
		channel.Pong(msgPing.ID)
		return
	})
	return ch
}
func NewMessageChannelTls(conn *tls.Conn) MessageChannel {
	con := net.Conn(conn)
	return NewMessageChannel(&con)
}

func (mc *MessageChannel) CloseConn() (err error) {
	(*mc.Conn).SetDeadline(time.Now().Add(time.Millisecond))
	return (*mc.Conn).Close()
}

func (mc *MessageChannel) RegMsg(msgType int, msg interface{}, fn msgCallback, fns ...msgCallback) {
	mc.msgTypeMap[msgType] = msg
	mc.msgHandler[msgType] = append(fns, fn)
}
func (mc *MessageChannel) SetMsgErrorHandler(fn func(channel *MessageChannel, rawMsg interface{}, err error)) {
	mc.msgErrorHandler = fn
}
func (mc *MessageChannel) encode(data interface{}) (msg []byte, err error) {
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
func (mc *MessageChannel) Write(msg interface{}) (err error) {
	defer mc.writelock.Unlock()
	mc.writelock.Lock()
	if reflect.TypeOf(msg).Kind().String() != "struct" {
		err = fmt.Errorf("error : message must be a struct , send to %s", mc.remoteAddr)
		return
	}
	if _, ok := reflect.TypeOf(msg).FieldByName("MsgType"); !ok {
		err = fmt.Errorf("error : message must be has MsgType field , send to %s", mc.remoteAddr)
		return
	}
	if reflect.ValueOf(msg).FieldByName("MsgType").Int() <= 0 {
		err = fmt.Errorf("error : message's  MsgType field must be great than 0, send to %s", mc.remoteAddr)
		return
	}
	v := reflect.ValueOf(msg).FieldByName("MsgType").Int()
	pack := message{
		Type: int(v),
		Data: msg,
	}
	var msgData []byte
	msgData, err = mc.encode(pack)
	if err != nil {
		err = fmt.Errorf("encode message error : %s , to : %s", err, mc.remoteAddr)
		return
	}
	_, err = mc.writer.Write(msgData)

	if err != nil {
		err = fmt.Errorf("write messasge fail to %s , ERR : %s", mc.remoteAddr, err)
		return
	}
	err = mc.writer.Flush()
	if err != nil {
		err = fmt.Errorf("flush messasge fail %s , ERR : %s", mc.remoteAddr, err)
		return
	}
	return
}
func (mc *MessageChannel) ReadTimeout(msg interface{}, timeout int) (err error) {
	if !mc.serveStared {
		err = fmt.Errorf("DoServe() must be called before call Read(),From : %s", mc.remoteAddr)
		return
	}
	type M struct {
		Err error
		Msg interface{}
	}
	msgChn := make(chan M, 1)
	mc.readHandler = func(rawMsg message) {
		_, err = mc.toStruct(rawMsg.Data, &msg)
		msgChn <- M{
			Msg: msg,
			Err: err,
		}
	}
	m := M{}
	if timeout > 0 {
		select {
		case m = <-msgChn:
			msg = m.Msg
			err = m.Err
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			err = fmt.Errorf("read channel message timeout from %s", mc.remoteAddr)
		}
	} else {
		m = <-msgChn
		msg = m.Msg
		err = m.Err
	}
	return
}

func (mc *MessageChannel) Read(msg interface{}) (err error) {
	return mc.ReadTimeout(msg, 0)
}

func (mc *MessageChannel) read() (msg message, err error) {
	defer func() {
		if err != nil {
			(*mc.Conn).Close()
		}
		mc.readlock.Unlock()
	}()
	mc.readlock.Lock()
	// 读取消息的长度
	lengthByte, err := mc.reader.Peek(4)
	if err != nil {
		err = fmt.Errorf("read message length error : %s , from %s", err, mc.remoteAddr)
		return
	}
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err = binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		err = fmt.Errorf("read message error : %s , from %s", err, mc.remoteAddr)
		return
	}
	if int32(mc.reader.Buffered()) < length+4 {
		err = fmt.Errorf("message data length error from %s", mc.remoteAddr)
		return
	}
	// 读取消息真正的内容
	pack := make([]byte, int(4+length))
	_, err = mc.reader.Read(pack)
	if err != nil {
		err = fmt.Errorf("read message error : %s , from %s", err, mc.remoteAddr)
		return
	}
	err = json.Unmarshal(pack[4:], &msg)
	if err != nil {
		err = fmt.Errorf("unmarshal message error : %s , from %s", err, mc.remoteAddr)
		return
	}
	return
}
func (mc *MessageChannel) DoServe(errfn func(err error)) {
	mc.serveStared = true
	go func() {
		var err error
		var msg message
		for {
			msg, err = mc.read()
			if err != nil {
				go errfn(err)
				mc.serveStared = false
				return
			}
			if mc.readHandler != nil {
				go mc.readHandler(msg)
				mc.readHandler = nil
				continue
			}

			h, ok := mc.msgHandler[msg.Type]
			var data interface{}

			if ok {
				data, err = mc.parse(msg)
			} else {
				err = fmt.Errorf("msg handler not found , msgType:%d", msg.Type)
			}
			if err != nil {
				go mc.msgErrorHandler(mc, msg.Data, err)
			} else {
				for i := len(h) - 1; i >= 0; i-- {
					go h[i](mc, data)
				}
			}
		}
	}()
	return
}
func (mc *MessageChannel) Ping() (id string, err error) {
	id = randStr(32)
	ping := MsgPing{
		Msg: Msg{MsgType: MSG_TYPE_PING},
		ID:  id,
	}
	err = mc.Write(ping)
	return
}
func (mc *MessageChannel) Pong(id string) (err error) {
	pong := MsgPong{
		Msg: Msg{MsgType: MSG_TYPE_PONG},
		ID:  id,
	}
	err = mc.Write(pong)
	return
}

func (mc *MessageChannel) RemoteAddr() net.Addr {
	return mc.remoteAddr
}
func (mc *MessageChannel) LocalAddr() net.Addr {
	return mc.localAddr
}

func (mc *MessageChannel) parse(msg message) (data interface{}, err error) {
	data, ok := mc.msgTypeMap[msg.Type]
	if !ok {
		err = fmt.Errorf("message type not registed")
		return
	}
	mbytes, err := json.Marshal(msg.Data)
	if err != nil {
		return
	}
	err = json.Unmarshal(mbytes, data)
	return
}
func (mc *MessageChannel) toStruct(msg, struc interface{}) (data interface{}, err error) {
	mbytes, err := json.Marshal(msg)
	if err != nil {
		return
	}
	err = json.Unmarshal(mbytes, &struc)
	data = struc
	return
}

func randStr(strlen int) string {
	codes := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	codeLen := len(codes)
	data := make([]byte, strlen)
	rand.Seed(time.Now().UnixNano() + rand.Int63() + rand.Int63() + rand.Int63() + rand.Int63())
	for i := 0; i < strlen; i++ {
		idx := rand.Intn(codeLen)
		data[i] = byte(codes[idx])
	}
	return string(data)
}
