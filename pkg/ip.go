package pkg

import (
	"net"
	"net/netip"

	"github.com/pkg/errors"
)

func FindPublicAddrs(names ...string) ([]netip.Addr, error) {
	var (
		ifaces []net.Interface
		err    error
	)
	if len(names) == 0 {
		ifaces, err = net.Interfaces()
		if err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		for _, n := range names {
			i, err := net.InterfaceByName(n)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			ifaces = append(ifaces, *i)
		}
	}
	return findPublicAddrs(ifaces)
}

func findPublicAddrs(ifaces []net.Interface) ([]netip.Addr, error) {
	resp := []netip.Addr{}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		for _, addr := range addrs {
			prefix, err := netip.ParsePrefix(addr.String())
			if err != nil {
				return nil, errors.WithStack(err)
			}

			addr := prefix.Addr()
			if !(addr.IsPrivate() || addr.IsLoopback() || addr.IsLinkLocalUnicast() || addr.IsLinkLocalMulticast() || addr.IsMulticast()) {
				resp = append(resp, addr)
			}
		}
	}
	return resp, nil
}
