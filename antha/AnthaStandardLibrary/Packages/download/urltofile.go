// package download, for downloading files
package download

import (
	"io"
	"net/http"
	"os"
)

// func File downloads a url to the given filename.
// On error, the file will be left behind.
func File(url string, filename string) (err error) {

	var f *os.File
	res, err := http.Get(url)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	f, err = os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, res.Body); err != nil {
		return err
	}

	return nil
}
