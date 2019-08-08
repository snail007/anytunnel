package at_common

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func HttpGetNotInternal(URL string) (body []byte, code int, err error) {
	return _HttpGet(URL, false)
}
func HttpPostNotInternal(URL string, data map[string]string, header map[string]string) (body []byte, code int, err error) {
	return _HttpPost(URL, data, header, false)
}
func HttpPost(URL string, data map[string]string, header map[string]string) (body []byte, code int, err error) {
	return _HttpPost(URL, data, header, true)
}
func HttpGet(URL string) (body []byte, code int, err error) {
	return _HttpGet(URL, true)
}
func _HttpGet(URL string, isInternal bool) (body []byte, code int, err error) {
	client, tr, err := getRequestClient(URL, isInternal)
	if err != nil {
		return
	}
	resp, err := client.Get(URL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	defer tr.CloseIdleConnections()
	code = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)
	return
}
func _HttpPost(URL string, data map[string]string, header map[string]string, isInternal bool) (body []byte, code int, err error) {
	postParamsString := ""
	if data != nil {
		postParams := []string{}
		for k, v := range data {
			postParams = append(postParams, url.QueryEscape(k)+"="+url.QueryEscape(v))
		}
		postParamsString = strings.Join(postParams, "&")
	}
	return HttpPostRaw(URL, postParamsString, header, isInternal)
}

func HttpPostRaw(URL, postParamsString string, header map[string]string, isInternal bool) (body []byte, code int, err error) {
	var resp *http.Response
	var client *http.Client
	var tr *http.Transport
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
		if tr != nil {
			tr.CloseIdleConnections()
		}
	}()
	req, err := http.NewRequest("POST", URL, strings.NewReader(postParamsString))
	if err != nil {
		return
	}
	client, tr, err = getRequestClient(URL, isInternal)
	if err != nil {
		return
	}
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	code = resp.StatusCode
	if err != nil {
		return
	}
	return
}

func getRequestClient(url string, isInternal bool) (client *http.Client, tr *http.Transport, err error) {
	var conf *tls.Config
	conf, err = getRequestTlsConfig(isInternal)
	if err != nil {
		return
	}
	if strings.Contains(url, "https://") {
		tr = &http.Transport{TLSClientConfig: conf}
		client = &http.Client{Timeout: time.Second * 5, Transport: tr}
	} else {
		tr = &http.Transport{}
		client = &http.Client{Timeout: time.Second * 5, Transport: tr}
	}
	return
}

func UrlArgs(url, args string) string {
	if strings.Contains(url, "?") {
		return url + "&" + args
	}
	return url + "?" + args
}
