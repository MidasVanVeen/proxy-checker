package proxyutils

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"h12.io/socks"
)

type Protocol string

const (
	HTTP   Protocol = "http"
	HTTPS  Protocol = "https"
	SOCKS4 Protocol = "socks4"
	SOCKS5 Protocol = "socks5"
)

func ConfigureTransport(protocol Protocol, proxy string) (*http.Transport, error) {
	switch protocol {
	case HTTP, HTTPS:
		proxyUrl, err := url.Parse(string(protocol) + "://" + proxy)
		if err != nil {
			return nil, err
		}
		return &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}, nil
	case SOCKS4, SOCKS5:
		dialSocksProxy := socks.Dial(string(protocol) + "://" + proxy)
		return &http.Transport{
			Dial: dialSocksProxy,
		}, nil
	default:
		return nil, nil
	}
}

// TODO: replace all this with regex
func ValidateProxy(proxy string) bool {
	proxyParts := strings.Split(proxy, ":")
	if len(proxyParts) != 2 {
		return false
	}
	ip := proxyParts[0]
	port := proxyParts[1]
	if !ValidateIp(ip) {
		return false
	}
	if !ValidatePort(port) {
		return false
	}
	return true
}

func GetIp(proxy string) string {
	return strings.Split(proxy, ":")[0]
}

func ValidateIp(ip string) bool {
	ipParts := strings.Split(ip, ".")
	if len(ipParts) != 4 {
		return false
	}
	for _, part := range ipParts {
		if !ValidateIpPart(part) {
			return false
		}
	}
	return true
}

func ValidateIpPart(part string) bool {
	partInt, err := strconv.Atoi(part)
	if err != nil {
		return false
	}
	if partInt < 0 || partInt > 255 {
		return false
	}
	return true
}

func ValidatePort(port string) bool {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	if portInt < 0 || portInt > 65535 {
		return false
	}
	return true
}
