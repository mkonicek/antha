package utils

import (
	"io"
	"os"
)

const (
	ReadOnly      os.FileMode = 0444
	ReadWrite     os.FileMode = 0666
	ReadWriteExec os.FileMode = 0777
)

// CreateFile creates a new file at the indicated path. The file is
// opened with O_CREATE|O_EXCL|O_RDWR meaning:
// a) there will be an error if the file already exists
// b) if there is no error, then the file handle returned will be
// available for reading and writing
//
// The perm parameter sets the permission bits on the file - the
// equivalent of chmod. Appropriate consts are provided. In general we
// rely on the user's umask to limit these, as desired.
func CreateFile(path string, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_RDWR, perm)
}

// CreateAndWriteFile is similar to ioutil.WriteAll, but it ensures
// that the file does not already exist.
func CreateAndWriteFile(path string, data []byte, perm os.FileMode) error {
	fh, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	n, err := fh.Write(data)
	if err == nil && n < len(data) {
		return io.ErrShortWrite
	}
	// always close, but preserve existing non-nil error:
	if errClose := fh.Close(); err == nil {
		err = errClose
	}
	return err
}

// MkdirAll is a simple wrapper around os.MkdirAll(). It exists to avoid
// having to consider which permissions to set on the created
// directories. We use 0777, and rely on the user's umask to limit
// these, as desired.
func MkdirAll(path string) error {
	return os.MkdirAll(path, ReadWriteExec)
}
