package cmd

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

const diskFolderName = "dddplayer"

type DiskWriter struct {
	content string
	name    string
	root    string
}

func NewDiskWriter(content, filename, mainPath string) (*DiskWriter, error) {
	rootDir, err := createDiskFolderIfNotExist(mainPath)
	if err != nil {
		return nil, err
	}

	dw := &DiskWriter{
		content: content,
		name:    filename,
		root:    rootDir,
	}

	return dw, nil
}

func createDiskFolderIfNotExist(mainPath string) (string, error) {
	projectRootDir, err := findProjectRootDir(mainPath)
	if err != nil {
		return "", err
	}

	rootDir := path.Join(projectRootDir, diskFolderName)
	err = os.MkdirAll(rootDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	return rootDir, nil
}

func (dw *DiskWriter) Write() error {
	if !dw.isUpdated() {
		return nil
	}

	err := dw.write(dw.content, dw.filename())
	if err != nil {
		return err
	}

	err = dw.write(dw.Hash(), dw.hashName())
	if err != nil {
		return err
	}
	return nil
}

func (dw *DiskWriter) Hash() string {
	return sha1Sum(dw.content)
}

func (dw *DiskWriter) isUpdated() bool {
	if !dw.fileExists(dw.hashName()) {
		return true
	}

	newCodeHash := dw.Hash()
	oldCodeHash, err := os.ReadFile(dw.hashName())
	if err != nil {
		return true
	}

	if newCodeHash != string(oldCodeHash) {
		return true
	}

	return false
}

func (dw *DiskWriter) fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func (dw *DiskWriter) write(content, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(content))
	if err != nil {
		return err
	}

	return nil
}

func (dw *DiskWriter) filename() string {
	return path.Join(dw.root, fmt.Sprintf("%s.dot", dw.name))
}

func (dw *DiskWriter) hashName() string {
	return path.Join(dw.root, fmt.Sprintf("%s.hash", dw.name))
}

func findProjectRootDir(startDir string) (string, error) {
	dir := startDir
	for {
		goModPath := filepath.Join(dir, "go.mod")
		_, err := os.Stat(goModPath)
		if err == nil {
			return dir, nil
		}

		if dir == "/" || dir == "." {
			return dir, fmt.Errorf("cannot find go.mod file")
		}

		// 向上级目录移动一级
		dir = filepath.Dir(dir)
	}
}

func sha1Sum(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	b := h.Sum(nil)
	return fmt.Sprintf("%x", b)
}

func writeToDisk(raw, filename, mainPkg string) error {
	dw, err := NewDiskWriter(raw, filename, mainPkg)
	if err != nil {
		return err
	}
	if err := dw.Write(); err != nil {
		return err
	}
	return nil
}
