// Package download provides convenience functions for downloading files
package download

import (
	"io"
	"net/http"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"bytes"
)

// File downloads the data at a url to the given filename. If there is an error, the file will contain the partially downloaded data.
func File(url string, filename string) (file wtype.File, err error) {

	//Downloading
	res, err := http.Get(url)
	if err != nil {
		return file, err
	}
	defer res.Body.Close()


	var buf bytes.Buffer

	if _, err := io.Copy(buf, res.Body); err != nil {
		return file, err
	}


	//returning wtype.File
	file.WriteAll(buf.Bytes())
	file.Name = filename


	return file, err
}