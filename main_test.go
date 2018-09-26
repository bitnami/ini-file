package main

import (
	"flag"
	"os"
	"testing"

	ca "github.com/juamedgod/cliassert"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name           string
		wantErr        bool
		stdin          string
		expectedErr    interface{}
		expectedResult string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{}
			res := runTool(os.Args[0], args, "")
			if tt.wantErr {
				if res.Success() {
					t.Errorf("the command was expected to fail but succeeded")
				} else if tt.expectedErr != nil {
					res.AssertErrorMatch(t, tt.expectedErr)
				}
			} else {
				res.AssertSuccessMatch(t, tt.expectedResult)
			}
		})
	}
}

func TestMain(m *testing.M) {
	if os.Getenv("BE_TOOL") == "1" {
		main()
		os.Exit(0)
		return
	}
	flag.Parse()
	c := m.Run()
	os.Exit(c)
}

func runTool(bin string, args []string, stdin string) ca.CmdResult {
	cmd := ca.NewCommand()
	if stdin != "" {
		cmd.SetStdin(stdin)
	}
	os.Setenv("BE_TOOL", "1")
	defer os.Unsetenv("BE_TOOL")
	return cmd.Exec(bin, args...)
}

func RunTool(args ...string) ca.CmdResult {
	return runTool(os.Args[0], args, "")
}
