package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Extract(archiveFile, destDir string) error {
	fileExt := strings.ToLower(filepath.Ext(archiveFile))

	switch fileExt {
	case ".zip":
		return UnZip(archiveFile, destDir)
	case ".tar":
		fallthrough
	case ".tar.gz":
		fallthrough
	case ".tgz":
		return UnTar(archiveFile, destDir)
	default:
		return fmt.Errorf("unsupported archive format: %s", fileExt)
	}
}

func UnZip(archiveFile, destDir string) error {
	r, err := zip.OpenReader(archiveFile)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(destDir, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}

	return nil
}

func UnTar(archiveFile, destDir string) error {
	file, err := os.Open(archiveFile)
	if err != nil {
		return err
	}
	defer file.Close()

	var r io.Reader = file

	if strings.HasSuffix(archiveFile, ".gz") || strings.HasSuffix(archiveFile, ".tgz") {
		gr, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer gr.Close()
		r = gr
	}

	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, header.FileInfo().Mode()); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, header.FileInfo().Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tr); err != nil {
				return err
			}
		}
	}

	return nil
}
