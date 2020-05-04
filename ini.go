package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/dchest/safefile"
	"github.com/go-ini/ini"
)

func init() {
	// Maintain the original file format
	ini.PrettyFormat = false
}

// iniLoad attempts to load the ini file.
func iniLoad(filename string) (*ini.File, error) {
	return ini.LoadSources(
		ini.LoadOptions{
			// Support mysql-style "boolean" values - a key wth no value.
			AllowBooleanKeys: true,
		},
		filename,
	)
}

// iniLoadOrEmpty attempts to load the ini file. If it does not exists,
// it will return an empty one
func iniLoadOrEmpty(filename string) (*ini.File, error) {
	f, err := iniLoad(filename)
	if err == nil {
		return f, nil
	}
	if os.IsNotExist(err) {
		return ini.Empty(), nil
	}
	return nil, err
}

// iniSave safely writes the ini file to the named file.
func iniSave(filename string, iniFile *ini.File) error {
	finfo, err := os.Stat(filename)
	if err != nil {
		return err
	}
	f, err := safefile.Create(filename, finfo.Mode())
	if err != nil {
		return err
	}
	defer f.Close()
	// safefile.Create doesn't seem to respect the permissions
	// so we need to recover the original permissions
	if err := os.Chmod(f.File.Name(), finfo.Mode()); err != nil {
		return err
	}
	sys := finfo.Sys().(*syscall.Stat_t)
	if err := os.Chown(f.File.Name(), int(sys.Uid), int(sys.Gid)); err != nil {
		return err
	}

	_, err = iniFile.WriteTo(f)
	if err != nil {
		return err
	}
	return f.Commit()
}

func iniFileGet(file string, s string, key string) (string, error) {
	iniFile, err := iniLoad(file)
	if err != nil {
		return "", err
	}
	section := iniFile.Section(s)
	if !section.HasKey(key) {
		return "", nil
	}
	k, err := section.GetKey(key)
	if err != nil {
		return "", err
	}
	return k.String(), nil
}

func iniFileSet(file string, s string, key string, value interface{}) error {
	iniFile, err := iniLoadOrEmpty(file)
	if err != nil {
		return err
	}
	section := iniFile.Section(s)
	switch v := value.(type) {
	case string:
		section.NewKey(key, v)
	case bool:
		section.NewBooleanKey(key)
	default:
		return fmt.Errorf("invalid key type %T", v)
	}

	return iniSave(file, iniFile)
}

func iniFileDel(file string, s string, key string) error {
	iniFile, err := iniLoad(file)
	if err != nil {
		return err
	}

	section := iniFile.Section(s)
	if !section.HasKey(key) {
		return nil
	}
	section.DeleteKey(key)
	return iniSave(file, iniFile)
}
