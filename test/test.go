package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

func main() {
	resp, err := http.Get("http://icanhazip.com")
	var data []byte
	if err != nil {
		data, err = ioutil.ReadAll(resp.Body)
	}
	resp.Body.Close()
	fmt.Println(resp, string(data))
}
func GetOutboundIP() string {
	iface, err := net.InterfaceByName("en0")
	if err != nil {
		panic(err)
	}
	addrs, err := iface.Addrs()
	if err != nil {
		panic(err)
	}

	var ip net.IP
	switch v := addrs[len(addrs)-1].(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	return ip.String()
}
func loop() {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	// handle err
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			panic(err)
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// process IP address
			fmt.Println(i, ip)
		}
	}

}
