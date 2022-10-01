package utils

import "os"

func CopyFile(src string, dst string) error {
	open, err := os.Open(src)
	if err != nil {
		return err
	}

	defer open.Close()

	newFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer newFile.Close()
	newFile.ReadFrom(open)

	return nil
}
