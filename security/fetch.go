package security

import (
    "errors"
    "net"
    "net/http"
    "net/url"
    "time"
)

var ErrDisallowedURL = errors.New("disallowed host/IP")

func SafeGet(rawURL string) (*http.Response, error) {
    u, err := url.Parse(rawURL)

    if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
        return nil, ErrDisallowedURL
    }

    ips, err := net.LookupIP(u.Hostname())
    if err != nil {
        return nil, err
    }

    for _, ip := range ips {
        if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() {
            return nil, ErrDisallowedURL
        }
    }

    client := &http.Client{Timeout: 10 * time.Second}
    return client.Get(u.String())
}
