package io

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
)

//ReadLines ...
func ReadLines(filePath string) (lines []string, err error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	content := string(bytes)
	lines = strings.Split(content, "\n")

	return
}

//FileOrDirExists ...
func FileOrDirExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

// CreateFile ...
func CreateFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}

// ReadFile ...
func ReadFile(filePath string) (result string, err error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	return string(bytes), nil
}

//WriteFile ...
func WriteFile(outPath string, lines ...string) error {
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	out := bufio.NewWriter(outFile)
	for _, l := range lines {
		if _, err := fmt.Fprintln(out, l); err != nil {
			return err
		}
	}
	return out.Flush()
}

//CreateDirIfNotExist ...
func CreateDirIfNotExist(dirPath string) (err error) {
	if FileOrDirExists(dirPath) {
		return
	}

	err = os.MkdirAll(dirPath, os.ModePerm)
	return
}

// CreateFileIfNotExist ...
func CreateFileIfNotExist(filePath string) (err error) {
	if err = CreateDirIfNotExist(path.Dir(filePath)); err != nil {
		return
	}

	if !FileOrDirExists(filePath) {
		file, err := os.Create(filePath)
		if err != nil {
			glog.Error(err)
			return err
		}
		defer file.Close()
	}
	return
}

// GetFilesInDictionary ...
func GetFilesInDictionary(dirPath string) (files []string, err error) {
	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		glog.Error(err)
		return
	}
	for _, f := range fs {
		if !f.IsDir() {
			files = append(files, filepath.Join(dirPath, f.Name()))
		}
	}
	return files, nil
}

// GetDirsInDictionary ...
func GetDirsInDictionary(dirPath string) (dirs []string, err error) {
	ds, err := ioutil.ReadDir(dirPath)
	if err != nil {
		glog.Error(err)
		return
	}
	for _, d := range ds {
		if d.IsDir() {
			dirs = append(dirs, filepath.Join(dirPath, d.Name()))
		}
	}
	return dirs, nil
}

// GetFilesWithSuffix ...
func GetFilesWithSuffix(dirPth, suffix string) (files []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	for _, fi := range dir {
		if fi.IsDir() {
			GetFilesWithSuffix(filepath.Join(dirPth, fi.Name()), suffix)
		} else {
			ok := strings.HasSuffix(fi.Name(), suffix)
			if ok {
				files = append(files, filepath.Join(dirPth, fi.Name()))
			}
		}
	}

	return files, nil
}

// AppendToFile ...
func AppendToFile(file, text string) (err error) {
	if err = CreateFileIfNotExist(file); err != nil {
		glog.Error(err)
		return
	}

	w, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		glog.Error(err)
		return
	}
	defer w.Close()
	if _, err = w.WriteString(text); err != nil {
		glog.Error(err)
		return
	}
	return
}

// CountFilesInDir ...
func CountFilesInDir(path string) (count int, err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		glog.Error(err)
		return
	}
	count = len(files)
	return
}

// DeCompress ...
func DeCompress(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := filepath.Join(dest, hdr.Name)
		file, err := CreateFile(filename)
		if err != nil {
			return err
		}
		io.Copy(file, tr)
	}
	return nil
}

// IsGz ...
func IsGz(fileName string) bool {
	if filepath.Ext(fileName) == ".gz" {
		return true
	}
	return false
}

// IsZip ...
func IsZip(fileName string) bool {
	if filepath.Ext(fileName) == ".zip" {
		return true
	}
	return false
}

// CopyFile ...
func CopyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

// FileExist returns true if the specified normal file exists
func FileExist(filename string) (ok bool) {
	info, err := os.Stat(filename)
	if err == nil {
		if ^os.ModePerm&info.Mode() == 0 {
			ok = true
		}
	}
	return ok
}
