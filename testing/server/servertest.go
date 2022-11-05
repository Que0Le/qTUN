package main

import (
	"fmt"
	"net"
)

func sendResponse( /* conn *net.UDPConn,*/ addr *net.UDPAddr) {
	conn, err := net.Dial("udp", "172.16.0.10:1234")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	fmt.Fprintf(
		conn,
		"This is server. Message received from you, "+addr.IP.String(),
	)
}

func main() {
	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("172.16.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	for {
		_, remoteaddr, err := ser.ReadFromUDP(p)
		fmt.Printf("Read a message from %v %s \n", remoteaddr, p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		go sendResponse(remoteaddr)
	}
}
