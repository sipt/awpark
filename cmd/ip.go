package cmd

import (
	"io/ioutil"
	"net"
	"net/http"

	aw "github.com/deanishe/awgo"
)

func init() {
	Register(&ipTool{})
}

type ipTool struct{ RunModeRun }

func (b *ipTool) Use() string {
	return "ip"
}

func (b *ipTool) ActionItem() *aw.Item {
	return wf.NewItem("Show IP Address").UID("B00000003-1").Valid(true).Icon(&aw.Icon{Value: "ip.png"})
}

func (b *ipTool) Action(args []string) {
	localIp, err := GetOutboundIP()
	if err != nil {
		wf.NewWarningItem("Local IP ", err.Error())
	} else {
		wf.NewItem("Local IP: " + localIp).Valid(true).Copytext(localIp).Arg(localIp).Icon(&aw.Icon{Value: "ip.png"}).Subtitle("Press [Enter], copy to the clipboard.")
	}
	resp, err := http.Get("http://icanhazip.com")
	var data []byte
	if err == nil {
		data, err = ioutil.ReadAll(resp.Body)
	}
	resp.Body.Close()
	if err != nil {
		wf.NewWarningItem("External IP ", err.Error())
	} else {
		if len(data) > 0 && data[len(data)-1] == '\n' {
			data = data[:len(data)-1]
		}
		wf.NewItem("External IP: " + string(data)).Valid(true).Copytext(string(data)).Arg(string(data)).Icon(&aw.Icon{Value: "ip.png"}).Subtitle("Press [Enter], copy to the clipboard.")
	}
}

func GetOutboundIP() (ipStr string, err error) {
	iface, err := net.InterfaceByName("en0")
	if err != nil {
		return
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return
	}

	var ip net.IP
	switch v := addrs[len(addrs)-1].(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	ipStr = ip.String()
	return
}
