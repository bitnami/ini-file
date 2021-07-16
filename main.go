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

func main() {
	setCmd := NewINIFileSetCmd()
	getCmd := NewINIFileGetCmd()
	delCmd := NewINIFileDelCmd()

	parser := flags.NewParser(globalOpts, flags.HelpFlag|flags.PassDoubleDash)

	parser.AddCommand("set", "INI File Set", "Sets values in a INI file", setCmd)
	parser.AddCommand("get", "INI FILE Get", "Gets values from a INI file", getCmd)
	parser.AddCommand("del", "INI FILE Delete", "Deletes values from a INI file", delCmd)

	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
