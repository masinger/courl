package util

import (
	"bytes"
	"io"
	"os"
	"strings"
)

func GetRawOrStream(expression string) (io.Reader, error) {
	if (strings.Index(expression, "@")) != 0 {
		return bytes.NewBufferString(expression), nil
	}
	remainder := expression[1:]
	if strings.EqualFold(remainder, "-") {
		return os.Stdin, nil
	}

	f, err := os.Open(remainder)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func CopyRawOrStream(expression string, writer io.Writer) (wasRef bool, err error) {
	if (strings.Index(expression, "@")) != 0 {
		_, err = writer.Write([]byte(expression))
		return false, err
	}
	remainder := expression[1:]
	if strings.EqualFold(remainder, "-") {
		_, err = io.Copy(writer, os.Stdin)
		return true, err
	}

	f, err := os.Open(remainder)
	if err != nil {
		return true, err
	}
	defer func() {
		if tempErr := f.Close(); tempErr != nil {
			err = tempErr
		}
	}()
	_, err = io.Copy(writer, f)
	return true, err
}
