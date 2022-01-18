package cmd

import (
	"fmt"
	"net/url"

	aw "github.com/deanishe/awgo"
)

func init() {
	Register(&urlDecoder{})
	Register(&urlEncoder{})
}

type urlDecoder struct{ RunModeRun }

func (b *urlDecoder) Use() string {
	return "url-decode"
}

func (b *urlDecoder) ActionItem() *aw.Item {
	return wf.NewItem("URL Decode").UID("A00000003-1").Valid(true).Icon(&aw.Icon{Value: "lock.png"})
}

func (b *urlDecoder) Action(args []string) {
	if len(args) > 0 && len(args[0]) > 0 {
		plain, err := url.QueryUnescape(args[0])
		if err != nil {
			wf.NewItem(fmt.Sprintf("Error: %s", err.Error()))
		} else {
			wf.NewItem("URL Unescape: " + plain).Valid(true).Copytext(plain).Arg(plain).Icon(&aw.Icon{Value: "lock.png"}).Subtitle("Press [Enter], copy to the clipboard.")
		}
	}
}

type urlEncoder struct{ RunModeRun }

func (b *urlEncoder) Use() string {
	return "url-encode"
}

func (b *urlEncoder) ActionItem() *aw.Item {
	return wf.NewItem("URL Encode").UID("A00000003-2").Valid(true).Icon(&aw.Icon{Value: "lock.png"})
}

func (b *urlEncoder) Action(args []string) {
	if len(args) > 0 && len(args[0]) > 0 {
		value := url.QueryEscape(args[0])
		wf.NewItem("URL Escape: " + value).Valid(true).Copytext(value).Arg(value).Icon(&aw.Icon{Value: "lock.png"}).Subtitle("Press [Enter], copy to the clipboard.")
	}
}
