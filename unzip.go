package IgorLib

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func popS(s *[]string) string {
	ret := (*s)[0]
	if len(*s) > 0 {
		*s = (*s)[1:]
	} else {
		*s = make([]string, 0)
	}
	return ret
}

func CreateDirectories(source string) error {
	currentPath, delim := "", string(os.PathSeparator)
	parts := strings.Split(source, delim)
	currentPath += popS(&parts) + delim
	for len(parts) > 0 {
		currentPath = fmt.Sprintf("%s%s%s", currentPath, delim, popS(&parts))
		_, err := os.Stat(currentPath)
		if os.IsNotExist(err) {
			err = os.Mkdir(currentPath, os.ModePerm)
			if IsErr(err) {
				return err
			}
		} else if err != nil {
			return err
		}
	}
	return nil
}

func UnzipSource(source, destination string) error {
	reader, err := zip.OpenReader(source)
	if IsErr(err) {
		return err
	}
	defer reader.Close()

	destination, err = filepath.Abs(destination)
	if IsErr(err) {
		return err
	}

	//err = CreateDirectories(destination)
	//if IsErr(err) {
	//	return err
	//}

	for _, f := range reader.File {
		err = UnzipFile(f, destination)
		if IsErr(err) {
			return err
		}

	}
	return nil
}

func UnzipFile(f *zip.File, destination string) error {
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("Invalid file path: %s\n", filePath)
	}

	if f.FileInfo().IsDir() {
		err := os.MkdirAll(filePath, os.ModePerm)
		if IsErr(err) {
			return err
		}
		return nil
	}

	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if IsErr(err) {
		return err
	}

	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if IsErr(err) {
		return err
	}
	defer destinationFile.Close()

	zippedFile, err := f.Open()
	if IsErr(err) {
		return err
	}
	defer zippedFile.Close()

	_, err = io.Copy(destinationFile, zippedFile)
	if IsErr(err) {
		return err
	}

	return nil
}
