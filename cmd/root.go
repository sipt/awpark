package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	aw "github.com/deanishe/awgo"
)

const CmdFlag = "cmd"

func init() {
	// Create a new Workflow using default settings.
	// Critical settings are provided by Alfred via environment variables,
	// so this *will* die in flames if not run in an Alfred-like environment.
	wf = aw.New()
	rootCmd.AddCommand(&cobra.Command{
		Use:   "store",
		Short: "workflow list",
		Run:   workflowList,
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "list kit",
		Run: func(cmd *cobra.Command, args []string) {
			wf.Run(func() {
				for _, action := range actionMap {
					item := action.ActionItem()
					if item != nil {
						item.Var(CmdFlag, action.Use())
					}
				}
				wf.SendFeedback()
			})
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use: "exec",
		Run: func(cmd *cobra.Command, args []string) {
			cmdFlag := wf.Config.Get(CmdFlag)
			if len(cmdFlag) == 0 {
				fmt.Printf("not found [%s] in vars", CmdFlag)
			}
			if action, ok := actionMap[cmdFlag]; !ok {
				fmt.Printf("not found [%s] in actions", CmdFlag)
			} else {
				switch action.GetRunMode() {
				case RunMode_Run:
					wf.Run(func() { action.Action(args) })
				case RunMode_Backgroud: // TODO
				case RunMode_None:
					action.Action(args)
				}
			}
		},
	})
}

var rootCmd = &cobra.Command{
	Use:   "awpark",
	Short: "alfred workflow park",
}

var wf *aw.Workflow

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var actionMap = make(map[string]Command)

func Register(cmd Command) {
	if _, ok := actionMap[cmd.Use()]; ok {
		panic(fmt.Sprintf("%s: is duplicate"))
	}
	actionMap[cmd.Use()] = cmd
}

type Command interface {
	Use() string
	ActionItem() *aw.Item
	Action(args []string)
	GetRunMode() RunMode
}

type RunModeRun struct{}

func (d *RunModeRun) GetRunMode() RunMode { return RunMode_Run }

type RunModeNone struct{}

func (d *RunModeNone) GetRunMode() RunMode { return RunMode_None }

type RunModeBackground struct{}

func (d *RunModeBackground) GetRunMode() RunMode { return RunMode_Backgroud }

type RunMode int

const (
	RunMode_None RunMode = iota
	RunMode_Run
	RunMode_Backgroud
)
