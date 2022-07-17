package terminal

import (
	"io"
	"os"
)

// FileWriter provides a minimal interface for Stdin.
type FileWriter interface {
	io.Writer
	Fd() uintptr
}

// FileReader provides a minimal interface for Stdout/Stderr.
type FileReader interface {
	io.Reader
	Stat() (os.FileInfo, error)
	Fd() uintptr
}
