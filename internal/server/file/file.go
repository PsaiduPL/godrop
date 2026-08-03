package file

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"strings"
	"time"
)

type ShareFile struct {
	Path  string
	IsDir bool
}

type FsProxy struct {
	fs fs.FS
}

func (fsProxy *FsProxy) Open(name string) (fs.File, error) {
	if strings.Contains(name, "DS_Store") {
		return &emptyFile{name}, nil
	}
	return fsProxy.fs.Open(name)
}

type emptyFile struct {
	name string
}

func (e *emptyFile) Stat() (fs.FileInfo, error) { return e, nil }
func (e *emptyFile) Read([]byte) (int, error)   { return 0, io.EOF } // Zawsze natychmiastowy EOF
func (e *emptyFile) Close() error               { return nil }

// Metody dla interface'u fs.FileInfo:
func (e *emptyFile) Name() string       { return e.name }
func (e *emptyFile) Size() int64        { return 0 } // Rozmiar 0 B
func (e *emptyFile) Mode() fs.FileMode  { return 0444 }
func (e *emptyFile) ModTime() time.Time { return time.Now() }
func (e *emptyFile) IsDir() bool        { return false }
func (e *emptyFile) Sys() any           { return nil }

func (file *ShareFile) WriteContent(writer io.Writer) error {

	if file.IsDir {
		return handleFolder(file.Path, writer)
	}
	return handleFile(file.Path, writer)

}
func handleFile(path string, writer io.Writer) error {
	slog.Info("Received info for writing content")
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.WriteTo(writer)
	slog.Info("Content written")
	if err != nil {
		slog.Error("Error while writing content", "error", err)
		return err
	}
	return nil
}

func handleFolder(pathToZip string, writer io.Writer) error {
	file, err := os.CreateTemp(".", "temp.*.zip")

	defer os.Remove(file.Name())
	if err != nil {
		return fmt.Errorf("Error while creating zip file %w", err)
	}

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	err = zipWriter.AddFS(&FsProxy{fs: os.DirFS(pathToZip)})
	if err != nil {
		return fmt.Errorf("Error while writing data to zip %w", err)
	}

	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("error while closing zip writer: %w", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("error while seeking to start of zip file: %w", err)
	}
	_, err = file.WriteTo(writer)
	if err != nil {
		return fmt.Errorf("Error while writing data to http body %w", err)
	}

	return nil
}
