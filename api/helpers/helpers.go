package helpers

import (
	"os"
	"strings"
)
func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}
func RemoveDomainError(url string) bool{
	if os.Getenv("DOMAIN") == url {
		return false
	}
	newURL := strings.Replace(url,"http://","",1)
	newURL = strings.Replace(newURL, "https://","",1)
	newURL = strings.Replace(newURL, "www.","",1)
	newURL = strings.Split(newURL,"/")[0]
	if os.Getenv("DOMAIN") == newURL {
		return false
	} 
	return true
}