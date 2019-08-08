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
	go func() {
		atClient := NewClient()
		atServer := NewServer()
		tunnel := NewTunnel()
		cluster := NewCluster()
		traffic := NewTraffic()
		user := NewUser()
		online := NewOnline()
		clusterApi := NewClusterAPI()
		//intra api router
		routerIntra := httprouter.New()
		//client
		routerIntra.POST("/client/add", timeoutFactory(atClient.Add))
		routerIntra.POST("/client/update", timeoutFactory(atClient.Update))
		routerIntra.GET("/client/delete", timeoutFactory(atClient.Delete))
		routerIntra.GET("/client/list", timeoutFactory(atClient.List))
		routerIntra.GET("/client/getClientById", timeoutFactory(atClient.GetClientByClientId))
		//server
		routerIntra.POST("/server/add", timeoutFactory(atServer.Add))
		routerIntra.POST("/server/update", timeoutFactory(atServer.Update))
		routerIntra.GET("/server/delete", timeoutFactory(atServer.Delete))
		routerIntra.GET("/server/list", timeoutFactory(atServer.List))
		routerIntra.GET("/server/getServerById", timeoutFactory(atServer.GetServerByServerId))
		//tunnel
		routerIntra.POST("/tunnel/add", timeoutFactory(tunnel.Add))
		routerIntra.POST("/tunnel/update", timeoutFactory(tunnel.Update))
		routerIntra.GET("/tunnel/delete", timeoutFactory(tunnel.Delete))
		routerIntra.GET("/tunnel/list", timeoutFactory(tunnel.List))
		routerIntra.GET("/tunnel/getTunnelByTunnelId", timeoutFactory(tunnel.GetTunnelByTunnelId))
		routerIntra.GET("/tunnel/getTunnelsByClientId", timeoutFactory(tunnel.GetTunnelsByClientId))
		routerIntra.GET("/tunnel/getTunnelsByServerId", timeoutFactory(tunnel.GetTunnelsByServerId))
		routerIntra.GET("/tunnel/open", timeoutFactory(tunnel.TunnelOpen))
		routerIntra.GET("/tunnel/close", timeoutFactory(tunnel.TunnelClose))
		//user
		routerIntra.GET("/user/getUserByName", timeoutFactory(user.GetUserByName))
		routerIntra.GET("/user/getUserByEmail", timeoutFactory(user.GetUserByEmail))
		routerIntra.GET("/user/getUserById", timeoutFactory(user.GetUserByUserId))
		routerIntra.POST("/user/create", timeoutFactory(user.Create))
		routerIntra.POST("/user/update", timeoutFactory(user.Update))
		routerIntra.GET("/user/getRolesByUserId", timeoutFactory(user.GetRolesByUserId))
		routerIntra.GET("/user/cluster", timeoutFactory(user.Cluster))
		//online
		routerIntra.GET("/online/getServerByClusterId", timeoutFactory(online.GetServerByClusterId))
		routerIntra.GET("/online/getClientByClusterId", timeoutFactory(online.GetClientByClusterId))

		//###########cluster api############
		//cs登录认证
		routerIntra.GET("/cs/auth", timeoutFactory(clusterApi.Auth))
		//cs上线下线上报
		routerIntra.GET("/cs/status", timeoutFactory(clusterApi.Status))
		//流量上报
		routerIntra.POST("/user/traffic", timeoutFactory(traffic.Traffic))
		//cluster机器上报
		routerIntra.POST("/cluster/report", timeoutFactory(cluster.Report))

		//#######外网api######
		//extra api router
		routerExtra := httprouter.New()
		routerExtra.GET("/cluster/get", timeoutFactory(clusterApi.Cluster))

		//intra api
		for _, ip := range cfg.GetStringSlice("port.ip-intra") {
			go func(ip string) {
				pool := x509.NewCertPool()
				pool.AppendCertsFromPEM(utils.GetRootCert())
				s := &http.Server{
					Addr:    fmt.Sprintf("%s:%d", ip, cfg.GetInt("port.intra")),
					Handler: routerIntra,
					TLSConfig: &tls.Config{
						ClientCAs:  pool,
						ClientAuth: tls.RequireAndVerifyClientCert,
						ServerName: "anytunnel-client",
					},
				}
				log.Infof("listening on [%s] for http intra api", (*s).Addr)
				if err := ListenAndServeTLS(s, utils.GetServerCert(), utils.GetServerKey()); err != nil {
					log.Errorf("ListenAndServeTLS err:%s", err)
				}
			}(ip)
		}
		//extra api
		for _, ip := range cfg.GetStringSlice("port.ip-extra") {
			go func(ip string) {
				pool := x509.NewCertPool()
				pool.AppendCertsFromPEM(utils.GetRootCert())
				s := &http.Server{
					Addr:    fmt.Sprintf("%s:%d", ip, cfg.GetInt("port.extra")),
					Handler: routerExtra,
					TLSConfig: &tls.Config{
						ClientCAs:  pool,
						ClientAuth: tls.RequireAndVerifyClientCert,
						ServerName: "anytunnel-client",
					},
					ReadTimeout: time.Millisecond * 3000,
				}
				log.Infof("listening on [%s] for http extra api", (*s).Addr)
				if err := ListenAndServeTLS(s, utils.GetServerCert(), utils.GetServerKey()); err != nil {
					log.Errorf("ListenAndServeTLS err:%s", err)
				}
			}(ip)
		}
	}()
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
