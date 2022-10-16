package app

import (
	"log"

	"github.com/net-byte/vtun/common/cipher"
	"github.com/net-byte/vtun/common/config"
	"github.com/net-byte/vtun/common/netutil"

	// "github.com/net-byte/vtun/grpc"
	"github.com/net-byte/vtun/tls"
	"github.com/net-byte/vtun/tun"
	"github.com/net-byte/vtun/udp"

	// "github.com/net-byte/vtun/ws"
	"github.com/net-byte/water"
)

var _banner = `
_                 
__ __ | |_   _  _   _ _  
\ V / |  _| | || | | ' \ 
 \_/   \__|  \_,_| |_||_|
						 
A simple VPN written in Go.
%s
`
var _srcUrl = "https://github.com/net-byte/vtun"

// vtun app struct
type App struct {
	Config  *config.Config
	Version string
	Iface   *water.Interface
	// IfaceMP1 *water.Interface
	// IfaceMP2 *water.Interface
	// IfaceMP3 *water.Interface
}

func NewApp(config *config.Config, version string) *App {

	return &App{
		Config:  config,
		Version: version,
	}
}

// InitConfig initializes the config
func (app *App) InitConfig() {
	log.Printf(_banner, _srcUrl)
	log.Printf("vtun version %s", app.Version)
	if !app.Config.ServerMode {
		app.Config.LocalGateway = netutil.DiscoverGateway(true)
		app.Config.LocalGatewayV6 = netutil.DiscoverGateway(false)
	}
	app.Config.BufferSize = 64 * 1024
	cipher.SetKey(app.Config.Key)
	app.Iface = tun.CreateTun(*app.Config)

	// 3 vnic
	// clonedConfig := app.Config
	// clonedConfig.CIDR = "172.16.0.11"
	// clonedConfig.DeviceName = "vtunMP1"
	// app.IfaceMP1 = tun.CreateTun(*clonedConfig)
	// clonedConfig.CIDR = "172.16.0.12"
	// clonedConfig.DeviceName = "vtunMP2"
	// app.IfaceMP2 = tun.CreateTun(*clonedConfig)
	// clonedConfig.CIDR = "172.16.0.13"
	// clonedConfig.DeviceName = "vtunMP3"
	// app.IfaceMP3 = tun.CreateTun(*clonedConfig)
	log.Printf("initialized config: %+v", app.Config)
	netutil.PrintStats(app.Config.Verbose)
}

// StartApp starts the app
func (app *App) StartApp() {

	switch app.Config.Protocol {
	case "udp":
		if app.Config.ServerMode {
			udp.StartServer(app.Iface, *app.Config)
		} else {
			udp.StartClient(app.Iface, *app.Config)
		}
	// case "ws", "wss":
	// 	if app.Config.ServerMode {
	// 		ws.StartServer(app.Iface, *app.Config)
	// 	} else {
	// 		ws.StartClient(app.Iface, *app.Config)
	// 	}
	case "tls":
		if app.Config.ServerMode {
			tls.StartServer(app.Iface, *app.Config)
		} else {
			tls.StartClient(app.Iface, *app.Config)
		}
	// case "grpc":
	// 	if app.Config.ServerMode {
	// 		grpc.StartServer(app.Iface, *app.Config)
	// 	} else {
	// 		grpc.StartClient(app.Iface, *app.Config)
	// 	}
	default:
		if app.Config.ServerMode {
			udp.StartServer(app.Iface, *app.Config)
		} else {
			udp.StartClient(app.Iface, *app.Config)
		}
	}
}

// StopApp stops the app
func (app *App) StopApp() {
	tun.ResetTun(*app.Config)
	app.Iface.Close()
	log.Println("vtun stopped")
}
