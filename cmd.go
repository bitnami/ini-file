package main

import (
	"fmt"
	"io"
	"os"
)

// INIFileSetCmd defines a ini set command operation
type INIFileSetCmd struct {
	Section string `short:"s" long:"section" description:"INI Section" value-name:"SECTION" required:"yes"`
	Key     string `short:"k" long:"key" description:"INI Key to set" value-name:"KEY" required:"yes"`
	Value   string `short:"v" long:"value" description:"Value to store" value-name:"VALUE"`
	Boolean bool   `short:"b" long:"boolean" description:"Create a boolean key"`
	Args    struct {
		File string `positional-arg-name:"file"`
	} `positional-args:"yes" required:"yes"`
}

// NewINIFileSetCmd returns a new INIFileSetCmd
func NewINIFileSetCmd() *INIFileSetCmd {
	return &INIFileSetCmd{}
}

// Execute runs the ini set command
func (c *INIFileSetCmd) Execute(args []string) error {
	var v interface{}
	if c.Boolean {
		v = true
	} else {
		v = c.Value
	}
	return iniFileSet(c.Args.File, c.Section, c.Key, v)
}

// INIFileGetCmd defines a ini get command operation
type INIFileGetCmd struct {
	Section string `short:"s" long:"section" description:"INI Section" value-name:"SECTION" required:"yes"`
	Key     string `short:"k" long:"key" description:"INI Key to get" value-name:"KEY" required:"yes"`
	Args    struct {
		File string `positional-arg-name:"file"`
	} `positional-args:"yes" required:"yes"`
	OutWriter io.Writer
}

// NewINIFileGetCmd returns a new INIFileGetCmd
func NewINIFileGetCmd() *INIFileGetCmd {
	return &INIFileGetCmd{OutWriter: os.Stdout}
}

// Execute runs the ini get command
func (c *INIFileGetCmd) Execute(args []string) error {
	v, err := iniFileGet(c.Args.File, c.Section, c.Key)
	if err != nil {
		return err
	}
	fmt.Fprint(c.OutWriter, v)
	return nil
}

// INIFileDelCmd defines a ini del command operation
type INIFileDelCmd struct {
	Section string `short:"s" long:"section" description:"INI Section" value-name:"SECTION" required:"yes"`
	Key     string `short:"k" long:"key" description:"INI Key to delete" value-name:"KEY" required:"yes"`
	Args    struct {
		File string `positional-arg-name:"file"`
	} `positional-args:"yes" required:"yes"`
}

// NewINIFileDelCmd returns a new INIFileDelCmd
func NewINIFileDelCmd() *INIFileDelCmd {
	return &INIFileDelCmd{}
}

// Execute runs the ini del command
func (c *INIFileDelCmd) Execute(args []string) error {
	return iniFileDel(c.Args.File, c.Section, c.Key)
}
