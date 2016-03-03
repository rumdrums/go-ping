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

	marshalled, err := ping.Marshal(nil)
	if err != nil {
		log.Fatal("Marshall: ", err)
	}
	// why doesn't this fail when I pass
	// a nonsense address, like 192.168.1.1000 ??
	new_address := net.ParseIP(address)
	_, err = c.WriteTo(marshalled, &net.IPAddr{IP: new_address})
	if err != nil {
		log.Fatal("Write to socket: ", err)
	}

	buf := make([]byte, 1500)
	// I don't know what n is:
	n, peer, err := c.ReadFrom(buf)
	if err != nil {
        log.Fatal("Read from socket: ", err)
    }

	message, err := icmp.ParseMessage(1, buf)
    if err != nil {
        log.Fatal("Parse response: ", err)
    }

//	switch message.Type {
//	case ipv4.ICMPTypeEchoReply:
		fmt.Println("Message Type: ", message.Type)
		fmt.Println("Message Code: ", message.Code)
		fmt.Println("Message Checksum: ", message.Checksum)
		fmt.Println("Message Body: ", message.Body)
//	case ipv4.ICMPTypeDestinationUnreachable:
//		fmt.Println("Unreachable")
//	default:
//	}
	// still don't know what n is:
	fmt.Println("num bytes?: ", n)
	fmt.Println("read from: ", peer)

	// fmt.Println("Address: ", address)
	// fmt.Println("icmp.Message.Marshal(): ", marshalled )
	// fmt.Println("icmp.Message contents: ", ping)
	// fmt.Println("This is what happens when you print the type directly:", ipv4.ICMPTypeParameterProblem)

}
