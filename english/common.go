package main

import (
	"fmt"
	"os"
)

//判断文件是否存在
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//判断目录是否存在
func pathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//创建目录，如果没有就创建。
func mkdir(dir string) (err error) {
	exist, err := pathExist(dir)
	if err != nil {
		return fmt.Errorf("get dir error!: %s", err)
	}
	if !exist {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("mkdir failed![%v]\n", err)
		}
	}
	return nil
}
