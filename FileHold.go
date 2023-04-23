package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

type FileHold struct {
	Size int64
	AccessTime int64
	Hash string
}

func (fh *FileHold) Compare(path string, szOnly bool) bool {
	info, err := os.Stat(path)
	if err!=nil {return false}
	if fh.Size==info.Size() && szOnly {return true}
	md5Hash, err := fh.CalcHash(path)
	if err!=nil {return false}
	if md5Hash==fh.Hash {return true}

	return false
}

func (fh FileHold) CalcHash(path string) (string, error) {
	hash := md5.New()
	fp, err := os.Open(path)
	if err!=nil {return "", err}
	if _, err = io.Copy(hash, fp); err!=nil {return "", err}
	return fmt.Sprintf("%x",hash.Sum(nil)), nil
}

func NewFileHold(path string) *FileHold {
	info, err := os.Stat(path)
	if err!=nil {return nil}
	fh := FileHold{
		Size:       info.Size(),
		AccessTime: info.ModTime().UnixNano(),
	}
	fh.Hash, _ = fh.CalcHash(path)

	return &fh
}
