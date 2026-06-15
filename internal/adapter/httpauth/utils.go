package httpauth

import (
	"net"
	"net/netip"
	"strings"
)

func parseTrustedProxies(values []string) []netip.Prefix {
	prefixes := make([]netip.Prefix, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if strings.Contains(value, "/") {
			prefix, err := netip.ParsePrefix(value)
			if err == nil {
				prefixes = append(prefixes, prefix)
			}
			continue
		}
		addr, err := netip.ParseAddr(value)
		if err == nil {
			prefixes = append(prefixes, netip.PrefixFrom(addr, addr.BitLen()))
		}
	}
	return prefixes
}

func parseIP(value string) (netip.Addr, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return netip.Addr{}, false
	}
	if addr, err := netip.ParseAddr(value); err == nil {
		return addr, true
	}
	host, _, err := net.SplitHostPort(value)
	if err != nil {
		return netip.Addr{}, false
	}
	addr, err := netip.ParseAddr(host)
	return addr, err == nil
}

func parseXForwardedFor(value string) (netip.Addr, bool) {
	for part := range strings.SplitSeq(value, ",") {
		if addr, ok := parseIP(part); ok {
			return addr, true
		}
	}
	return netip.Addr{}, false
}

func ipInPrefixes(addr netip.Addr, prefixes []netip.Prefix) bool {
	for _, prefix := range prefixes {
		if prefix.Contains(addr) {
			return true
		}
	}
	return false
}
