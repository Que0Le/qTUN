package tls

import (
	"crypto/tls"
	"log"
	"net"
	"time"

	"github.com/golang/snappy"
	"github.com/net-byte/vtun/common/cache"
	"github.com/net-byte/vtun/common/cipher"
	"github.com/net-byte/vtun/common/config"
	"github.com/net-byte/vtun/common/counter"
	"github.com/net-byte/vtun/common/netutil"
	"github.com/net-byte/water"
)

// StartClient starts the tls client
func StartClient(iface *water.Interface, config config.Config) {
	log.Println("vtun tls client started")
	go tunToTLS(config, iface)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: config.TLSInsecureSkipVerify,
	}
	if config.TLSSni != "" {
		tlsconfig.ServerName = config.TLSSni
	}
	for {
		conn, err := tls.Dial("tcp", config.ServerAddr, tlsconfig)
		if err != nil {
			time.Sleep(3 * time.Second)
			netutil.PrintErr(err, config.Verbose)
			continue
		}
		cache.GetCache().Set("tlsconn", conn, 24*time.Hour)
		tlsToTun(config, conn, iface)
		cache.GetCache().Delete("tlsconn")
	}
}

// tunToTLS sends packets from tun to tls
func tunToTLS(config config.Config, iface *water.Interface) {
	packet := make([]byte, config.BufferSize)
	for {
		n, err := iface.Read(packet)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			break
		}
		if v, ok := cache.GetCache().Get("tlsconn"); ok {
			b := packet[:n]
			if config.Obfs {
				b = cipher.XOR(b)
			}
			if config.Compress {
				b = snappy.Encode(nil, b)
			}
			tlsconn := v.(net.Conn)
			tlsconn.SetWriteDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
			_, err = tlsconn.Write(b)
			if err != nil {
				netutil.PrintErr(err, config.Verbose)
				continue
			}
			counter.IncrWrittenBytes(n)
		}
	}
}

// tlsToTun sends packets from tls to tun
func tlsToTun(config config.Config, tlsconn net.Conn, iface *water.Interface) {
	defer tlsconn.Close()
	packet := make([]byte, config.BufferSize)
	for {
		tlsconn.SetReadDeadline(time.Now().Add(time.Duration(config.Timeout) * time.Second))
		n, err := tlsconn.Read(packet)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			break
		}
		b := packet[:n]
		if config.Compress {
			b, err = snappy.Decode(nil, b)
			if err != nil {
				netutil.PrintErr(err, config.Verbose)
				break
			}
		}
		if config.Obfs {
			b = cipher.XOR(b)
		}
		_, err = iface.Write(b)
		if err != nil {
			netutil.PrintErr(err, config.Verbose)
			break
		}
		counter.IncrReadBytes(n)
	}
}
