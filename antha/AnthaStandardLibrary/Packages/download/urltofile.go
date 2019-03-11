// Package download provides convenience functions for downloading files

package download

import (
	"bytes"
	"io"
	"net/http"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory"
)

//File takes a URL, and returns the whole content of the response into a wtype.File object.
func File(lab *laboratory.Laboratory, url string, filename string) (*wtype.File, error) {
	var buf bytes.Buffer

	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint

	if resp.StatusCode == http.StatusOK { // OK
		_, err := io.Copy(&buf, resp.Body)
		if err != nil {
			return nil, err
		}
	}

	return lab.FileManager.WriteAll(buf.Bytes(), filename)
}
