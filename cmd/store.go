package cmd

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
const cacheImagesKey = "awpark-workflows-cache-images"

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
			}).Valid(true).Arg(item.Url).Var(CmdFlag, (&downloadFile{}).Use())
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
		go w.loadDataRemote()
	}
	return nil
}

func (w *workflowStore) loadDataRemote() error {
	wg.Add(1)
	defer wg.Done()
	store.items = []*workflowItem{}
	data := `[
    {
        "website": "http://alfred-spotify-mini-player.com/",
        "icon": "https://www.alfredapp.com/media/pages/workflows/macapps/spotify.png",
        "name": "Spotify Mini Player",
        "desc": "The Mini Player gives you control over the Spotify app on your Mac. Find albums, search for artists & play songs to groove to."
    },
    {
        "website": "https://github.com/Kapeli/Dash-Alfred-Workflow",
        "icon": "https://www.alfredapp.com/media/pages/workflows/macapps/dash.png",
        "name": "Dash",
        "desc": "Add blisteringly fast search of the Dash documentation / API browser app, with in-line results and responsive integration."
    },
    {
        "website": "https://github.com/rhydlewis/search-omnifocus#readme",
        "icon": "https://www.alfredapp.com/media/pages/workflows/macapps/omnifocus.png",
        "name": "OmniFocus",
        "desc": "Search for your projects, folders and tasks in OmniFocus 3.",
		"url":"https://github.com/rhydlewis/search-omnifocus/releases/download/v2.1.3/Search.OmniFocus.v2.1.3.alfredworkflow"
    },
    {
        "website": "https://github.com/core-code/MacUpdater-Alfred-Workflow#readme",
        "icon": "https://www.alfredapp.com/media/pages/workflows/macapps/macupdater.png",
        "name": "MacUpdater",
        "desc": "Keep your Mac software up to date effortlessly."
    },
    {
        "website": "http://www.packal.org/workflow/transmit",
        "icon": "https://www.alfredapp.com/media/pages/workflows/macapps/transmit.png",
        "name": "Transmit",
        "desc": "Search and open Favorites in the Transmit 4 FTP client."
    },
    {
        "website": "http://www.packal.org/workflow/launch-url-3-browsers",
        "icon": "https://www.alfredapp.com/media/pages/workflows/handy/launch-in-3-browsers.png",
        "name": "Launch a URL in 3 browsers",
        "desc": "Launch a URL in three browsers (Safari, Chrome and Firefox) to speed up website testing."
    },
    {
        "website": "http://www.packal.org/workflow/lastpass-cli-workflow-alfred",
        "icon": "https://www.alfredapp.com/media/pages/workflows/macapps/lastpass.png",
        "name": "LastPass",
        "desc": "Search for your logins saved in the LastPass password manager."
    },
    {
        "website": "https://www.apptorium.com/sidenotes/articles/alfred-workflow",
        "icon": "https://www.alfredapp.com/media/pages/workflows/macapps/sidenotes.png",
        "name": "SideNotes",
        "desc": "Create new notes, folders, and search your existing SideNotes."
    },
    {
        "website": "/help/workflows/templates/#suggest",
        "icon": "https://www.alfredapp.com/media/pages/workflows/websites/google.png",
        "name": "Google Suggest",
        "desc": "An example included in Alfred: Search Google from Alfred's search box and see results in-line."
    },
    {
        "website": "https://github.com/deanishe/alfred-stackexchange#readme",
        "icon": "https://www.alfredapp.com/media/pages/workflows/websites/stackoverflow.png",
        "name": "Stack Exchange",
        "desc": "Search for answers to your programming questions on Stack Overflow."
    },
    {
        "website": "https://github.com/edgarjs/alfred-github-repos#readme",
        "icon": "https://www.alfredapp.com/media/pages/workflows/websites/github.png",
        "name": "GitHub",
        "desc": "Quickly pick which GitHub repository you want to open and launch it from Alfred."
    },
    {
        "website": "https://github.com/deanishe/alfred-reddit#readme",
        "icon": "https://www.alfredapp.com/media/pages/workflows/websites/reddit.png",
        "name": "Reddit",
        "desc": "Browse and search Reddit directories (subreddits), and search hot results within a subreddit."
    },
    {
        "website": "https://github.com/alfredapp/google-drive-workflow#readme",
        "icon": "https://www.alfredapp.com/media/pages/workflows/websites/googledrive.png",
        "name": "Google Drive",
        "desc": "List File Stream contents from Google Drive."
    },
    {
        "website": "https://github.com/alfredapp/tinypng-workflow#readme",
        "icon": "https://www.alfredapp.com/media/pages/workflows/websites/tinypng.png",
        "name": "TinyPNG",
        "desc": "Optimise your images to be more lightweight with TinyPNG"
    },
    {
        "website": "/help/workflows/templates/#suggest",
        "icon": "https://www.alfredapp.com/media/pages/workflows/websites/amazon.png",
        "name": "Amazon Suggest",
        "desc": "An example included in Alfred: Search Amazon from Alfred's search box and see results in-line."
    },
    {
        "website": "https://github.com/vitorgalvao/alfred-workflows/tree/master/PinPlus#readme",
        "icon": "https://www.alfredapp.com/media/pages/workflows/websites/pinboard.png",
        "name": "PinPlus",
        "desc": "Add and view your Pinboard Bookmarks."
    },
    {
        "website": "https://github.com/tmcknight/Movie-and-TV-Show-Search-Alfred-Workflow#readme",
        "icon": "https://www.alfredapp.com/media/pages/workflows/websites/movieandtv.png",
        "name": "Movie & TV Show Search",
        "desc": "Search for a movie or TV show, and get a few ratings."
    },
    {
        "website": "https://github.com/clarencecastillo/alfred-powerthesaurus#readme",
        "icon": "https://www.alfredapp.com/media/pages/workflows/websites/powerthesaurus.jpg",
        "name": "Power Thesaurus",
        "desc": "Search in-line for synonyms and antonyms on Power Thesaurus."
    },
    {
        "website": "https://github.com/mdreizin/chrome-bookmarks-alfred-workflow#readme",
        "icon": "https://www.alfredapp.com/media/pages/workflows/websites/chromium.png",
        "name": "Chromium Bookmarks",
        "desc": "Search for your browser bookmarks in Chrome, Chromium, Edge and Vivaldi."
    },
    {
        "website": "https://github.com/stuartcryan/advanced-google-maps-alfred-workflow",
        "icon": "https://www.alfredapp.com/media/pages/workflows/websites/googlemaps.png",
        "name": "Advanced Maps",
        "desc": "Advanced Google and Apple Maps search, including the ability to configure a Home and Work location, to see Google traffic reports before travelling."
    },
    {
        "website": "https://github.com/deanishe/alfred-convert#readme",
        "icon": "https://www.alfredapp.com/media/pages/workflows/handy/convert.png",
        "name": "Convert",
        "desc": "Convert between different units of quantities, distances, time and more. No internet connection required."
    },
    {
        "website": "https://github.com/deanishe/alfred-smartfolders",
        "icon": "https://www.alfredapp.com/media/pages/workflows/handy/smart-folders.png",
        "name": "Smart Folders",
        "desc": "View and explore your Smart Folders."
    },
    {
        "website": "https://github.com/jaroslawhartman/TimeZones-Alfred#readme",
        "icon": "https://www.alfredapp.com/media/pages/workflows/handy/timezones.png",
        "name": "Timezones",
        "desc": "A customised world clock — shows a list of user-configured cities with the current local times."
    },
    {
        "website": "https://github.com/vitorgalvao/alfred-workflows/tree/master/CoffeeCoffee#readme",
        "icon": "https://www.alfredapp.com/media/pages/workflows/handy/coffeecoffee.png",
        "name": "CoffeeCoffee",
        "desc": "Prevent your Mac from going to sleep indefinitely or for a set amount of time."
    },
    {
        "website": "https://github.com/shmulvad/alfred-nightshift#readme",
        "icon": "https://www.alfredapp.com/media/pages/workflows/handy/nightshift.png",
        "name": "NightShift",
        "desc": "Save your eyes by turning macOS's Night Shift feature on and off quickly."
    },
    {
        "website": "https://github.com/alfredapp/magic-8-ball-workflow#readme",
        "icon": "https://www.alfredapp.com/media/pages/workflows/handy/magic-8ball.png",
        "name": "Magic 8 Ball",
        "desc": "Let the Magic 8 Ball answer your important questions."
    },
    {
        "website": "http://www.packal.org/workflow/wi-fi",
        "icon": "https://www.alfredapp.com/media/pages/workflows/handy/wi-fi.png",
        "name": "Wi-Fi On/Off",
        "desc": "Set your Mac's Wi-Fi function to On or Off."
    },
    {
        "website": "http://www.packal.org/workflow/caffeinate-control",
        "icon": "https://www.alfredapp.com/media/pages/workflows/handy/caffeinate-control.png",
        "name": "Caffeinate Control",
        "desc": "Controls \"caffeinate\" — system command to prevent the sleep function."
    }
]`
	err := json.Unmarshal([]byte(data), &store.items)
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
