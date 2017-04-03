// Package download provides convenience functions for downloading files
package download

import (
	"net/http"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"bytes"
	"io/ioutil"
)

// File downloads the data at a url to the given filename. If there is an error, the file will contain the partially downloaded data.
func File(url string, filename string) (file wtype.File, err error) {

	//intializing global buffer object
	var buf bytes.Buffer

	//Downloading

	//requesting
	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		return file, err
	}
	defer resp.Body.Close()

	//converting body to bytes
	if resp.StatusCode == 200 { // OK
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return file, err
		}
		buf.Write(bodyBytes)
	}

	//returning wtype.File
	file.WriteAll(buf.Bytes())
	file.Name = filename


	return file, err
}