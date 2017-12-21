package main

import (
	"flag"
	"sync"
	"net"
	"io"
	"log"
	"fmt"
	"strings"
	"github.com/cheikhshift/gos/core"
)

// Data structures to manage 
// web server instances
type StaticHost struct {
	Lock  *sync.RWMutex
	Cache map[string]int
}

func NewCache() StaticHost {
	return StaticHost{Lock: new(sync.RWMutex), Cache: make(map[string]int)}
}


var (
	PortApp,App string 
	Limit,IpInc int
	Count int = 2
	LANnet = []string{}

	Host StaticHost
	
	bspath = "./launcher.sh"
)

func GetServerAvailable() string {
	var index string
	Host.Lock.Lock()
	defer Host.Lock.Unlock()

	for index,concount := range Host.Cache {

		if concount < Limit  {
			Host.Cache[index] += IpInc
			return index
		}
	}
	if App != "" {
		core.RunCmd(App)
	}

	lsize := len(LANnet) - IpInc
	
	LANnet[lsize] = fmt.Sprintf("%v:%s",Count, PortApp )
	Count++
	index = strings.Join(LANnet,".")
	//run bash
	Host.Cache[index] = IpInc
	return index
}

func main() {
	lnch := flag.String("app", "", "Run specified terminal command each time a new instance is needed.")
	dhcpstr := flag.String("net", "192.168.0.1", "Router LAN IP address (DHCP subnet).")
	maxcon := flag.Int("max", 100, "Maximum number of connections per instance. clm-static will divide tasks.")
	addby := flag.Int("incby", 1, "Increase last octet of IP when picking the next server.")
	dhcpstart := flag.Int("start", 1, "First DHCP assigned IP of your instances (Initial value of last octet of IPv4 address). Example : with value 21, this tool will assume your first instance's ip will be 192.168.0.21")
	apport := flag.String("appPort", "8080", "Port your instances will listen on.")
	port := flag.String("port", "9000",  "Port clm-static should listen on.")

	flag.Parse()


	App = *lnch
	Limit = *maxcon
	IpInc = *addby
	LANnet = strings.Split(*dhcpstr, ".")
	Count = *dhcpstart
	Host = NewCache()
	PortApp = *	apport

	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", *port) )
	if err != nil {
		panic(err)
	}



	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {

	ipaddr := GetServerAvailable()
	proxy, err := net.Dial("tcp",  ipaddr )
	if err != nil {
		
		Host.Lock.Lock()
		defer Host.Lock.Unlock()
		defer handleRequest(conn)
		delete(Host.Cache, ipaddr)
		return
	}
	
	go copyIO(conn, proxy,"")
	go copyIO(proxy, conn,ipaddr)
}

func copyIO(src, dest net.Conn, index string) {
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
	if index != "" {
		Host.Lock.Lock()
		defer Host.Lock.Unlock()
		Host.Cache[index] -=  IpInc
		// fmt.Println(Host)
	}


}
