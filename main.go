package main

import (
	"fmt"
	"os"

	"github.com/micahhausler/route53-dynamic-dns/pkg"
	lib "github.com/micahhausler/route53-dynamic-dns/pkg"
	flag "github.com/spf13/pflag"
)

func main() {
	ifaces := flag.StringSliceP("interfaces", "i", []string{"eth0"}, "List of interfaces to evaluate")
	zoneId := flag.StringP("zone", "z", "", "Zone ID to update")
	record := flag.StringP("record", "r", "", "Record to update")
	ttl := flag.Int("ttl", 3600, "Time to live on record")
	flag.Parse()
	ips, err := lib.FindPublicAddrs(*ifaces...)

	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}
	fmt.Println("Found IPs:")
	for _, ip := range ips {
		fmt.Println(ip.String())
	}

	updater, err := pkg.NewUpdater(*zoneId, *ttl)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}
	err = updater.Update(*record, ips)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}
}
