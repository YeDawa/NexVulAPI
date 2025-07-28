package utils

import (
	"net"
	"net/url"
)

func GetIPFromURL(targetURL string) string {
	u, err := url.Parse(targetURL)
	if err != nil {
		return ""
	}

	host := u.Hostname()
	ips, err := net.LookupIP(host)
	if err != nil || len(ips) == 0 {
		return ""
	}

	return ips[0].String()
}
