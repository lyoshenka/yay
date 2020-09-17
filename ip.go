package main

import (
	"bytes"
	"net"
	"net/http"
	"strings"
)

type ipRange struct {
	start, end net.IP
}

func (r ipRange) Contains(ip net.IP) bool {
	return bytes.Compare(ip, r.start) >= 0 && bytes.Compare(ip, r.end) < 0
}

var privateRanges = []ipRange{
	{
		start: net.ParseIP("10.0.0.0"),
		end:   net.ParseIP("10.255.255.255"),
	},
	{
		start: net.ParseIP("100.64.0.0"),
		end:   net.ParseIP("100.127.255.255"),
	},
	{
		start: net.ParseIP("127.0.0.1"),
		end:   net.ParseIP("127.255.255.255"),
	},
	{
		start: net.ParseIP("172.16.0.0"),
		end:   net.ParseIP("172.31.255.255"),
	},
	{
		start: net.ParseIP("192.0.0.0"),
		end:   net.ParseIP("192.0.0.255"),
	},
	{
		start: net.ParseIP("192.168.0.0"),
		end:   net.ParseIP("192.168.255.255"),
	},
	{
		start: net.ParseIP("198.18.0.0"),
		end:   net.ParseIP("198.19.255.255"),
	},
}

func isPrivateSubnet(ipAddress net.IP) bool {
	if ipAddress.To4() == nil {
		return false
	}
	for _, r := range privateRanges {
		if r.Contains(ipAddress) {
			return true
		}
	}
	return false
}

func requestIP(r *http.Request) string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || isPrivateSubnet(realIP) {
				// bad address, go to next
				continue
			}
			return ip
		}
	}

	ipParts := strings.Split(r.RemoteAddr, ":")
	ip := strings.Join(ipParts[:len(ipParts)-1], ":")

	if ip == "[::1]" {
		return "127.0.0.1"
	}
	return ip
}
