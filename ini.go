package main

import (
	"fmt"
	"os"

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

// iniSave writes the ini file to the named file.
func iniSave(filename string, iniFile *ini.File) error {
	// The third argument, perm, is ignored when the file doesn't exist
	// So we can safely set it to '0644', it won't modify the existing permissions
	// if the file exists.
	f, err := os.OpenFile(filename, os.O_SYNC|os.O_RDWR|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		return err
	}
	// Clear file content
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = iniFile.WriteTo(f)
	if err != nil {
		return err
	}
	return f.Close()
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
