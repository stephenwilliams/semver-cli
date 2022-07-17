package terminal

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
)

var (
	ErrStdInIsNotPiped = errors.New("stdin is not piped")
)

func IsPiped(in FileReader) bool {
	fi, _ := in.Stat()

	return (fi.Mode() & os.ModeCharDevice) == 0
}

func ReadPipedInput(in FileReader) ([]byte, error) {
	if !IsPiped(in) {
		return nil, ErrStdInIsNotPiped
	}

	reader := bufio.NewReader(os.Stdin)
	out := &bytes.Buffer{}

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		out.WriteRune(input)
	}

	return out.Bytes(), nil
}
