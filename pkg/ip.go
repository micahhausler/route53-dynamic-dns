package pkg

import (
	"net"

	"github.com/pkg/errors"
	"inet.af/netaddr"
)

func FindPublicAddrs(names ...string) ([]netaddr.IP, error) {
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

func findPublicAddrs(ifaces []net.Interface) ([]netaddr.IP, error) {
	resp := []netaddr.IP{}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		for _, addr := range addrs {
			ipp, err := netaddr.ParseIPPrefix(addr.String())
			if err != nil {
				return nil, errors.WithStack(err)
			}
			ip := ipp.IP()
			if !(ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsMulticast()) {
				resp = append(resp, ip)
			}
		}
	}
	return resp, nil
}
