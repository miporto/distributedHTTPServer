package filemanager

import (
	"io/ioutil"
	"os"
)

func FileExists(filename string) bool {
	return true
}

func SaveFile(filename string, body []byte) error {
	return ioutil.WriteFile(filename, body, 0600)
}

func LoadFile(filename string) ([]byte, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func UpdateFile(filename string, body []byte) error {
	return SaveFile(filename, body)
}

func DeleteFile(filename string) error {
	return os.Remove(filename)
}
