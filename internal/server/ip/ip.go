package ip

import (
	"fmt"
	"godrop/internal/common/colored"

	"log/slog"
	"net"
	"os"
)

func InformAboutIp() {
	ip, err := getLocalIPs()
	if err != nil {
		slog.Error("Error while getting local address", "error", err)
		os.Exit(1)
	}
	colored.PrintColoredWithTags("Current local ip addresses <RED>%v<RED> share it with your friend to download file\n", ip)

	// slog.Info("Current local ip address", "ip", ip)
}

func getLocalIPs() ([]net.IP, error) {
	var ips []net.IP
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addresses {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP)
			}
		}
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf(colored.BuildColoredString("<RED>No addreses found<RED>"))
	}
	return ips, nil
}
