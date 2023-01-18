package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
)

// Options defines options supported by all subcommands
type Options struct {
	IgnoreInlineComments bool `long:"ignore-inline-comments" description:"Ignore inline comments"`
}

var globalOpts = &Options{}

var (
	version   = "1.4.5"
	buildDate = ""
	commit    = ""
)

func versionText() string {
	msg := fmt.Sprintf("%-12s %s", "Version:", version)
	if buildDate != "" {
		msg += fmt.Sprintf("\n%-12s %s", "Built on:", buildDate)
	}
	if commit != "" {
		msg += fmt.Sprintf("\n%-12s %s", "Git Commit:", commit)
	}
	return msg
}

func main() {
	setCmd := NewINIFileSetCmd()
	getCmd := NewINIFileGetCmd()
	delCmd := NewINIFileDelCmd()

	parser := flags.NewParser(globalOpts, flags.HelpFlag|flags.PassDoubleDash)

	parser.LongDescription = versionText()

	parser.AddCommand("set", "INI File Set", "Sets values in a INI file", setCmd)
	parser.AddCommand("get", "INI FILE Get", "Gets values from a INI file", getCmd)
	parser.AddCommand("del", "INI FILE Delete", "Deletes values from a INI file", delCmd)

	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
