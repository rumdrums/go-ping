package main

import(
	"log"
	"os"
	"fmt"
	"net"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/icmp"
)

func main() {
	var address string
	if len(os.Args) > 1 {
		address = os.Args[1]
	} else {
		fmt.Println("Pass an address")
		os.Exit(69)
	}

	// c, err := icmp.ListenPacket("udp4", "localhost")
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	ping := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff,
			Seq: 1,
			Data: []byte("Hello"),
		},
	}
	// 
	marshalled, err := ping.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}
	new_address := net.ParseIP(address)
	_, err = c.WriteTo(marshalled, &net.IPAddr{IP: new_address})
	if err != nil {
		log.Fatal("Fucked: ",err)
	}
	fmt.Println("Address: ", address)
	fmt.Println("This is what happens when you print the type directly:", ipv4.ICMPTypeParameterProblem)
	fmt.Println("icmp.Message contents: ", ping)
	fmt.Println("icmp.Message.Marshal(): ", marshalled )
}
