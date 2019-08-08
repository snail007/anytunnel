package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"

	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	utils "anytunnel/at-common"

	"github.com/julienschmidt/httprouter"
	"github.com/snail007/mini-logger"
	"github.com/valyala/fasthttp"
)

var (
	apiTimeout = time.Second * 30
)

func apiTrafficCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	d, err := trafficCounter.AllDataJSON()
	if err != nil {
		jsonError(w, err, nil)
		return
	}
	fmt.Fprint(w, `{"code":1,"message":"","data":`+string(d)+`}`)
}
func apiTrafficCountTunnel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_TunnelID := ps.ByName("TunnelID")
	TunnelID, err := strconv.ParseUint(_TunnelID, 10, 64)
	if err != nil {
		jsonError(w, err, nil)
		return
	}
	_, err = poolTunnel.Get(TunnelID)
	if err != nil {
		jsonError(w, err.Error(), nil)
		return
	}
	d, err := trafficCounter.TunnelDataJSON(TunnelID)
	if err != nil {
		jsonError(w, err, nil)
		return
	}
	fmt.Fprint(w, `{"code":1,"message":"","data":`+string(d)+`}`)
}
func apiServerStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ServerToken := ps.ByName("ServerToken")
	_, err := poolServerControlChannel.Get(ServerToken)
	if err != nil {
		jsonError(w, err.Error(), nil)
		return
	}
	jsonSuccess(w, "", nil)
}
func apiClientStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ClientToken := ps.ByName("ClientToken")
	_, err := poolClientControlChannel.Get(ClientToken)
	if err != nil {
		jsonError(w, err.Error(), nil)
		return
	}
	jsonSuccess(w, "", nil)
}
func apiServerOffline(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ServerToken := ps.ByName("ServerToken")
	c, err := poolServerControlChannel.Get(ServerToken)
	if err != nil {
		jsonError(w, err.Error(), nil)
		return
	}
	err = (*c.ServerMessageChannel.Conn).Close()
	if err != nil {
		jsonError(w, err, nil)
		return
	}
	//server offline should be delete data file associate it
	dataRecovery.DeleteServer(ServerToken)
	jsonSuccess(w, "", nil)
}
func apiClientOffline(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ClientToken := ps.ByName("ClientToken")
	c, err := poolClientControlChannel.Get(ClientToken)
	if err != nil {
		jsonError(w, err.Error(), nil)
		return
	}
	err = (*c.ClientMessageChannel.Conn).Close()
	if err != nil {
		jsonError(w, err, nil)
		return
	}
	jsonSuccess(w, "", nil)
}
func apiPortStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_TunnelID := ps.ByName("TunnelID")
	TunnelID, err := strconv.ParseUint(_TunnelID, 10, 64)
	if err != nil {
		jsonError(w, err, nil)
		return
	}
	clusterTunnel, err := poolTunnel.Get(TunnelID)
	if err != nil {
		jsonError(w, err.Error(), nil)
		return
	}
	cmd := utils.MsgServerStatusPort{
		Msg:      utils.Msg{MsgType: utils.MSG_SERVER_STATUS_PORT},
		TunnelID: TunnelID,
		Protocol: clusterTunnel.Protocol,
	}
	c, err := poolServerControlChannel.Get(clusterTunnel.ServerToken)
	if err != nil {
		jsonError(w, err.Error(), nil)
		return
	}
	err = c.ServerMessageChannel.Write(cmd)
	if err != nil {
		jsonError(w, err.Error(), nil)
		return
	}
	var resp utils.MsgResponse
	err = c.ServerMessageChannel.ReadTimeout(&resp, 5000)
	if err != nil {
		jsonError(w, err.Error(), nil)
		return
	}
	if !resp.IsSuccess() {
		jsonError(w, resp.Message, nil)
		return
	}
	jsonSuccess(w, resp.Message, nil)
}
func apiPortClose(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_TunnelID := ps.ByName("TunnelID")
	TunnelID, err := strconv.ParseUint(_TunnelID, 10, 64)
	if err != nil {
		jsonError(w, err, nil)
		return
	}
	clusterTunnel, err := poolTunnel.Get(TunnelID)
	if err != nil {
		jsonError(w, err.Error(), nil)
		return
	}
	cmd := utils.MsgServerClosePort{
		Msg:      utils.Msg{MsgType: utils.MSG_SERVER_CLOSE_PORT},
		TunnelID: TunnelID,
		Protocol: clusterTunnel.Protocol,
	}
	c, err := poolServerControlChannel.Get(clusterTunnel.ServerToken)
	if err != nil {
		jsonError(w, err.Error(), nil)
		return
	}
	err = c.ServerMessageChannel.Write(cmd)
	if err != nil {
		jsonError(w, err.Error(), nil)
		return
	}
	var resp utils.MsgResponse
	err = c.ServerMessageChannel.ReadTimeout(&resp, 5000)
	if err != nil {
		jsonError(w, err.Error(), nil)
		return
	}
	if !resp.IsSuccess() {
		jsonError(w, resp.Message, nil)
		return
	}
	trafficCounter.DeleteTunnel(TunnelID)
	poolTunnel.Delete(TunnelID)
	err = dataRecovery.Delete(clusterTunnel.ServerToken, clusterTunnel.TunnelID)
	if err != nil {
		jsonError(w, err.Error(), nil)
		return
	}
	jsonSuccess(w, "success", nil)
}
func apiPortOpen(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_TunnelID := ps.ByName("TunnelID")
	ServerToken := ps.ByName("ServerToken")
	ServerBindIP := ps.ByName("ServerBindIP")
	_ServerListenPort := ps.ByName("ServerListenPort")
	ClientToken := ps.ByName("ClientToken")
	ClientLocalHost := ps.ByName("ClientLocalHost")
	_ClientLocalPort := ps.ByName("ClientLocalPort")
	_Protocol := ps.ByName("Protocol")
	_BytesPerSec := ps.ByName("BytesPerSec")
	TunnelID, err := strconv.ParseUint(_TunnelID, 10, 64)
	if err != nil {
		jsonError(w, err, nil)
		return
	}
	ServerListenPort, err := strconv.Atoi(_ServerListenPort)
	if err != nil {
		jsonError(w, err, nil)
		return
	}
	ClientLocalPort, err := strconv.Atoi(_ClientLocalPort)
	if err != nil {
		jsonError(w, err, nil)
		return
	}
	Protocol, err := strconv.Atoi(_Protocol)
	if err != nil {
		jsonError(w, err, nil)
		return
	}
	if Protocol != 1 && Protocol != 2 {
		jsonError(w, "protocol error", nil)
		return
	}
	BytesPerSec, err := strconv.Atoi(_BytesPerSec)
	if err != nil {
		jsonError(w, err, nil)
		return
	}
	if BytesPerSec < 0 {
		jsonError(w, "BytesPerSec error", nil)
		return
	}
	clusterTunnel, err := poolTunnel.Get(TunnelID)
	if err == nil {
		jsonError(w, "tunnel already opened , please close first", nil)
		return
	}
	clusterTunnel = ClusterTunnel{}
	clusterTunnel.TunnelID = TunnelID
	clusterTunnel.ServerToken = ServerToken
	clusterTunnel.ServerBindIP = ServerBindIP
	clusterTunnel.ServerListenPort = ServerListenPort
	clusterTunnel.ClientToken = ClientToken
	clusterTunnel.ClientLocalHost = ClientLocalHost
	clusterTunnel.ClientLocalPort = ClientLocalPort
	clusterTunnel.Protocol = Protocol
	clusterTunnel.BytesPerSec = float64(BytesPerSec)
	poolTunnel.Set(clusterTunnel)
	cmd := utils.MsgServerOpenPort{
		Msg:      utils.Msg{MsgType: utils.MSG_SERVER_OPEN_PORT},
		TunnelID: TunnelID,
		BindIP:   ServerBindIP,
		BindPort: ServerListenPort,
		Protocol: Protocol,
	}
	c, err := poolServerControlChannel.Get(ServerToken)
	if err != nil {
		jsonError(w, err.Error(), nil)
		return
	}
	err = c.ServerMessageChannel.Write(cmd)
	if err != nil {
		jsonError(w, err.Error(), nil)
		return
	}
	var resp utils.MsgResponse
	err = c.ServerMessageChannel.ReadTimeout(&resp, 5000)
	if err != nil {
		jsonError(w, err.Error(), nil)
		return
	}
	if !resp.IsSuccess() {
		jsonError(w, resp.Message, nil)
		return
	}
	trafficCounter.InitTunnel(TunnelID)
	serverTunnelPool.AddTunnel(ServerToken, TunnelID)
	err = dataRecovery.Store(clusterTunnel)
	if err != nil {
		jsonError(w, resp.Message, nil)
		return
	}
	jsonSuccess(w, "success", nil)
	//log.Infof("open port cmd:%s", cmd)

}
func jsonSuccess(w http.ResponseWriter, message, data interface{}) {
	jsonEcho(w, 1, message, data)
}
func jsonError(w http.ResponseWriter, message, data interface{}) {
	jsonEcho(w, 0, message, data)
}
func jsonEcho(w http.ResponseWriter, code int, message, data interface{}) {
	type JSONObj struct {
		Code    int         `json:"code"`
		Message interface{} `json:"message"`
		Data    interface{} `json:"data"`
	}
	jsonObj := JSONObj{
		Code:    code,
		Message: message,
		Data:    data,
	}
	d, err := json.Marshal(jsonObj)
	if err != nil {
		jsonObj.Code = 0
		jsonObj.Message = err.Error()
		d, _ = json.Marshal(jsonObj)
		fmt.Fprint(w, string(d))
		return
	}
	fmt.Fprint(w, string(d))
}
func initHttp(fn func(err error)) {
	router := httprouter.New()
	router.GET("/port/open/:TunnelID/:ServerToken/:ServerBindIP/:ServerListenPort/:ClientToken/:ClientLocalHost/:ClientLocalPort/:Protocol/:BytesPerSec", timeoutFactory(apiPortOpen))
	router.GET("/port/close/:TunnelID", timeoutFactory(apiPortClose))
	router.GET("/port/status/:TunnelID", timeoutFactory(apiPortStatus))
	router.GET("/server/offline/:ServerToken", timeoutFactory(apiServerOffline))
	router.GET("/server/status/:ServerToken", timeoutFactory(apiServerStatus))
	router.GET("/client/offline/:ClientToken", timeoutFactory(apiClientOffline))
	router.GET("/client/status/:ClientToken", timeoutFactory(apiClientStatus))
	router.GET("/traffic/count", timeoutFactory(apiTrafficCount))
	router.GET("/traffic/count/:TunnelID", timeoutFactory(apiTrafficCountTunnel))
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(utils.GetRootCert())

	for _, ip := range cfg.GetStringSlice("port.ip-api") {
		go func(ip string) {
			s := &http.Server{
				Addr:    fmt.Sprintf("%s:%d", ip, cfg.GetInt("port.api")),
				Handler: router,
				TLSConfig: &tls.Config{
					ClientCAs:  pool,
					ClientAuth: tls.RequireAndVerifyClientCert,
					ServerName: "anytunnel-client",
				},
				ReadTimeout: time.Millisecond * 3000,
			}
			log.Infof("listening on [%s] for http api", (*s).Addr)
			if err := ListenAndServeTLS(s, utils.GetServerCert(), utils.GetServerKey()); err != nil {
				log.Errorf("ListenAndServeTLS err:%s", err)
			}
		}(ip)
	}
}
func timeoutFactory(fn func(w http.ResponseWriter, r *http.Request, ps httprouter.Params)) (handle func(w http.ResponseWriter, r *http.Request, ps httprouter.Params)) {

	handle = func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		chn := make(chan bool, 1)
		go func() {
			fn(w, r, ps)
			chn <- true
		}()
		select {
		case <-chn:
		case <-time.After(apiTimeout):
			fmt.Fprint(w, "timeout")
		}
	}
	return
}
func access(ctx *fasthttp.RequestCtx) {
	post := ""
	if cfg.GetBool("log.post") {
		post = string(ctx.Request.Body())
	}
	fields := logger.Fields{
		"code":       strconv.Itoa(ctx.Response.StatusCode()),
		"uri":        string(ctx.RequestURI()),
		"remoteAddr": strings.Split(ctx.RemoteAddr().String(), ":")[0],
		"method":     string(ctx.Method()),
		"host":       string(ctx.Request.Host()),
		"referer":    string(ctx.Request.Header.Referer()),
		"userAgent":  string(ctx.Request.Header.UserAgent()),
		"response":   string(ctx.Response.Body()),
		"post":       post,
	}
	accessLog.With(fields).Info("")
}
func ListenAndServeTLS(srv *http.Server, certPEMBlock, keyPEMBlock []byte) error {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}
	config := &tls.Config{}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, config)
	return srv.Serve(tlsListener)
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetReadDeadline(time.Now().Add(time.Millisecond * 300))
	return tc, nil
}
