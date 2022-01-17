package cmd

import (
	"crypto/md5"
	"fmt"
	"strings"

	aw "github.com/deanishe/awgo"
)

func init() {
	Register(&md5Lower{})
	Register(&md5Upper{})
}

type md5Lower struct{ RunModeRun }

func (b *md5Lower) Use() string {
	return "md5-lower"
}

func (b *md5Lower) ActionItem() *aw.Item {
	return wf.NewItem("MD5 Lower").UID("A00000002-1").Valid(true).Icon(&aw.Icon{Value: "lock.png"})
}

func (b *md5Lower) Action(args []string) {
	if len(args) > 0 && len(args[0]) > 0 {
		data := md5.Sum([]byte(args[0]))
		str := strings.ToLower(fmt.Sprintf("%x", data))
		wf.NewItem("MD5: " + str).Valid(true).Copytext(str).Arg(str).Icon(&aw.Icon{Value: "lock.png"})
	}
}

type md5Upper struct{ RunModeRun }

func (b *md5Upper) Use() string {
	return "md5-upper"
}

func (b *md5Upper) ActionItem() *aw.Item {
	return wf.NewItem("MD5 Upper").UID("A00000002-2").Valid(true).Icon(&aw.Icon{Value: "lock.png"})
}

func (b *md5Upper) Action(args []string) {
	if len(args) > 0 && len(args[0]) > 0 {
		data := md5.Sum([]byte(args[0]))
		str := strings.ToUpper(fmt.Sprintf("%x", data))
		wf.NewItem("MD5: " + str).Valid(true).Copytext(str).Arg(str).Icon(&aw.Icon{Value: "lock.png"})
	}
}
