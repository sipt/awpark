package cmd

import (
	"encoding/base64"
	"fmt"
	"strings"

	aw "github.com/deanishe/awgo"
)

func init() {
	Register(&jwtDecoder{})
}

type jwtDecoder struct{ RunModeRun }

func (b *jwtDecoder) Use() string {
	return "jwt-decode"
}

func (b *jwtDecoder) ActionItem() *aw.Item {
	return wf.NewItem("JWT Decode").UID("B00000002-1").Valid(true).Icon(&aw.Icon{Value: "jwt.png"})
}

func (b *jwtDecoder) Action(args []string) {
	if len(args) > 0 && len(args[0]) > 0 {
		subStrs := strings.Split(args[0], ".")
		if len(subStrs) < 3 {
			wf.NewItem(fmt.Sprintf("Error: %s", "invalid jwt token"))
		}
		// header
		plain, err := base64.RawStdEncoding.DecodeString(subStrs[0])
		if err != nil {
			wf.NewItem(fmt.Sprintf("Error: %s", err.Error()))
		} else {
			wf.NewItem("Header: " + string(plain)).Valid(true).Copytext(string(plain)).Arg(string(plain)).Icon(&aw.Icon{Value: "jwt.png"})
		}
		// header
		plain, err = base64.RawStdEncoding.DecodeString(subStrs[1])
		if err != nil {
			wf.NewItem(fmt.Sprintf("Error: %s", err.Error()))
		} else {
			wf.NewItem("Payload: " + string(plain)).Valid(true).Copytext(string(plain)).Arg(string(plain)).Icon(&aw.Icon{Value: "jwt.png"})
		}
	}
}
