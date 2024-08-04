package utils

import (
	"io"
	"os"
	"path/filepath"
)

func SaveImage(file io.Reader, filename string) (string, error) {
	uploadPath := "./uploads/"
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.MkdirAll(uploadPath, 0755)
	}

	filePath := filepath.Join(uploadPath, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return "", err
	}

	return filePath, nil
}
