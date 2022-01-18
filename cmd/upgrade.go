package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/hashicorp/go-version"

	aw "github.com/deanishe/awgo"
)

const CacheVersionKey = "awpark-version"
const UpgradeUrl = "https://github.com/TTNomi/awpark/releases/latest/download/awpark.alfredworkflow"

var vm = new(versionManager)

func NeedUpgrade() bool {
	ver, _ := vm.CheckVersion()
	if ver != "" {
		wf.NewItem(fmt.Sprintf("AWPark has new version [%s]", ver)).Subtitle("Press [Enter] key to upgrade.").
			Var("upgrade", "true").Arg(UpgradeUrl).Icon(&aw.Icon{
			Value: "icon.png",
		}).Valid(true).Var("title", "Downloading [AWPark]").Var(CmdFlag, (&downloadFile{}).Use()).
			Var("website", "https://github.com/TTNomi/awpark")
	}
	return ver != ""
}

type versionManager struct{}

func (v *versionManager) FetchVersion() (string, error) {
	log.Printf("[DEBUG] start fetch version\n")
	req, err := http.NewRequest("HEAD", "https://github.com/TTNomi/awpark/releases/latest", nil)
	if err != nil {
		log.Printf("[ERROR] make request failed [%s]", err.Error())
		return "", err
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] check version failed [%s]", err.Error())
		return "", err
	}
	lv := resp.Header.Get("Location")
	ver := lv[strings.LastIndex(lv, "/")+1:]
	if ver == "releases" {
		ver = "1.0"
	}
	log.Printf("[DEBUG] end fetch version [%s]\n", ver)
	_, err = version.NewVersion(ver)
	if err != nil {
		log.Printf("[ERROR] new version invalid [%s][%s]", ver, err.Error())
		return "", err
	}
	err = wf.Cache.Store(CacheVersionKey, []byte(ver))
	if err != nil {
		log.Printf("[ERROR] store new version failed [%s][%s]", ver, err.Error())
		return "", err
	}
	return ver, nil
}

func (v *versionManager) CheckVersion() (string, error) {
	nowVer, err := version.NewVersion(wf.Version())
	if err != nil {
		log.Printf("[ERROR] current version invalid [%s][%s]", wf.Version(), err.Error())
		return "", err
	}
	ver, err := wf.Cache.Load(CacheVersionKey)
	if err != nil {
		log.Printf("[ERROR] load new version failed [%s]", err.Error())
		v.FetchVersionInBg()
		return "", err
	}
	d, _ := wf.Cache.Age(CacheVersionKey)
	if d > time.Minute*30 {
		v.FetchVersionInBg()
	}
	newVer, _ := version.NewVersion(string(ver))
	if newVer.GreaterThan(nowVer) {
		return string(ver), nil
	}
	return "", nil
}

func (v *versionManager) FetchVersionInBg() {
	ci := &versionGetter{}
	bgCmd := exec.Command("./awpark", "exec", ci.Use())
	if !wf.IsRunning(ci.Use()) {
		log.Printf("[DEBUG] start run in background [%s]\n", bgCmd.String())
		err := wf.RunInBackground(ci.Use(), bgCmd)
		if err != nil {
			log.Printf("[ERROR] run in background failed [%s]:[%s]\n", bgCmd.String(), err.Error())
		}
	}
}

func init() {
	Register(&versionGetter{})
}

type versionGetter struct{ RunModeNone }

func (v *versionGetter) Use() string {
	return "version-getter"
}

func (v *versionGetter) ActionItem() *aw.Item {
	return nil
}

func (v *versionGetter) Action(args []string) {
	vm.FetchVersion()
}
