package file

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)

func GetSize(f multipart.File) (int, error) {
	content, err := ioutil.ReadAll(f)

	return len(content), err
}

func GetExt(fileName string) string {
	return path.Ext(fileName)
}

func CheckNotExist(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

func IsNotExistMkDir(src string) error {
	dir, _ := path.Split(src)
	if notExist := CheckNotExist(dir); notExist == true {
		if err := MkDir(dir); err != nil {
			return err
		}
	}

	return nil
}

func MkDir(dir string) error {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func Open(src string) (*os.File, error) {
	if err := IsNotExistMkDir(src); err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src:%v", err)
	}

	p := CheckPermission(src)
	if p == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	f, err := os.OpenFile(src, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return f, nil
}