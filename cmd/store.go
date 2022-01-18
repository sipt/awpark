package cmd

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/spf13/cobra"

	aw "github.com/deanishe/awgo"
)

var userHomeDir, _ = os.UserHomeDir()
var searchLimit = 20
var store = &workflowStore{}

const cacheKey = "awpark-workflows"
const localDataFile = "workflows.json"
const remoteDataUrl = "https://raw.githubusercontent.com/TTNomi/awpark/main/static/workflows.json"

type workflowItem struct {
	Name    string   `json:"name"`
	Icon    string   `json:"icon"`
	Desc    string   `json:"desc"`
	Tags    []string `json:"tags"`
	Url     string   `json:"url"`
	Author  string   `json:"author"`
	Version string   `json:"version"`
	Website string   `json:"website"`
	Query   string   `json:"query"`
}

type workflowStore struct {
	items []*workflowItem
}

func (w *workflowStore) Search(keywords []string) {
	err := w.LoadData()
	if err != nil {
		wf.NewWarningItem(fmt.Sprintf("Error: %s", err.Error()), "")
		return
	}
	count := 0
	lowers := make([]string, len(keywords))
	for i, keyword := range keywords {
		lowers[i] = strings.ToLower(keyword)
	}
	for _, item := range w.items {
		notFound := false
		for _, keyword := range lowers {
			if !strings.Contains(item.Query, keyword) {
				notFound = true
				break
			}
		}
		if !notFound {
			wf.NewItem(item.Name+" @"+item.Author).Subtitle(item.Desc).Icon(&aw.Icon{
				Value: fmt.Sprintf(wf.Cache.Dir+"/icons/%x", md5.Sum([]byte(item.Icon))),
			}).Valid(true).Var("title", fmt.Sprintf("Downloading [%s]", item.Name)).
				Arg(item.Url).Var(CmdFlag, (&downloadFile{}).Use())
			count += 1
			if count > searchLimit {
				return
			}
		}
	}
}

func (w *workflowStore) LoadData() error {
	err := wf.Cache.LoadJSON(cacheKey, &w.items)
	if err != nil {
		err = w.loadDataRemote()
		return err
	}
	for _, item := range store.items {
		item.Query = strings.ToLower(item.Name + " " + strings.Join(item.Tags, " "))
	}
	duration, _ := wf.Cache.Age(cacheKey)
	if duration > 10*time.Minute {
		ci := &refreshCacheWorkflows{}
		bgCmd := exec.Command("./awpark", "exec", ci.Use())
		if !wf.IsRunning(ci.Use()) {
			log.Printf("[DEBUG] start run in background [%s]\n", bgCmd.String())
			err := wf.RunInBackground(ci.Use(), bgCmd)
			if err != nil {
				log.Printf("[ERROR] run in background failed [%s]:[%s]\n", bgCmd.String(), err.Error())
			}
		}
	}
	return nil
}

func (w *workflowStore) initData() error {
	data, err := ioutil.ReadFile(localDataFile)
	if err != nil {
		log.Printf("[ERROR] load local backup data failed [%s]", err.Error())
	} else {
		err = json.Unmarshal(data, &store.items)
		if err != nil {
			log.Printf("[ERROR] load local backup data failed [%s]", err.Error())
		}
		for _, item := range store.items {
			item.Query = strings.ToLower(item.Name + " " + strings.Join(item.Tags, " "))
		}
	}
	return w.loadDataRemote()
}

func (w *workflowStore) loadDataRemote() error {
	wg.Add(1)
	defer wg.Done()
	store.items = []*workflowItem{}

	resp, err := http.Get(remoteDataUrl)
	if err != nil {
		log.Printf("[ERROR] get remoteData failed [%s]\n", err.Error())
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] get remoteData failed [%s]\n", err.Error())
		return err
	}
	err = json.Unmarshal(data, &store.items)
	if err != nil {
		return err
	}
	for _, item := range store.items {
		item.Query = strings.ToLower(item.Name + " " + strings.Join(item.Tags, " "))
	}
	err = wf.Cache.Store(cacheKey, []byte(data))
	if err != nil {
		return err
	}

	ci := &cacheAllImages{}
	bgCmd := exec.Command("./awpark", "exec", ci.Use())
	if !wf.IsRunning(ci.Use()) {
		log.Printf("[DEBUG] start run in background [%s]\n", bgCmd.String())
		err := wf.RunInBackground(ci.Use(), bgCmd)
		if err != nil {
			log.Printf("[ERROR] run in background failed [%s]:[%s]\n", bgCmd.String(), err.Error())
		}
	}
	return nil
}

func (w *workflowStore) CacheAllImage() {
	err := w.LoadData()
	if err != nil {
		return
	}

	log.Printf("[DEBUG] starg cache image \n")
	wg.Add(len(w.items))
	for _, item := range w.items {
		go func(icon string) {
			file := fmt.Sprintf("%x", md5.Sum([]byte(icon)))
			localPath := wf.Cache.Dir + "/icons/" + file
			if _, err := os.Stat(localPath); os.IsNotExist(err) {
				_, err = grab.Get(localPath, icon)
				if err != nil {
					log.Printf("[ERROR] cache image [%s] to [%s] failed[%s]\n", icon, localPath, err.Error())
				} else {
					log.Printf("[DEBUG] cache image [%s] to [%s] success\n", icon, localPath)
				}
			}
			wg.Done()
		}(item.Icon)
	}
	wg.Wait()
}

func workflowList(cmd *cobra.Command, args []string) {
	wf.Run(func() {
		defer func() { wf.SendFeedback() }()
		if len(args) > 0 && len(args[0]) > 0 {
			urlPlain := args[0]
			validUrl, err := url.Parse(urlPlain)
			if err == nil && len(validUrl.Scheme) > 0 && len(validUrl.Host) > 0 && len(validUrl.Path) > 0 {
				wf.NewItem("Download & Install ...").Valid(true).Arg(urlPlain).Icon(&aw.Icon{Value: "lock.png"}).Var(CmdFlag, (&downloadFile{}).Use())
				return
			}
		}
		store.Search(args)
	})
}
func init() {
	Register(&downloadFile{})
	Register(&cacheAllImages{})
	Register(&refreshCacheWorkflows{})
}

type downloadFile struct{ RunModeNone }

func (d *downloadFile) Use() string {
	return "download"
}

func (d *downloadFile) ActionItem() *aw.Item {
	return nil
}

func (d *downloadFile) Action(args []string) {
	if len(args) > 0 && len(args[0]) > 0 {
		log.Printf("[DEBUG] start to download [%s]", args[0])
		resp, err := grab.Get(userHomeDir+"/Downloads/", args[0])
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}
		cmd := exec.Command("open", resp.Filename)
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}
		fmt.Printf("Download success. Wait Install...")
	}
}

type cacheAllImages struct{ RunModeNone }

func (d *cacheAllImages) Use() string {
	return "cache-images"
}

func (d *cacheAllImages) ActionItem() *aw.Item {
	return nil
}

func (d *cacheAllImages) Action(args []string) {
	store.CacheAllImage()
}

func copyFile(src, dest string) (err error) {
	defer func() {
		if err != nil {
			log.Printf("[ERROR] copy file failed [%s] to [%s], [%s]", src, dest, err.Error())
		}
	}()
	bytesRead, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dest, bytesRead, 0644)
}

type refreshCacheWorkflows struct{ RunModeNone }

func (d *refreshCacheWorkflows) Use() string {
	return "refresh-cache"
}

func (d *refreshCacheWorkflows) ActionItem() *aw.Item {
	return nil
}

func (d *refreshCacheWorkflows) Action(args []string) {
	log.Printf("[DEBUG] start to load remote data\n")
	err := store.loadDataRemote()
	if err != nil {
		log.Printf("[ERROR] load remote data failed [%s]\n", err.Error())
	}
}
