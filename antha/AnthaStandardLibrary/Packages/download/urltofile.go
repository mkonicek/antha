// Package download provides convenience functions for downloading files
package download

import (
	"io"
	"net/http"
	"os"
)

// File downloads the data at a url to the given filename. If there is an error, the file will contain the partially downloaded data.
func File(url string, filename string) (err error) {

  res, err := http.Get(url)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, res.Body); err != nil {
		return err
	}

	return nil
}
