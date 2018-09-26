package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
)

func main() {
	setCmd := NewINIFileSetCmd()
	getCmd := NewINIFileGetCmd()
	delCmd := NewINIFileDelCmd()

	parser := flags.NewParser(nil, flags.HelpFlag|flags.PassDoubleDash)

	parser.AddCommand("set", "INI File Set", "Sets values in a INI file", setCmd)
	parser.AddCommand("get", "INI FILE Get", "Gets values from a INI file", getCmd)
	parser.AddCommand("del", "INI FILE Delete", "Deletes values from a INI file", delCmd)

	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
