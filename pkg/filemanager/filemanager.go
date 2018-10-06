package filemanager

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/manuporto/distributedHTTPServer/pkg/cache"
	"github.com/manuporto/distributedHTTPServer/pkg/lockpool"
)

type FileManager struct {
	locks lockpool.LockPool
	cache cache.Cache
}

func NewFileManager(lockpoolSize uint, cacheSize uint) FileManager {
	lp := lockpool.NewLockPool(lockpoolSize)
	cch := cache.NewCache(cacheSize)
	return FileManager{lp, cch}
}

func (fm *FileManager) Save(filepath string, body []byte) error {
	l := fm.locks.GetLock(filepath)
	l.Lock()
	defer l.Unlock()
	err := saveFile(filepath, body)
	if err != nil {
		return err
	}
	fm.cache.Update(filepath, body)
	return err
}

func (fm *FileManager) Load(filepath string) ([]byte, error) {
	if body, ok := fm.cache.Get(filepath); ok {
		return body, nil
	}
	l := fm.locks.GetLock(filepath)
	l.RLock()
	defer l.RUnlock()
	body, err := loadFile(filepath)
	if err != nil {
		return nil, err
	}
	fm.cache.Insert(filepath, body)
	return body, nil
}

func (fm *FileManager) Update(filepath string, body []byte) error {
	l := fm.locks.GetLock(filepath)
	l.Lock()
	defer l.Unlock()
	err := updateFile(filepath, body)
	if err != nil {
		return err
	}
	fm.cache.Update(filepath, body)
	return nil
}

func (fm *FileManager) Delete(filepath string) error {
	l := fm.locks.GetLock(filepath)
	l.Lock()
	defer l.Unlock()
	err := os.Remove(filepath)
	if err != nil {
		return err
	}
	fm.cache.Delete(filepath)
	return nil
}

func saveFile(filepath string, body []byte) error {
	dir, file := path.Split(filepath)
	os.MkdirAll(dir, os.ModePerm)
	fd, err := os.OpenFile(path.Join(dir, file), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	_, err = fd.Write(body)
	fd.Close()
	return err
}

func loadFile(filename string) ([]byte, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func updateFile(filename string, body []byte) error {
	fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = fd.Write(body)
	return err
}
