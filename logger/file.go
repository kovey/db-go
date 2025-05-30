package logger

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type File struct {
	file *os.File
	path string
	date string
}

func NewFile(path string) *File {
	if !strings.HasSuffix(path, "/db") {
		path = fmt.Sprintf("%s/db", path)
	}

	err := os.MkdirAll(path, 0755)
	if err != nil && !os.IsExist(err) {
		return nil
	}

	file := &File{path: path, date: time.Now().Format(time.DateOnly)}
	return file
}

func (f *File) Open() error {
	filePath := fmt.Sprintf("%s/%s.log", f.path, f.date)
	fi, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if f.file != nil {
		f.file.Close()
	}

	f.file = fi
	return nil
}

func (f *File) Check(now time.Time) error {
	nowDate := time.Now().Format(time.DateOnly)
	if f.date == nowDate {
		return nil
	}

	f.date = nowDate
	return f.Open()
}

func (f *File) Write(data []byte) (int, error) {
	if f.file == nil {
		return 0, nil
	}

	return f.file.Write(data)
}

func (f *File) Close() {
	if f.file == nil {
		return
	}

	f.file.Close()
}
