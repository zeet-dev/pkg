package utils

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func WriteLinesToFile(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return errors.Wrapf(err, "unable to create file: %s", filename)
	}
	defer func(file *os.File) {
		errClose := file.Close()
		if errClose != nil {
			err = errClose
		}
	}(file)

	bufferedWriter := bufio.NewWriter(file)
	for _, data := range lines {
		line := data
		// all lines must end with newline character
		if !strings.HasSuffix(line, "\n") {
			line = line + "\n"
		}
		_, err := bufferedWriter.WriteString(line)
		if err != nil {
			return errors.Wrap(err, "error writing lines to file")
		}
	}
	if err := bufferedWriter.Flush(); err != nil {
		return errors.Wrap(err, "error writing lines to file")
	}

	err = file.Sync()
	if err != nil {
		return errors.Wrap(err, "error syncing file")
	}

	return err
}

func CopyFile(dst, src string) error {
	if dst == src {
		return os.ErrInvalid
	}

	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()

	info, err := srcF.Stat()
	if err != nil {
		return err
	}

	dstF, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer dstF.Close()

	if _, err := io.Copy(dstF, srcF); err != nil {
		return err
	}
	return nil
}
