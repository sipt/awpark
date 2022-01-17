package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	aw "github.com/deanishe/awgo"
)

func init() {
	Register(&timestampGetter{})
	Register(&timestampFormatter{})
}

type timestampGetter struct{ RunModeRun }

func (t *timestampGetter) Use() string {
	return "timestamp-getter"
}

func (t *timestampGetter) ActionItem() *aw.Item {
	return wf.NewItem("Timestamp Getter").Valid(true).Icon(&aw.Icon{Value: "clock.png"})
}

func (t *timestampGetter) Action(args []string) {
	now := time.Now()
	seconds := fmt.Sprintf("%d", now.Unix())
	milliseconds := fmt.Sprintf("%d", now.UnixNano()/int64(time.Millisecond))
	wf.NewItem("Seconds: " + seconds).Valid(true).Copytext(seconds).Arg(seconds).Icon(&aw.Icon{Value: "clock.png"})
	wf.NewItem("Milliseconds: " + milliseconds).Valid(true).Copytext(milliseconds).Arg(milliseconds).Icon(&aw.Icon{Value: "clock.png"})
}

type timestampFormatter struct{ RunModeRun }

func (t *timestampFormatter) Use() string {
	return "timestamp-formatter"
}

func (t *timestampFormatter) ActionItem() *aw.Item {
	return wf.NewItem("Timestamp Formatter").Valid(true).Icon(&aw.Icon{Value: "clock.png"})
}

func (t *timestampFormatter) Action(args []string) {
	if len(args) > 0 && len(args[0]) > 0 {
		strs := strings.Split(args[0], ":")
		if len(strs) <= 1 {
			return
		}
		timestamp, err := strconv.ParseInt(strs[1], 10, 64)
		if err != nil {
			wf.NewWarningItem("Invalid Timestamp", "convert to int64 failed")
			return
		}
		var inputTime time.Time
		switch strings.ToLower(strs[0]) {
		case "", "s":
			inputTime = time.Unix(timestamp, 0)
		case "ms":
			inputTime = time.Unix(0, timestamp*int64(time.Millisecond))
		case "us":
			inputTime = time.Unix(0, timestamp*int64(time.Microsecond))
		case "ns":
			inputTime = time.Unix(0, timestamp*int64(time.Nanosecond))
		default:
			wf.NewWarningItem("Invalid Timestamp", "flag:timestamp [flag not in (s, ms, us, ns)]")
			return
		}
		formatted := inputTime.Format(time.RFC3339)
		wf.NewItem("Formatted: " + formatted).Valid(true).Copytext(formatted).Arg(formatted).Icon(&aw.Icon{Value: "clock.png"})
	}
}
