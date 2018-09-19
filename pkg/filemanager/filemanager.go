package filemanager

import (
	"io/ioutil"
	"os"
	"path"
)

func SaveFile(filepath string, body []byte) error {
	dir, file := path.Split(filepath)
	file = file + ".json"
	os.MkdirAll(dir, os.ModePerm)
	fd, err := os.OpenFile(path.Join(dir, file), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = fd.Write(body)
	return err
}

func LoadFile(filename string) ([]byte, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func UpdateFile(filename string, body []byte) error {
	fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = fd.Write(body)
	return err
}

func DeleteFile(filename string) error {
	return os.Remove(filename)
}
