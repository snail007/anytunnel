package utils

import (
	"net/url"
	"strings"
)

type Urls struct {}

func NewUrls() *Urls  {
	return &Urls{}
}

func (Urls *Urls) HttpQueryBuild(queryValues map[string]string) (queryString string) {
	queryString = ""
	for queryKey, queryValue := range queryValues {
		queryString = queryString + "&" + queryKey + "=" + url.QueryEscape(queryValue)
	}
	queryString = strings.Replace(queryString, "&", "", 1)
	return
}