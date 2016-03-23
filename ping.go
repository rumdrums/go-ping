package main

import(
	"log"
	"os"
	"fmt"
	"net"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/icmp"
)

func incr(ip net.IP) {
    for j:= len(ip)-1; j>=0; j-- {
        ip[j]++
        if ip[j] > 0 {
            break
        }
    }
}

// closure to cache stuff that only needs
// to be done once:
func pingBuilder() func() []byte {
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
	return func() []byte {
		return marshalled
	}
}

func ping(c *icmp.PacketConn, address net.IP, queue chan int) {

	// this is pointless if done here:
	getPacket := pingBuilder()

	ping := getPacket()
	_, err := c.WriteTo(ping, &net.IPAddr{IP: address})
	if err != nil {
		log.Fatal("Write to socket: ", err)
	}
	queue <- 1
}

func getResponses(c *icmp.PacketConn, queue chan int, quit chan int) {

	fmt.Println("queue length: ", len(queue))
	//for len(queue) > 0 {
	for msg := range queue {
		fmt.Println("message: ", msg)
		buf := make([]byte, 1500)
		// I don't know what n is:
		// n, peer, err := c.ReadFrom(buf)
		_, peer, err := c.ReadFrom(buf)

		<- queue

		if err != nil {
			log.Fatal("Read from socket: ", err)
		}

		message, err := icmp.ParseMessage(1, buf)
		if err != nil {
			log.Fatal("Parse response: ", err)
		}

		fmt.Println("Message Type: ", message.Type)
		fmt.Println("Message Code: ", message.Code)
		//fmt.Println("Message Checksum: ", message.Checksum)
		//fmt.Println("Message Body: ", message.Body)

		// still don't know what n is:
		//fmt.Println("num bytes?: ", n)
		fmt.Println("read from: ", peer)
	}
	quit <- 1
}

func main() {
	var address string
	queue := make(chan int)
	quit := make(chan int)
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

	go func() {
		parsedIP, network, err := net.ParseCIDR(address)
		// not great -- currently assumes any error 
		// is caused by passing IP instead of CIDR...
		if err != nil {
			parsedIP := net.ParseIP(address)
			if parsedIP == nil {
				log.Fatal("Couldn't parse address")
			}
			fmt.Printf("Pinging %v\n",parsedIP)
			ping(c, parsedIP, queue)
		// ... and otherwise assumes it's CIDR and 
		// tries to iterate through subnet:
		} else {
			for ip := parsedIP.Mask(network.Mask); network.Contains(ip); incr(ip) {
				fmt.Printf("Pinging %v\n",ip)
				ping(c, ip, queue)
			}
		}
		close(queue)
	}()
	go getResponses(c, queue, quit)
	<-quit
}
