package file

import (
	"io"
	"log/slog"
	"os"
)

type ShareFile struct {
	Path  string
	IsDir bool
}

func (file *ShareFile) WriteContent(writer io.Writer) error {
	if file.IsDir {

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

func handleFolder(path string, writer io.Writer) {

}
