package main

import(
	"time"
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

func getPacket() []byte {
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

	return marshalled
}

func ping(c *icmp.PacketConn, address net.IP, queue *map[string]Ping) {

	fmt.Println(queue)
	ping := getPacket()
	_, err := c.WriteTo(ping, &net.IPAddr{IP: address})
	if err != nil {
		log.Fatal("Write to socket: ", err)
	}
}

func getResponses(c *icmp.PacketConn, queue *map[string]Ping, quit chan int) {

	fmt.Println(queue)
	for {
		buf := make([]byte, 1500)
		// this will block:
		_, peer, err := c.ReadFrom(buf)

		if err != nil {
			log.Fatal("Read from socket: ", err)
		}

		message, err := icmp.ParseMessage(1, buf)
		if err != nil {
			log.Fatal("Parse response: ", err)
		}

		fmt.Println("Message Type: ", message.Type)
		fmt.Println("read from: ", peer)
	}
	quit <- 1
}

type Ping struct {
	send time.Time
	receive time.Time
}

func main() {
	var address string
	queue := make(map[string]Ping)
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

	// closure -- execute whole block in background:
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
			ping(c, parsedIP, &queue)
		// ... and otherwise assumes it's CIDR and 
		// tries to iterate through subnet:
		} else {
			for ip := parsedIP.Mask(network.Mask); network.Contains(ip); incr(ip) {
				fmt.Printf("Pinging %v\n",ip)
				ping(c, ip, &queue)
			}
		}
	}()
	go getResponses(c, &queue, quit)
	<-quit
}
