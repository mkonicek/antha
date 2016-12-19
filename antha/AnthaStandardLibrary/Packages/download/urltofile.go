// package download, for downloading files
package download

import (
	"io"
	"net/http"
	"os"
)

// func File downloads a url to the given filename.
// On error, the file will be left behind.
func File(url string, filename string) (f *os.File, err error) {

	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	f, err = os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := io.Copy(f, res.Body); err != nil {
		return nil, err
	}

	return f, nil
}
