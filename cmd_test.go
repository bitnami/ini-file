package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	tu "github.com/bitnami/gonit/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type iniTestValue struct {
	globalOpts *Options
	section    string
	key        string
	value      string
	isBoolean  bool
}
type iniSetTest struct {
	name          string
	values        []iniTestValue
	wantErr       bool
	initialText   string
	expectedText  string
	createIniFile bool
	expectedErr   interface{}
}

type iniGetTest struct {
	name string
	iniTestValue
	wantErr       bool
	initialText   string
	expectedText  string
	createIniFile bool
	expectedErr   interface{}
}

type iniDelTest struct {
	name          string
	values        []iniTestValue
	wantErr       bool
	initialText   string
	expectedText  string
	createIniFile bool
	expectedErr   interface{}
}

var delTests = []iniDelTest{
	{
		name: "Deletes key non existent",
		values: []iniTestValue{
			{
				section: "general", key: "mykey",
			},
		},
		createIniFile: true,
		expectedText:  ``,
	},
	{
		name:        "Deletes boolean value",
		initialText: "[general]\nboolkey\nkey1=val1\n",
		values: []iniTestValue{
			{
				section: "general", key: "boolkey",
			},
		},
		expectedText: `\[general\]\nkey1=val1\n\s*$`,
	},
	{
		name: "Deletes key in file containing semicolon values in regular mode",
		values: []iniTestValue{
			{
				section: "general", key: "key2",
			},
		},
		initialText:  "[general]\nkey1=`my ; value`\nkey2=val2",
		expectedText: "\\[general\\]\nkey1=`my ; value`\n\\s*$",
	},
	{
		name: "Deletes key in file containing semicolon values in regular mode (2)",
		values: []iniTestValue{
			{
				section: "general", key: "key2",
			},
		},
		initialText:  "[general]\nkey1=my ; my comment\nkey2=val2",
		expectedText: "\\[general\\]\n; my comment\nkey1=my\n\\s*$",
	},
	{
		name: "Deletes key in file containing semicolon values in ignore-inline-comments mode",
		values: []iniTestValue{
			{
				section: "general", key: "key2", globalOpts: &Options{IgnoreInlineComments: true},
			},
		},
		initialText:  "[general]\nkey1=`my ; value`\nkey2=val2",
		expectedText: "\\[general\\]\nkey1=my ; value\n\\s*$",
	},
	{
		name: "Deletes key in file containing semicolon values in ignore-inline-comments mode (2)",
		values: []iniTestValue{
			{
				section: "general", key: "key2", globalOpts: &Options{IgnoreInlineComments: true},
			},
		},
		initialText:  "[general]\nkey1=my ; my comment\nkey2=val2",
		expectedText: "\\[general\\]\nkey1=my ; my comment\n\\s*$",
	},
	{
		name:        "Deletes regular value",
		initialText: "[general]\nkey1=val1\nkey2=val2\n",
		values: []iniTestValue{
			{
				section: "general", key: "key1",
			},
		},
		expectedText: `\[general\]\nkey2=val2\n\s*$`,
	},
	{
		name:        "Fails if ini file does not exists",
		values:      []iniTestValue{{section: "general", key: "key1"}},
		expectedErr: "no such file or directory",
	},
	{
		name:          "Preserve comments",
		createIniFile: true,
		initialText:   "# this is a comment\n[general]\n# key 1 sample\nkey1=value1\n# mykey comment\nmykey=myvalue",
		values: []iniTestValue{
			{section: "general", key: "key1"},
		},
		expectedText: `^# this is a comment\n\[general\]\n# mykey comment\nmykey=myvalue\n\s*$`,
	},
}
var getTests = []iniGetTest{
	{
		name:          "Get non-existent",
		createIniFile: true,
		iniTestValue: iniTestValue{
			section: "general", key: "mykey", value: "",
		},
	},
	{
		name:        "Get regular key",
		initialText: "[general]\nmykey=myvalue\n",
		iniTestValue: iniTestValue{
			section: "general", key: "mykey", value: "myvalue",
		},
	},
	{
		name:        "Get value with semicolon in regular mode",
		initialText: "[general]\nmykey=`my ; value`\n",
		iniTestValue: iniTestValue{
			section: "general", key: "mykey", value: "my ; value",
		},
	},
	{
		name:        "Get value with semicolon in regular mode (2)",
		initialText: "[general]\nmykey=my ; this is a comment\n",
		iniTestValue: iniTestValue{
			section: "general", key: "mykey", value: "my",
		},
	},
	{
		name:        "Get value with semicolon in ignore-inline-comments mode",
		initialText: "[general]\nmykey=`my ; value`\n",
		iniTestValue: iniTestValue{
			globalOpts: &Options{IgnoreInlineComments: true},
			section:    "general", key: "mykey", value: "my ; value",
		},
	},
	{
		name:        "Get value with semicolon in ignore-inline-comments mode (2)",
		initialText: "[general]\nmykey=my ; this is a comment\n",
		iniTestValue: iniTestValue{
			globalOpts: &Options{IgnoreInlineComments: true},
			section:    "general", key: "mykey", value: "my ; this is a comment",
		},
	},
	{
		name:        "Get present boolean key",
		initialText: "[general]\nmykey\n",
		iniTestValue: iniTestValue{
			section: "general", key: "mykey", value: "true",
		},
	},
	{
		name:         "Fails if ini file does not exists",
		iniTestValue: iniTestValue{section: "general", key: "key1"},
		expectedErr:  "no such file or directory",
	},
	{
		name:          "Get from malformed file",
		createIniFile: true,
		initialText:   "this is not a\nINI\nfile\nmykey\n[general]\nmykey=myvalue",
		iniTestValue: iniTestValue{
			section: "general", key: "mykey", value: "myvalue",
		},
	},
}
var setTests = []iniSetTest{
	{
		name: "Sets regular key non existent",
		values: []iniTestValue{
			{
				section: "general", key: "mykey", value: "myvalue",
			},
		},
		expectedText: `mykey=myvalue\n`,
	},
	{
		name: "Sets value with semicolon in regular mode",
		values: []iniTestValue{
			{
				section: "general", key: "mykey", value: "my ; value",
			},
		},
		expectedText: "mykey=`my ; value`\n",
	},
	{
		name: "Sets value with semicolon in ignore-inline-comments mode",
		values: []iniTestValue{
			{
				globalOpts: &Options{IgnoreInlineComments: true},
				section:    "general", key: "mykey", value: "my ; value",
			},
		},
		expectedText: "mykey=my ; value\n",
	},
	{
		name: "Sets boolean value",
		values: []iniTestValue{
			{
				section: "testbool", key: "mykey", isBoolean: true,
			},
		},
		expectedText: `\[testbool\]\nmykey\n\s*$`,
	},
	{
		name:        "Override with boolean value",
		initialText: `\[testbool\]\nmykey=value1\n\s*$`,
		values: []iniTestValue{
			{
				section: "testbool", key: "mykey", isBoolean: true,
			},
		},
		expectedText: `\[testbool\]\nmykey\n\s*$`,
	},
	{
		name:        "Override boolean value with regular one",
		initialText: `\[testbool\]\nmykey\n\s*$`,
		values: []iniTestValue{
			{
				section: "testbool", key: "mykey", value: "myvalue",
			},
		},
		expectedText: `\[testbool\]\nmykey=myvalue\n\s*$`,
	},
	{
		name: "Set multiple keys",
		values: []iniTestValue{
			{section: "general", key: "key1", value: "value1"},
			{section: "general", key: "key2", value: "value2"},
			{section: "general", key: "key3", value: "value3"},
			{section: "general", key: "key4", isBoolean: true},
		},
		expectedText: `^\[general\]\nkey1=value1\nkey2=value2\nkey3=value3\nkey4\n\s*$`,
	},
	{
		name: "Sets regular keys in existing file",
		values: []iniTestValue{
			{section: "general", key: "mykey", value: "myvalue"},
			{section: "general", key: "key2", value: "newvalue2"},
			{section: "newsection", key: "key5", value: "value5"},
		},
		initialText: `
[general]
key1=value1
key2=value2
key3=value3
[newsection]
key4=value4
		`,
		expectedText: `^\[general\]\nkey1=value1\nkey2=newvalue2\nkey3=value3\nmykey=myvalue\n\s+` +
			`\[newsection\]\nkey4=value4\nkey5=value5\n.*`,
	},
	{
		name:          "Preserve comments",
		createIniFile: true,
		initialText:   "# this is a comment\n[general]\n# key 1 sample\nkey1=value1",
		values: []iniTestValue{
			{section: "general", key: "mykey", value: "myvalue"},
		},
		expectedText: `^# this is a comment\n\[general\]\n# key 1 sample\nkey1=value1\nmykey=myvalue\n\s*$`,
	},
}

func testFile(t *testing.T, path string, fn func(t *testing.T, contents string) bool, msgAndArgs ...interface{}) bool {
	if !assert.FileExists(t, path) {
		return false
	}
	data, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	return fn(t, string(data))
}
func AssertFileContains(t *testing.T, path string, expected interface{}, msgAndArgs ...interface{}) bool {
	return testFile(t, path, func(t *testing.T, contents string) bool {
		return assert.Regexp(t, expected, contents, msgAndArgs...)
	})
}
func AssertFileDoesNotContain(t *testing.T, path string, expected interface{}, msgAndArgs ...interface{}) bool {
	return testFile(t, path, func(t *testing.T, contents string) bool {
		return assert.NotRegexp(t, expected, contents, msgAndArgs...)
	})
}
func TestINIFileSetCmd_Execute(t *testing.T) {
	type testFn func(file, section, key, value string, isBoolean bool, opts *Options) error
	var testViaCommand = func(file, section, key, value string, isBoolean bool, opts *Options) error {
		cmd := NewINIFileSetCmd()
		cmd.Section = section
		cmd.Key = key
		cmd.Value = value
		cmd.Boolean = isBoolean
		cmd.Args.File = file
		if opts != nil {
			globalOpts = opts
		} else {
			globalOpts = &Options{}
		}
		return cmd.Execute([]string{})
	}
	var testViaCli = func(file, section, key, value string, isBoolean bool, opts *Options) error {
		args := []string{"set", "-k", key, "-s", section, "-v", value}
		if isBoolean {
			args = append(args, "-b")
		}
		args = append(args, file)
		if opts != nil && opts.IgnoreInlineComments {
			args = append([]string{"--ignore-inline-comments"}, args...)
		}
		res := RunTool(args...)
		if !res.Success() {
			return fmt.Errorf("%s", res.Stderr())
		}
		return nil
	}

	for _, tt := range setTests {
		for id, fn := range map[string]testFn{
			"command": testViaCommand,
			"cli":     testViaCli,
		} {
			var err error

			file := ""
			sb := tu.NewSandbox()
			defer sb.Cleanup()
			if tt.initialText != "" || tt.createIniFile {
				file, err = sb.Write("my.ini", tt.initialText)
				require.NoError(t, err)
			} else {
				file = sb.Normalize("my.ini")
			}
			testTitle := fmt.Sprintf("%s (via %s)", tt.name, id)
			t.Run(testTitle, func(t *testing.T) {
				for _, v := range tt.values {
					err = fn(file, v.section, v.key, v.value, v.isBoolean, v.globalOpts)
					if err != nil {
						break
					}
				}
				if tt.expectedErr != nil {
					if err == nil {
						t.Errorf("the command was expected to fail but succeeded")
					} else {
						assert.Regexp(t, tt.expectedErr, err)
					}
				} else {
					require.NoError(t, err)
					AssertFileContains(t, file, tt.expectedText)
				}
			})
		}
	}
}

func TestINIFileGetCmd_Execute(t *testing.T) {
	type testFn func(file, section, key string, opts *Options) (string, error)
	var testViaCommand = func(file, section, key string, opts *Options) (string, error) {
		b := &bytes.Buffer{}
		cmd := NewINIFileGetCmd()
		cmd.Section = section
		cmd.Key = key
		cmd.Args.File = file
		cmd.OutWriter = b
		if opts != nil {
			globalOpts = opts
		} else {
			globalOpts = &Options{}
		}
		err := cmd.Execute([]string{})
		stdout := b.String()
		return stdout, err
	}
	var testViaCli = func(file, section, key string, opts *Options) (string, error) {
		args := []string{"get", "-k", key, "-s", section, file}
		if opts != nil && opts.IgnoreInlineComments {
			args = append([]string{"--ignore-inline-comments"}, args...)
		}
		res := RunTool(args...)
		stdout := res.Stdout()
		var err error
		if !res.Success() {
			err = fmt.Errorf("%s", res.Stderr())
		}
		return stdout, err
	}

	for _, tt := range getTests {
		for id, fn := range map[string]testFn{
			"command": testViaCommand,
			"cli":     testViaCli,
		} {
			var err error
			file := ""
			sb := tu.NewSandbox()
			defer sb.Cleanup()
			if tt.initialText != "" || tt.createIniFile {
				file, err = sb.Write("my.ini", tt.initialText)
				require.NoError(t, err)
			} else {
				file = sb.Normalize("my.ini")
			}
			testTitle := fmt.Sprintf("%s (via %s)", tt.name, id)

			t.Run(testTitle, func(t *testing.T) {

				stdout, err := fn(file, tt.section, tt.key, tt.globalOpts)

				if tt.expectedErr != nil {
					if err == nil {
						t.Errorf("the command was expected to fail but succeeded")
					} else {
						assert.Regexp(t, tt.expectedErr, err)
					}
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.value, stdout)
				}
			})
		}
	}
}

func TestINIFileDelCmd_Execute(t *testing.T) {
	type testFn func(file, section, key string, opts *Options) error
	var testViaCommand = func(file, section, key string, opts *Options) error {
		cmd := NewINIFileDelCmd()
		cmd.Section = section
		cmd.Key = key
		cmd.Args.File = file
		if opts != nil {
			globalOpts = opts
		} else {
			globalOpts = &Options{}
		}
		return cmd.Execute([]string{})
	}
	var testViaCli = func(file, section, key string, opts *Options) error {
		args := []string{"del", "-k", key, "-s", section, file}
		if opts != nil && opts.IgnoreInlineComments {
			args = append([]string{"--ignore-inline-comments"}, args...)
		}
		res := RunTool(args...)
		if !res.Success() {
			return fmt.Errorf("%s", res.Stderr())
		}
		return nil
	}
	for _, tt := range delTests {
		for id, fn := range map[string]testFn{
			"command": testViaCommand,
			"cli":     testViaCli,
		} {
			var err error
			file := ""
			sb := tu.NewSandbox()
			defer sb.Cleanup()
			if tt.initialText != "" || tt.createIniFile {
				file, err = sb.Write("my.ini", tt.initialText)
				require.NoError(t, err)
			} else {
				file = sb.Normalize("my.ini")
			}
			testTitle := fmt.Sprintf("%s (via %s)", tt.name, id)

			t.Run(testTitle, func(t *testing.T) {
				for _, v := range tt.values {

					err = fn(file, v.section, v.key, v.globalOpts)
					if err != nil {
						break
					}
				}
				if tt.expectedErr != nil {
					if err == nil {
						t.Errorf("the command was expected to fail but succeeded")
					} else {
						assert.Regexp(t, tt.expectedErr, err)
					}
				} else {
					require.NoError(t, err)
					AssertFileContains(t, file, tt.expectedText)
				}
			})
		}
	}
}
