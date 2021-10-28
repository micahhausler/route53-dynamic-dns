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
	zoneId := flag.StringP("default-zone", "z", "", "Default Zone ID to update")
	ttl := flag.Int64("default-ttl", 0, "Default time to live on record")
	configFile := flag.String("config", "", "Path to a config JSON file")
	dryRun := flag.Bool("dry-run", false, "Dry run, don't actually update records")
	flag.Parse()

	f, err := os.Open(*configFile)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	configs, err := pkg.ParseConfigFile(f)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}

	ips, err := lib.FindPublicAddrs(*ifaces...)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}
	fmt.Println("Found IPs:")
	for _, ip := range ips {
		fmt.Println(ip.String())
	}

	updater, err := pkg.NewRoute53Updater(*zoneId, pkg.SetDefaultTTL(*ttl))
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}
	for _, config := range configs {
		err = updater.Update(config, ips, *dryRun)
		if err != nil {
			fmt.Printf("Error updating %#v: %+v\n", config, err)
			os.Exit(1)
		}
	}
	fmt.Printf("Updated %d configs successfully\n", len(configs))

}
