package udp

import (
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/golang/snappy"
	"github.com/net-byte/vtun/common/cipher"
	"github.com/net-byte/vtun/common/config"
	"github.com/net-byte/vtun/common/counter"
	"github.com/net-byte/vtun/common/netutil"
	"github.com/net-byte/water"
)

// StartClient starts the udp client
func StartClient(iface *water.Interface, config config.Config) {
	serverAddr, err := net.ResolveUDPAddr("udp", config.ServerAddr)
	if err != nil {
		log.Fatalln("failed to resolve server addr:", err)
	}
	localAddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		log.Fatalln("failed to get udp socket:", err)
	}
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatalln("failed to listen on udp socket:", err)
	}
	/*  */
	localAddr1, err := net.ResolveUDPAddr("udp", "192.168.122.101:0")
	if err != nil {
		log.Fatalln("failed to get udp socket:", err)
	}
	conn1, err := net.ListenUDP("udp", localAddr1)
	if err != nil {
		log.Fatalln("failed to listen on udp socket:", err)
	}
	localAddr2, err := net.ResolveUDPAddr("udp", "192.168.122.102:0")
	if err != nil {
		log.Fatalln("failed to get udp socket:", err)
	}
	conn2, err := net.ListenUDP("udp", localAddr2)
	if err != nil {
		log.Fatalln("failed to listen on udp socket:", err)
	}
	localAddr3, err := net.ResolveUDPAddr("udp", "192.168.122.103:0")
	if err != nil {
		log.Fatalln("failed to get udp socket:", err)
	}
	conn3, err := net.ListenUDP("udp", localAddr3)
	if err != nil {
		log.Fatalln("failed to listen on udp socket:", err)
	}
	/*  */
	defer conn.Close()
	log.Printf("vtun udp client started on %v", conn.LocalAddr().String())
	c := &Client{
		config: config, iface: iface, localConn: conn, serverAddr: serverAddr,
		localConnMP1: conn1, localConnMP2: conn2, localConnMP3: conn3,
	}
	go c.udpToTun()
	c.tunToUdp()
}

// The client struct
type Client struct {
	config       config.Config
	iface        *water.Interface
	localConn    *net.UDPConn
	localConnMP1 *net.UDPConn
	localConnMP2 *net.UDPConn
	localConnMP3 *net.UDPConn
	serverAddr   *net.UDPAddr
}

// udpToTun sends packets from udp to tun
func (c *Client) udpToTun() {
	packet := make([]byte, c.config.BufferSize)
	for {
		n, _, err := c.localConn.ReadFromUDP(packet)
		if err != nil || n == 0 {
			netutil.PrintErr(err, c.config.Verbose)
			continue
		}
		b := packet[:n]
		if c.config.Compress {
			b, err = snappy.Decode(nil, b)
			if err != nil {
				netutil.PrintErr(err, c.config.Verbose)
				continue
			}
		}
		if c.config.Obfs {
			b = cipher.XOR(b)
		}
		c.iface.Write(b)
		counter.IncrReadBytes(n)
	}
}

// tunToUdp sends packets from tun to udp
func (c *Client) tunToUdp() {
	rand.Seed(time.Now().UnixNano())
	packet := make([]byte, c.config.BufferSize)
	for {
		n, err := c.iface.Read(packet)
		if err != nil {
			netutil.PrintErr(err, c.config.Verbose)
			break
		}
		b := packet[:n]
		if c.config.Obfs {
			b = cipher.XOR(b)
		}
		if c.config.Compress {
			b = snappy.Encode(nil, b)
		}
		rand := 1 + rand.Intn(3-1+1)
		println(rand)
		if rand == 1 {
			c.localConnMP1.WriteToUDP(b, c.serverAddr)
		} else if rand == 2 {
			c.localConnMP2.WriteToUDP(b, c.serverAddr)
		} else if rand == 3 {
			c.localConnMP3.WriteToUDP(b, c.serverAddr)
		}
		// c.localConn.WriteToUDP(b, c.serverAddr)
		counter.IncrWrittenBytes(n)
	}
}
