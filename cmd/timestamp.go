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
	Register(&timestampParser{})
}

type timestampGetter struct{ RunModeRun }

func (t *timestampGetter) Use() string {
	return "timestamp-getter"
}

func (t *timestampGetter) ActionItem() *aw.Item {
	return wf.NewItem("Timestamp Getter").UID("B00000004-1").Valid(true).Icon(&aw.Icon{Value: "clock.png"})
}

func (t *timestampGetter) Action(args []string) {
	now := time.Now()
	seconds := fmt.Sprintf("%d", now.Unix())
	milliseconds := fmt.Sprintf("%d", now.UnixNano()/int64(time.Millisecond))
	wf.NewItem("Seconds: " + seconds).Valid(true).Copytext(seconds).Arg(seconds).Icon(&aw.Icon{Value: "clock.png"})
	wf.NewItem("Milliseconds: " + milliseconds).Valid(true).Copytext(milliseconds).Arg(milliseconds).Icon(&aw.Icon{Value: "clock.png"}).Subtitle("Press [Enter], copy to the clipboard.")
}

type timestampFormatter struct{ RunModeRun }

func (t *timestampFormatter) Use() string {
	return "timestamp-formatter"
}

func (t *timestampFormatter) ActionItem() *aw.Item {
	return wf.NewItem("Timestamp Formatter").UID("B00000004-2").Valid(true).Icon(&aw.Icon{Value: "clock.png"})
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
		wf.NewItem("Formatted: " + formatted).Valid(true).Copytext(formatted).Arg(formatted).Icon(&aw.Icon{Value: "clock.png"}).Subtitle("Press [Enter], copy to the clipboard.")
	}
}

type timestampParser struct{ RunModeRun }

func (t *timestampParser) Use() string {
	return "timestamp-parser"
}

func (t *timestampParser) ActionItem() *aw.Item {
	return wf.NewItem("Timestamp Parser").UID("B00000004-3").Valid(true).Icon(&aw.Icon{Value: "clock.png"})
}

func (t *timestampParser) Action(args []string) {
	formatTemplate := []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05.999999999Z07:00",
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02 15:04:05.999999999Z07:00",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05 -0700 MST",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05 -07:00",
		"2006-01-02 15:04:05.999999999 -0700 MST",
		"2006-01-02 15:04:05.999999999 -0700",
		"2006-01-02 15:04:05.999999999 -07:00",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"01/02 03:04:05PM '06 -0700",
		"Mon Jan _2 15:04:05 2006",
		"Mon Jan _2 15:04:05 MST 2006",
		"Mon Jan 02 15:04:05 -0700 2006",
		"02 Jan 06 15:04 MST",
		"02 Jan 06 15:04 -0700",
		"Monday, 02-Jan-06 15:04:05 MST",
		"Mon, 02 Jan 2006 15:04:05 MST",
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"3:04PM",
		"Jan _2 15:04:05",
		"Jan _2 15:04:05.000",
		"Jan _2 15:04:05.000000",
		"Jan _2 15:04:05.000000000",
	}
	tips := func() {
		wf.NewItem("Reference the following date format.").Valid(false).Icon(&aw.Icon{Value: "clock.png"})
		for _, t := range formatTemplate {
			wf.NewItem("eg. " + t).Valid(false).Icon(&aw.Icon{Value: "clock.png"})
		}
	}
	if len(args) > 0 && len(args[0]) > 0 {
		timestr := args[0]
		var (
			inputTime time.Time
			err       error
		)
		for _, format := range formatTemplate {
			inputTime, err = time.ParseInLocation(format, timestr, time.Local)
			if err == nil {
				break
			}
		}
		if err != nil {
			tips()
			return
		}
		seconds := fmt.Sprintf("%d", inputTime.Unix())
		milliseconds := fmt.Sprintf("%d", inputTime.UnixNano()/int64(time.Millisecond))
		wf.NewItem("Seconds: " + seconds).Valid(true).Copytext(seconds).Arg(seconds).Icon(&aw.Icon{Value: "clock.png"})
		wf.NewItem("Milliseconds: " + milliseconds).Valid(true).Copytext(milliseconds).Arg(milliseconds).Icon(&aw.Icon{Value: "clock.png"}).Subtitle("Press [Enter], copy to the clipboard.")
	} else {
		tips()
	}
}
