package file

import (
	"archive/zip"
	"os"
	"testing"
)

func Test_(t *testing.T) {
	f, _ := os.Create("interesting.zip")
	defer f.Close()
	writer := zip.NewWriter(f)
	defer writer.Close()
	err := writer.AddFS(os.DirFS("../api"))
	if err != nil {
		panic(err)
	}
}
