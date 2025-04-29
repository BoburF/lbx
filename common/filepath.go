package common

import (
	"path/filepath"
)

type FilePathInfo struct {
	AbsolutePath string
	Dir          string
	Base         string
	Name         string
	Ext          string
}

func ParseFilePath(inputPath string) (*FilePathInfo, error) {
	absolutePath, err := filepath.Abs(inputPath)
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(absolutePath)
	base := filepath.Base(absolutePath)
	ext := filepath.Ext(absolutePath)
	name := base[:len(base)-len(ext)]

	return &FilePathInfo{
		AbsolutePath: absolutePath,
		Dir:          dir,
		Base:         base,
		Name:         name,
		Ext:          ext,
	}, nil
}
