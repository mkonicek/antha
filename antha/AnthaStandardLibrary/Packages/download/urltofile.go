// package for downloading files
package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// get a file from a url link, the file will be given the specified filename
func UrlToFile(url string, filename string) (f *os.File, err error) {

	fmt.Println("getting url: ", url)

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

	f.Close()

	fmt.Println("made file: ", filename)

	return f, nil
}
