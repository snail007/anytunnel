package at_common

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"
)

func TestNotInternal(t *testing.T) {
	for {
		fmt.Println(rand.Uint64() / 10000)
		time.Sleep(time.Second)
	}
}
func TestOpenPort(t *testing.T) {
	//  /port/open/:TunnelID/:ServerToken/:ServerBindIP/:ServerListenPort
	//  /:ClientToken/:ClientLocalHost/:ClientLocalPort/:Protocol/:BytesPerSec
	body, code, err := HttpGet("https://127.0.0.1:37080/port/open/3/guest/0.0.0.0/20090/guest_client/8.8.8.8/53/2/102400")
	if err != nil {
		t.Error(err)
	} else {
		b := string(body)
		fmt.Println(code)
		fmt.Println(b)
	}
}

func TestClosePort(t *testing.T) {
	// /port/close/:TunnelID
	body, code, err := HttpGet("https://127.0.0.1:37080/port/close/3")
	if err != nil {
		t.Error(err)
	} else {
		b := string(body)
		fmt.Println(code)
		fmt.Println(b)
	}
}

func TestStatusPort(t *testing.T) {
	// /port/close/:TunnelID
	body, code, err := HttpGet("https://127.0.0.1:37080/port/status/1")
	if err != nil {
		t.Error(err)
	} else {
		b := string(body)
		fmt.Println(code)
		fmt.Println(b)
	}
}
func TestTrafficCount(t *testing.T) {
	// /traffic/count
	body, code, err := HttpGet("https://127.0.0.1:37080/traffic/count")
	if err != nil {
		t.Error(err)
	} else {
		b := string(body)
		fmt.Println(code)
		fmt.Println(b)
	}
}

func TestTrafficCountTunnel(t *testing.T) {
	// /traffic/count/:TunnelID
	body, code, err := HttpGet("https://127.0.0.1:37080/traffic/count/1")
	if err != nil {
		t.Error(err)
	} else {
		b := string(body)
		fmt.Println(code)
		fmt.Println(b)
	}
}
func TestServerOffline(t *testing.T) {
	// /server/offline/:ServerToken
	body, code, err := HttpGet("https://127.0.0.1:37080/server/offline/guest")
	if err != nil {
		t.Error(err)
	} else {
		b := string(body)
		fmt.Println(code)
		fmt.Println(b)
	}
}
func TestClientOffline(t *testing.T) {
	// /server/offline/:ClientToken
	body, code, err := HttpGet("https://127.0.0.1:37080/client/offline/guest_client")
	if err != nil {
		t.Error(err)
	} else {
		b := string(body)
		fmt.Println(code)
		fmt.Println(b)
	}
}

func TestCsAdd(t *testing.T) {

	values := map[string]string{
		"server_ip":    "127.0.0.1",
		"server_token": "akdhkajhdkashdksdhk",
		"client_ip":    "124.176.89.90",
		"client_token": "akhdkahdah",
		"cluster_ip":   "34.67.78.89",
		"comment":      "djadjllajd",
	}
	body, _, err := HttpPost("https://127.0.0.1:37081/cs/add", values, nil)
	if err != nil {
		t.Error(err)
	} else {
		b := string(body)
		fmt.Println(b)
	}
}

func TestCsUpdate(t *testing.T) {

	values := map[string]string{
		"cs_id":        "2",
		"server_ip":    "127.0.0.1",
		"server_token": "akdhkajhdkashdksdhk",
		"client_ip":    "124.176.89.90",
		"client_token": "tokentoken",
		"cluster_ip":   "34.67.78.89",
		"comment":      "djadjllajd",
	}
	body, _, err := HttpPost("https://127.0.0.1:37081/cs/update", values, nil)
	if err != nil {
		t.Error(err)
	} else {
		b := string(body)
		fmt.Println(b)
	}
}

func TestCsList(t *testing.T) {
	//body, err := HttpGet("https://127.0.0.1:37081/cs/list")
	body, _, err := HttpGet("https://127.0.0.1:37081/cs/list?keyword=tokentoken")
	if err != nil {
		t.Error(err)
	} else {
		b := string(body)
		fmt.Println(b)
	}
}

func TestCsDelete(t *testing.T) {
	body, _, err := HttpGet("https://127.0.0.1:37081/cs/delete?cs_id=1")
	if err != nil {
		t.Error(err)
	} else {
		b := string(body)
		fmt.Println(b)
	}
}

func TestCsGetCsByCsId(t *testing.T) {

	body, _, err := HttpGet("https://127.0.0.1:37081/cs/getCsByCsId?cs_id=1")
	if err != nil {
		t.Error(err)
	} else {
		b := string(body)
		fmt.Println(b)
	}
}

func TestClusterReport(t *testing.T) {

	values := map[string]string{
		"sys_conn_number":    "1312",
		"tunnel_conn_number": "3443",
		"bandwidth":          "1232",
	}
	body, code, err := HttpPost("https://127.0.0.1:37081/cluster/report", values, nil)
	if err != nil {
		t.Error(err)
	} else {
		b := string(body)
		fmt.Println(code)
		fmt.Println(b)
	}
}
func TestCSAuth(t *testing.T) {

	body, code, err := HttpGet("https://127.0.0.1:37081/cs/auth?token=guest&type=server")
	if err != nil {
		t.Error(err)
	} else {
		b := string(body)
		fmt.Println(code)
		fmt.Println(b)
	}
}
func TestCSOnlineOffline(t *testing.T) {

	body, code, err := HttpGet("https://127.0.0.1:37081/cs/status?token=guest_client&type=client&action=offline")
	if err != nil {
		t.Error(err)
	} else {
		b := string(body)
		fmt.Println(code)
		fmt.Println(b)
	}
}

func TestUserTraffic(t *testing.T) {
	data := `{"1":{"positive":32300100,"negative":5324},"2":{"negative":3997696,"positive":0},"3":{"positive":3231,"negative":5324}}`
	body, code, err := HttpPostRaw("https://127.0.0.1:37081/user/traffic", data, nil, true)
	if err != nil {
		t.Error(err)
	} else {
		b := string(body)
		fmt.Println(code)
		fmt.Println(b)
	}
}

func TestCSCluster(t *testing.T) {
	body, code, err := HttpGet("https://127.0.0.1:29531/cluster/get?token=guest&type=server")
	if err != nil {
		t.Error(err)
	} else {
		b := string(body)
		fmt.Println(code)
		fmt.Println(b)
	}
}

type ClusterTunnel struct {
	TunnelID         uint64
	ServerToken      string
	ServerBindIP     string
	ServerListenPort int
	ClientToken      string
	ClientLocalHost  string
	ClientLocalPort  int
	Protocol         int
	BytesPerSec      float64
}

func Test(t *testing.T) {
	_, code, err := HttpGet("https://127.0.0.1:37080/port/open/2/guest/0.0.0.0/20080/guest_client/127.0.0.1/80/1/102400")
	fmt.Println(code, err)
}
func getMyInterfaceAddr() ([]net.IP, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	addresses := []net.IP{}
	for _, iface := range ifaces {

		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			addresses = append(addresses, ip)
		}
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("no address Found, net.InterfaceAddrs: %v", addresses)
	}
	//only need first
	return addresses, nil
}
func TestReadTimeout(t *testing.T) {
	conn, err := TlsConnect("127.0.0.1", 29531, 3000)
	if err != nil {
		t.Error(err)
	} else {
		byt := make([]byte, 1)
		_, err := conn.Read(byt)
		if err != nil {
			t.Error(err)
		}
	}
}
