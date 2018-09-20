package filemanager

import (
	"github.com/manuporto/distributedHTTPServer/pkg/cache"
	"github.com/manuporto/distributedHTTPServer/pkg/lockpool"
	"io/ioutil"
	"os"
	"path"
)

type FileManager struct {
	lp  lockpool.LockPool
	cch cache.Cache
}

func NewFileManager(lockpoolSize uint, cacheSize uint) FileManager {
	lp := lockpool.NewLockPool(lockpoolSize)
	cch := cache.NewCache(cacheSize)
	return FileManager{lp, cch}
}

func (fm FileManager) Save(filepath string, body []byte) error {
	l := fm.lp.GetLock(filepath)
	l.Lock()
	SaveFile(filepath, body)
	l.Unlock()
	return nil
}

func (fm FileManager) Load(filepath string) ([]byte, error) {
	if body, ok := fm.cch.Get(filepath); ok {
		return body, nil
	}
	l := fm.lp.GetLock(filepath)
	l.RLock()
	defer l.RUnlock()
	body, err := LoadFile(filepath)
	if err != nil {
		return nil, err
	}
	fm.cch.Update(filepath, body) // OJO
	return body, nil
}

func (fm FileManager) Update(filepath string, body []byte) error {
	l := fm.lp.GetLock(filepath)
	l.Lock()
	defer l.Unlock()
	err := UpdateFile(filepath, body)
	if err != nil {
		return err
	}
	fm.cch.Update(filepath, body)
	return nil
}

func (fm FileManager) Delete(filepath string) error {
	l := fm.lp.GetLock(filepath)
	l.Lock()
	defer l.Unlock()
	err := os.Remove(filepath)
	if err != nil {
		return err
	}
	fm.cch.Delete(filepath)
}

func SaveFile(filepath string, body []byte) error {
	dir, file := path.Split(filepath)
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
