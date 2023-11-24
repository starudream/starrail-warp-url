package main

import (
	"net"
)

func localIP() string {
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if nip, ok := addr.(*net.IPNet); ok {
			if ipv4 := nip.IP.To4(); ipv4 != nil && ipv4[0] != 127 {
				return nip.IP.String()
			}
		}
	}
	return ""
}
