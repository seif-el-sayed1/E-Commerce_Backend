package constants

import "github.com/fatih/color"

func init() {
	color.NoColor = false
}

var (
	Success = color.New(color.FgGreen, color.Bold).SprintFunc()
	Error   = color.New(color.FgRed, color.Bold).SprintFunc()
	Warning = color.New(color.FgYellow, color.Bold).SprintFunc()
	Info    = color.New(color.FgCyan, color.Bold).SprintFunc()
)
