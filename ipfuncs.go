package main

import (
	"net"
)

func ExpandRoute(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)

	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	return ips[1 : len(ips)-1], nil
}

func RouteCount(cidr string) int {
	count := 0
	ip, ipnet, err := net.ParseCIDR(cidr)

	if err != nil {
		return count
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		count++
	}

	return count - 1
}

func GetTotalIpCountByIpv4Prefixes(cidrs []Prefix) int {

	count := 0

	for i := 0; i < len(cidrs); i++ {
		count += RouteCount(cidrs[i].Name)
	}

	return count
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func ExpandRage(start string, end string) []string {

	rangeBegin := net.ParseIP(start)
	rangeEnd := net.ParseIP(end)

	ip := dupIP(rangeBegin)

	out := []string{ip.String()}

	for !ip.Equal(rangeEnd) {
		ip = nextIP(ip)
		out = append(out, ip.String())
	}

	return out
}

func dupIP(ip net.IP) net.IP {
	// To save space, try and only use 4 bytes
	if x := ip.To4(); x != nil {
		ip = x
	}
	dup := make(net.IP, len(ip))
	copy(dup, ip)
	return dup
}

func nextIP(ip net.IP) net.IP {
	next := dupIP(ip)
	for j := len(next) - 1; j >= 0; j-- {
		next[j]++
		if next[j] > 0 {
			break
		}
	}
	return next
}

func isCidrIpV4(ipaddr string) bool {
	ip, _, _ := net.ParseCIDR(ipaddr)

	return ip.To4() != nil
}
