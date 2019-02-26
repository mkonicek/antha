// Package download provides convenience functions for downloading files

package download

import (
	"bytes"
	"io"
	"net/http"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

//File takes a URL and a desired file name, and returns the whole content of the response into a wtype.File object.
func File(url string, fileName string) (file wtype.File, err error) {
	var buf bytes.Buffer

	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		return file, err
	}
	defer resp.Body.Close() // nolint

	if resp.StatusCode == http.StatusOK { // OK
		_, err := io.Copy(&buf, resp.Body)
		if err != nil {
			return file, err
		}
	}

	if err := file.WriteAll(buf.Bytes()); err != nil {
		return file, err
	}

	file.Name = fileName

	return file, err
}
