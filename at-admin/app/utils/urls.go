package utils

import (
	"fmt"
	"net/url"
)

type Urls struct {}

func NewUrls() *Urls  {
	return &Urls{}
}

func (urls *Urls) UrlEncode(urlStr string) string {

	fmt.Println(urlStr)
	urlParse, _ := url.Parse(urlStr)
	urlEncode := urlParse.Query().Encode()
	fmt.Println(urlParse.EscapedPath());
	return "http://" + urlParse.Host + urlParse.EscapedPath() + "?" + urlEncode;
}