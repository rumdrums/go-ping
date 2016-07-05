package main

import (
	"fmt"
	"flag"
	"net"
)

type ip_addr net.IPAddr

func (ip *ip_addr) String() string {
	return fmt.Sprint(*ip)
}

func (ip *ip_addr) Set(value string) error {
	var err error
	ip.IP = net.ParseIP(value)
	if ip.IP == nil {
		return err
	}
	return nil
}
func main() {
	//var ip = flag.Int("flagname", 1234, "help message for flagname")
	var ip ip_addr
	flag.Var(&ip, "ip", "ip address, must be parseable by net.ParseIP")
	flag.Parse()
	fmt.Println(ip)
}

