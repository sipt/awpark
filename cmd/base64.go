package cmd

import (
	"encoding/base64"
	"fmt"

	aw "github.com/deanishe/awgo"
)

func init() {
	Register(&base64Decoder{})
	Register(&base64Encoder{})
}

type base64Decoder struct{ RunModeRun }

func (b *base64Decoder) Use() string {
	return "base64-decode"
}

func (b *base64Decoder) ActionItem() *aw.Item {
	return wf.NewItem("Base64 Decode").UID("A00000001-1").Valid(true).Icon(&aw.Icon{Value: "lock.png"})
}

func (b *base64Decoder) Action(args []string) {
	if len(args) > 0 && len(args[0]) > 0 {
		plain, err := base64.StdEncoding.DecodeString(args[0])
		if err != nil {
			wf.NewItem(fmt.Sprintf("Error: %s", err.Error()))
		} else {
			wf.NewItem("Plain: " + string(plain)).Valid(true).Copytext(string(plain)).Arg(string(plain)).Icon(&aw.Icon{Value: "lock.png"}).Subtitle("Press [Enter], copy to the clipboard.")
		}
	}
}

type base64Encoder struct{ RunModeRun }

func (b *base64Encoder) Use() string {
	return "base64-encode"
}

func (b *base64Encoder) ActionItem() *aw.Item {
	return wf.NewItem("Base64 Encode").UID("A00000001-2").Valid(true).Icon(&aw.Icon{Value: "lock.png"})
}

func (b *base64Encoder) Action(args []string) {
	if len(args) > 0 && len(args[0]) > 0 {
		plain := base64.StdEncoding.EncodeToString([]byte(args[0]))
		wf.NewItem("Base64: " + plain).Valid(true).Copytext(plain).Arg(plain).Icon(&aw.Icon{Value: "lock.png"}).Subtitle("Press [Enter], copy to the clipboard.")
	}
}
