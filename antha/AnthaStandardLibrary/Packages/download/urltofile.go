// Package download provides convenience functions for downloading files


package download

import (
	"net/http"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"bytes"
	"io"
)

func File(url string, fileName string) (file wtype.File, err error) {

	//--------------------------------------------------------------------------------
	//Given a URL and a desired file name, this function will return the whole content of the response into a wtype.DownloadFile object.
	//--------------------------------------------------------------------------------

	//intializing global buffer object
	var buf bytes.Buffer

	//requesting
	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		return file, err
	}
	defer resp.Body.Close()

	//passing the response body to the bytes buffer
	if resp.StatusCode == http.StatusOK { // OK
		_, err := io.Copy(&buf, resp.Body)
		if err != nil {
			return file, err
		}
	}

	//creating the wtype.DownloadFile object
	if err := file.WriteAll(buf.Bytes());
	err != nil {
		return file, err
	}

	file.Name = fileName


	return file, err
}