package main

import (
	des "anytunnel/at-common/des"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/testutils"
)

func main() {
	key := flag.String("k", "abcd1234", "encrypt or decrypt key")
	isEncrypt := flag.Bool("e", false, "encrypt text")
	isDecrypt := flag.Bool("d", false, "decrypt text")
	text := flag.String("t", "", "text to decrypt or encrypt")
	ip := flag.String("i", "", "ip to bind")
	port := flag.Int("p", 58080, "http port to listen")
	flag.Parse()
	if *isEncrypt {
		res, err := des.Encrypt([]byte(*text), *key)
		if err != nil {
			fmt.Printf("ERR:%s\n", err)
			return
		}
		fmt.Printf("%s\n", res)
		return
	}
	if *isDecrypt {
		res, err := des.Decrypt(*text, *key)
		if err != nil {
			fmt.Printf("ERR:%s\n", err)
			return
		}
		fmt.Printf("%s\n", res)
		return
	}
	logrus.SetLevel(logrus.FatalLevel)
	fwd, _ := forward.New()
	redirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_hostarr := strings.Split(req.Host, ".")
		if len(_hostarr) < 3 {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden"))
			return
		}
		hostAndPort, err := des.Decrypt(_hostarr[0], *key)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden"))
			return
		}
		hostAndPortArr := strings.Split(string(hostAndPort), ":")
		if len(hostAndPortArr) != 2 {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden"))
			return
		}
		host := hostAndPortArr[0]
		if host == "" || host == "0.0.0.0" {
			host = "127.0.0.1"
		}
		req.URL = testutils.ParseURI(fmt.Sprintf("http://%s:%s/", host, hostAndPortArr[1]))
		fwd.ServeHTTP(w, req)

	})
	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", *ip, *port),
		Handler: redirect,
	}
	log.Printf("http proxy on %s", s.Addr)
	log.Fatalf("ERR:%s", s.ListenAndServe())
}
