package cmd

import (
	aw "github.com/deanishe/awgo"
	"github.com/google/uuid"
)

var u = uuid.New()

func init() {
	Register(&uuidGen{})
}

type uuidGen struct{ RunModeRun }

func (b *uuidGen) Use() string {
	return "uuid-gen"
}

func (b *uuidGen) ActionItem() *aw.Item {
	return wf.NewItem("UUID Generator").UID("B00000001-1").Valid(true).Icon(&aw.Icon{Value: "lock.png"})
}

func (b *uuidGen) Action(args []string) {
	plain := u.String()
	wf.NewItem("UUID: " + plain).Valid(true).Copytext(plain).Arg(plain).Icon(&aw.Icon{Value: "lock.png"})
}
