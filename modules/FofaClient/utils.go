package FofaClient

import (
	"net"
	"strings"
)

// SplitDomainsAndIPs 解析输入数组并分为 domain 和 ip 两个数组
func SplitDomainsAndIPs(input []string) (domains []string, ips []string) {
	for _, item := range input {
		// Remove URL scheme if present
		cleanItem := strings.TrimPrefix(item, "http://")
		cleanItem = strings.TrimPrefix(cleanItem, "https://")

		// Extract only the protocol + host (removing any directory path)
		if idx := strings.Index(cleanItem, "/"); idx != -1 {
			cleanItem = cleanItem[:idx]
		}

		// Split host and port if present
		host, _, err := net.SplitHostPort(cleanItem)
		if err != nil {
			host = cleanItem // No port present
		}

		// Check if host is an IP address or a domain
		if net.ParseIP(host) != nil {
			ips = append(ips, "\""+cleanItem+"\"")
		} else {
			domains = append(domains, "host=\""+cleanItem+"\"")
		}
	}
	return
}
