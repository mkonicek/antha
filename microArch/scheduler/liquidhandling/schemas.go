// Code generated by go-bindata. DO NOT EDIT.
// sources:
// schemas/actions.schema.json (16.769kB)
// schemas/layout.schema.json (8.148kB)

package liquidhandling

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _actionsSchemaJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x5b\xdd\x6f\xdc\xb8\x11\x7f\xb6\xff\x8a\x81\xee\x50\xdc\xa1\x1b\x3b\x79\x2a\xea\xb7\x00\xf7\x72\x45\xd1\x04\xb8\x6b\x5f\x82\x74\x41\x4b\xb3\x16\x2f\x14\xa9\x90\x94\x37\xdb\xc0\xff\x7b\xc1\x2f\x7d\x52\x12\x65\xaf\xd3\x3b\xd4\x7e\xf0\xee\x8a\xe4\x70\x38\x9c\xf9\xcd\x07\xa9\xaf\x97\x17\xd9\xf7\xb4\xc8\x6e\x20\x2b\xb5\xae\xd5\xcd\xf5\x35\xe1\xba\x24\x57\xb9\xa8\xae\x49\xae\xa9\xe0\xea\x95\xca\x4b\xac\x48\xb6\x33\x7d\xfd\x77\xdf\xff\xe6\xfa\xfa\x37\x25\xb8\xef\x71\x25\xe4\xdd\x75\x21\xc9\x41\xbf\x7a\xfd\x97\x6b\xf7\xec\x3b\x3b\xac\x40\x95\x4b\x5a\x1b\x72\x66\xe8\xdf\x7e\x79\xf7\x0f\xf8\xc5\xb6\xc3\x41\x48\x70\xcd\xb7\x94\xdf\x81\x9f\x13\x72\x22\x25\xc5\x02\x44\xa3\xa1\x68\xa4\x69\x62\xf4\x73\x43\x8b\x92\xf0\x82\x51\x7e\x97\xed\x2e\x01\x00\x32\x7d\xaa\xd1\xd0\x14\xb7\xbf\x61\xae\xc3\x53\x89\x9f\x1b\x2a\xd1\x2c\xec\x43\x76\x8f\x52\x99\x99\x77\x90\x79\xf2\xd9\x47\xdf\xaf\x96\xa2\x46\xa9\x29\xaa\xec\x06\xbe\xda\x67\xf6\x79\x18\xd2\x7f\x68\x1b\x46\x2b\xd1\x25\x82\xef\x0b\xe2\x00\xe6\xa7\x5b\xf7\xce\x2e\xec\x9e\x30\x5a\x10\xdb\x79\x37\xa4\x93\x0b\xae\xb4\xa1\xf0\xe6\xea\x4d\xd6\x36\x3d\x74\xbd\x5a\x56\x53\x58\x58\x90\x9a\x69\x1e\x4a\x0e\xcc\x92\xa3\x4c\x05\x59\x12\x29\xc9\x69\xdc\x58\x51\xfe\xb3\xc6\xca\x30\xf4\x66\xd4\x44\xfd\xf3\x21\xa3\xb6\x49\x70\x7c\x77\x30\xbb\x30\x69\x32\x7f\x5f\x21\xfb\x5e\xa2\x69\xcf\xbe\xbb\x2e\xf0\x40\x39\xb5\x0b\xb9\xd6\x92\x70\x75\x40\xf9\xd6\x2e\x2c\xeb\x0b\x26\x69\x7c\x2d\x45\x55\xeb\xc7\x8e\x3e\x12\xda\x8d\x9d\x0c\xfd\x38\x78\xd2\xb5\xbb\x6f\x0f\x4e\xdf\x5b\x62\x56\x2c\x17\x17\x99\xdb\x03\xff\x6b\x62\x11\x3f\x39\x0b\x40\x20\x7e\xb3\x80\x72\x20\x50\x33\xa2\x11\x8e\xc8\x98\x31\xa3\x8b\x8b\xa9\xb6\x9b\x87\x03\x65\xe7\xa4\x42\xa3\xe9\x5a\x68\xc2\xf6\xf7\x82\x35\x15\x1a\x75\x37\x1d\x47\xda\x7e\x61\x9e\xd9\xfe\xe1\xd7\x84\xaf\x5f\x4b\x84\x82\xaa\x9a\x91\x13\x98\x9e\x41\xc9\xfd\x6a\x76\x7e\x54\x60\x4b\x69\xa3\x73\x99\x7d\xfa\xe0\x1a\x87\x8c\x2c\x4e\x64\x7b\x82\xeb\x09\xb5\x44\x85\x5c\xcf\x4c\x18\xdf\xb7\x0a\x89\x6a\x24\x56\xc8\xf5\x90\x07\x25\x1a\x99\xf7\x56\x3d\x99\xde\x23\x10\xaa\xde\x64\x0a\x8e\x25\xcd\x4b\x38\xa2\x44\xc8\x45\x75\x4b\x39\x16\xa0\x85\xb1\xec\x0a\x74\x49\x55\x7f\xaf\xe4\x2d\xd5\x92\xc8\x13\x08\x59\xa0\x9c\x48\x26\x98\x94\x7b\xda\x19\xcc\x45\x94\x9d\xb7\xe0\x38\xf6\xc2\x08\xe3\xe2\xfb\x3f\xab\x02\x61\x30\x64\x0d\xa7\xda\x2b\xc1\x8c\x1e\x4c\x55\x21\xc2\x96\x91\x4d\x5f\x0b\x1c\x93\x3b\x2f\x27\x55\x8a\x86\x15\x50\x92\x7b\x84\x0a\x09\xb7\xe8\x23\x6c\xc7\x46\x75\x22\x99\xd3\x97\x6e\xbb\x4c\x8f\xb1\xbe\xcc\x70\xe3\x95\xc5\xf2\x43\x55\x90\x9a\x2e\x89\x86\x23\x51\x40\x8a\x02\x8b\xd8\xc4\xbc\xa9\x6e\x87\x2c\x55\x94\xd3\xaa\xa9\xb2\x1b\x78\x7d\xf5\x3a\xc2\x90\x15\xe1\x1a\x3b\xa6\x93\x32\xea\xe0\x24\xd2\xe3\x90\x2a\xc0\x2f\x46\xa7\x55\x9c\xa1\xa9\x24\x62\x88\xd5\xc2\x6a\xc0\xc7\x38\xea\x4e\xba\x47\x58\xa5\xbc\xc0\x2f\xa8\xc2\x56\x52\x5e\xd0\x7b\x5a\x34\x84\x41\xa0\xfd\x83\xfa\xd1\xaf\xc3\x8a\xb1\x2f\xe1\x1d\xd0\x03\x10\x3e\xf6\x11\x71\x56\x17\x9c\x4a\x74\xc0\xbc\x37\x59\xa4\xdf\x6e\x69\xd2\xa8\xaa\x61\x9a\xd6\xcc\x39\xa7\x37\x57\xaf\x53\x87\x0d\xb4\x64\xad\xfb\xd4\x77\x0c\x5b\x47\x4a\x96\x91\xa2\xb0\x48\x46\xd8\xfb\xbe\x85\x1e\x08\x53\x78\x39\xe8\xbb\xde\xd5\x52\x77\xdd\xd7\x3a\xdb\x5e\x59\x1f\x3c\x67\xbc\xd4\x5b\xe8\x75\xda\x81\x3e\xd5\x34\x27\x8c\x9d\x8c\x12\x11\xf8\x57\x0f\xab\x12\x1c\xd5\x3d\x61\xcd\x18\x9c\xa2\x1e\xca\x75\x5c\xf4\x1c\xb6\x4b\x50\xe5\xfe\x42\xc6\x30\xec\x95\x64\xe0\x1c\x86\x96\x1d\x23\x1f\x31\xeb\xde\x24\x51\xdb\x8e\xbb\xc4\x4b\xff\xaf\x1f\xe7\x19\xef\xfe\x77\x91\x13\x9d\x1a\x6f\x32\xdf\xd9\x49\xdd\x0c\x87\x23\xd5\xa5\x0d\x18\x46\x91\x9e\x14\xb7\x42\xcf\x45\x79\x83\x88\xb9\x6d\xed\xef\x91\x99\x3e\xff\xb4\x37\x06\xb9\x37\xce\x17\x32\x29\x8e\xe6\x23\x17\x2c\x83\x8f\xa3\x91\x33\xb1\x74\x6f\x29\x3d\x5a\x73\xe6\x3d\x96\x5b\xdc\x2e\x63\x9b\xf4\xf3\x4f\x56\x20\x1c\xcc\x14\xf0\x43\x2d\xa9\x90\xc6\xff\x0c\x45\xf2\x63\x36\x21\x18\xc1\x5a\xbb\xce\x55\x16\x17\x11\x27\xb6\x73\x52\x1c\xdb\x78\x26\xec\xf8\xcc\xe8\x34\x7c\x5a\xc1\xa3\xd8\xca\xcc\xd6\x9d\x7f\x65\xb9\xb1\x7d\xfe\x6d\x17\x77\xb9\xb0\xd4\x65\xc8\x8b\x0c\x32\x09\x99\x46\xae\xff\x59\x17\x44\xe3\xaa\x1d\x36\xb6\x9b\x8b\x16\xfd\x48\xd5\xd9\xe3\x53\x2c\x8e\x89\xdc\x58\x18\xc7\xe3\xde\x13\xde\x6e\x69\x86\xc6\xcd\x42\x7a\xd3\x47\x9c\xa8\x92\xf4\x67\x5f\x20\xe4\x83\xf2\x73\x6f\x45\x89\xf4\xae\xd4\xab\x7b\x60\x8c\xde\x75\x0d\x7a\x87\xbc\x08\x5f\x35\xad\x41\x22\x23\x9a\xde\x9b\xc4\x02\x94\xa8\xd0\x60\xe4\x18\x51\x36\xed\x8d\xc4\x03\x4a\xe4\xb9\x8b\xae\x07\xfe\x6b\xf3\x1e\x75\xb4\x66\xad\x31\x96\x8b\x11\xdd\x54\xb6\xb4\xa0\xdb\xd5\xcf\x59\x19\x72\x6b\x3b\x1f\x9c\x8b\xd9\xdf\x0a\xad\x45\x65\x18\xb6\x3f\xb5\xa8\xcd\x77\xb7\x85\x7b\x86\xf7\x68\x20\x3d\x06\x21\x8f\x72\xc2\x2d\x6b\xcf\xe4\x7f\xfd\xc6\x6f\x70\xbd\x67\x55\x51\x45\xf9\x1d\xc3\x5f\x7d\xa4\xbc\xaa\xaa\xbd\xf4\xde\x8d\x84\xbc\x24\x9c\x63\x17\x6c\x3f\x45\x2d\x6d\x28\x6f\x36\xf3\x20\xdd\x06\x6b\x31\x4c\xff\x8e\x44\x69\xb4\xfe\xbb\x16\x8c\xe6\x27\x5b\x09\x53\x76\xff\x0b\xf7\x51\x22\x29\xb6\xeb\xb0\x9b\x38\x55\x7f\xdb\xac\x03\x84\x34\xde\xba\x51\x2e\x9d\x6e\x38\xfd\xdc\x20\x3b\x81\x35\x09\x97\x32\x52\x35\x27\x99\x89\x84\x96\x7d\xd5\xb3\xf9\x51\x2b\xeb\x2d\x4b\xf7\xb9\x69\x3f\x7a\xb3\x40\xb5\xb2\xcc\x38\xf2\x0e\xbd\x55\x12\xc3\x5a\x6c\x62\xb7\x40\xa5\x29\xb7\xac\x9a\x44\x70\xc4\xed\x0e\x82\x60\x81\x56\x35\xa3\xa8\xdc\x83\x57\x05\x55\x35\x72\x85\x6b\xbb\xb6\x94\x0d\x76\x19\x60\xda\xf2\xe7\xea\x7b\x83\x92\x65\x92\x8c\x7a\x55\x87\x54\x39\x75\x15\x08\x5f\x0c\x0a\x22\x92\x68\x0b\x43\x56\x97\xe7\x2a\xae\x2d\xe1\x84\x7a\x56\x0a\xff\xde\xd0\x9f\xc4\xbf\xa3\xf1\xcd\x59\xf7\xc8\xb4\x85\xf5\x69\x41\xd2\xa4\xa4\xe8\xdd\xe3\x06\x08\x09\x8e\x22\x85\x4f\x83\x9b\xa9\x4c\xbe\x27\x92\x54\xa8\x51\x2a\x17\x8b\x38\xbc\x33\xec\x1e\xc9\xa9\xf3\x67\x44\xd5\xd4\xc9\xd8\x16\xad\x6a\x94\x07\x21\x2b\xeb\xd0\x36\x48\xbc\xa6\x35\x6a\x4d\xf9\xdd\xbb\xda\x15\x9e\x93\x96\x53\x9c\x7f\x39\x1e\x02\x4c\xfe\x99\xb4\x1c\xc2\xd8\xc2\x11\x01\x2c\x96\xf9\xc7\x6b\x9e\x83\x02\x4b\x65\xa5\xa4\xb6\xe4\x71\x27\xbd\x87\x1e\x58\x8b\x26\x2f\xc5\xe1\x30\x71\xa1\x93\x71\x75\x2b\xc3\xb4\xca\xdd\x2d\x13\x47\xd1\x4c\x83\xe1\xd9\x01\xa3\x3d\xab\xa3\x7b\x46\xc0\xd3\x5d\x35\xef\x09\xfd\x4d\x62\x6a\x47\x0d\xc5\xd5\x45\x26\x21\x4c\x84\xec\xc0\xc4\x71\x2f\x2d\x9c\xaf\xc8\xb0\x25\xba\x12\x9e\xcc\x8e\x5b\x41\xf9\xd9\x71\xf3\xe8\x69\x43\x19\x03\x3b\xad\x60\x6d\x85\x74\x4d\xf5\x67\x67\x7a\x04\xa6\xce\xfd\xad\x54\x92\x07\xd3\xce\x24\x5e\xab\xe3\x22\x82\xf1\xf1\xb9\x89\x92\x42\xd4\x2e\x82\x3c\x3a\x29\x9d\x47\x2e\x9e\xed\x67\x11\x49\xa7\x96\xe7\x90\x8a\xa1\x06\x86\x5a\x9b\xc2\x3d\xd6\x0e\xdb\x59\xce\xa9\x2a\x49\x3d\xd7\x7b\x25\xc8\xb7\x43\xcc\xc7\x02\xdb\xb1\x44\x5d\xa2\x34\xb9\x04\x17\x1a\x08\x04\x8a\x26\x27\xdc\x6a\x78\x2d\xaa\xdd\x0a\xc1\x90\xf0\x75\x99\xad\x95\xf6\xb7\xb5\x44\xf3\xee\x09\x93\x36\x49\xdb\x12\x20\x99\x01\xb6\x04\x2a\xcc\x27\xad\x95\xf9\x21\x24\x34\xbc\x7b\xe2\x72\xc6\xff\x7d\x8a\x75\xde\xfc\x7c\x12\x18\xac\x65\xe8\xdd\x01\xb0\x70\x23\x80\xdc\x1a\xbb\x2c\xc5\x11\x5a\x62\x36\x9e\xe9\xdd\xb1\x78\x4a\xd2\x1e\xf3\x7c\xa3\x7a\xcc\xfe\x20\x18\x13\xc7\xed\x79\x79\x87\xe3\x2b\x90\x19\xcf\x6c\xfb\x90\x97\x02\x30\x51\x2a\xb1\x85\xa4\xea\x2e\x3d\x80\x96\x0d\xee\xc0\x8d\xeb\x07\xf8\x96\x5e\xff\x6e\xcb\x6a\x7e\xb2\x6a\xda\x31\xee\x2b\xfa\xc5\x64\x04\xa9\x0c\xab\xa6\xaa\x88\xa4\xff\x41\xcb\x52\xb7\x3d\xfe\x44\xc2\xa9\x14\x61\xe0\xc8\xa6\xf3\xbc\x18\x64\x8d\x14\x2a\x3f\xe5\x0c\xd5\xb0\xf0\xf3\x04\x25\x6b\x67\x49\x8c\xb2\xc2\xfc\x6b\x78\x1e\xab\xef\x39\x7c\x31\xc2\xaa\xe8\x17\x70\x84\xfc\x21\xbe\xbd\x7b\x31\x6b\x71\xb3\x82\x4b\x3a\x0c\xde\x78\x08\x9c\x78\xf8\xbb\xe0\xfd\x7a\x91\xe7\x63\xed\xaa\x25\xf5\x24\x1b\x6f\xa9\x9c\xc5\xd6\x5b\x6a\x5b\x6c\xbe\x1d\xf4\x8c\xb6\xdf\xce\x91\xec\xde\xe3\xee\x79\xae\xbe\x94\xe4\x96\xe2\xa4\x1f\xe2\x6e\x8b\x48\xc2\x18\xb2\xe4\xc2\xb2\x44\x7f\x4f\x4a\xd9\xc3\x5f\xe5\x8e\x42\xc2\xe5\x10\x67\x43\x64\x68\x42\xa0\xa8\xd1\x7b\xc2\x51\x34\x8a\x9d\xe0\xf6\x04\x04\x14\xda\x91\xbe\x2a\xad\x9e\xe2\xd8\x3a\x1a\x90\x7d\xa2\xbc\xc8\x60\x07\x99\xa6\x15\xee\x51\x69\x5a\x79\x08\xca\x9b\xaa\x71\xa7\x33\xfb\x61\xdb\x56\x5f\x67\xa7\xb0\xca\xdb\x5e\xa9\x0c\x52\xdc\xb7\xf5\x9f\xf8\x49\x68\x60\x74\x4b\x48\x65\xeb\x96\x7d\x51\xed\xfc\x4d\x9a\xc2\x08\x32\x54\xf5\x97\x03\xa6\x34\x74\xaf\x89\xd6\x28\xf9\xfb\x44\xf8\xfd\xf7\x87\xd7\xaf\xfe\xfa\xf1\xcf\x4b\x76\x3c\x3a\xb4\x78\x56\x65\x0f\x97\xff\x06\x7b\x0b\xe1\x94\x27\x26\xdd\xd0\xcb\x6a\x30\xad\x10\x34\xf9\x84\xbc\x2b\xe6\x51\xae\xb4\x6c\xec\xad\x4c\x23\x73\x50\x98\x0b\x5e\x8c\x55\x75\x22\xe4\x95\xd8\x75\x72\xeb\x2b\x70\x3e\xab\xa0\xe9\x8b\xb0\xd7\x19\x27\x4b\x41\x90\x0d\x07\x25\xe0\x40\x24\x90\x83\xc6\xe9\xfa\xa0\x34\x71\xa6\xa8\x6a\x86\x1a\x8b\x67\x5d\xed\x79\xc3\x6e\x4d\x6b\x7f\x6f\x36\x35\xde\x86\xee\xfa\xb2\xad\x41\x0b\x52\xd8\x08\x29\x24\x29\x16\xea\x69\x7d\x5e\x44\x0a\x07\x5e\xdf\x04\x98\xe2\x5b\xd1\x1d\xce\x9a\x55\xba\x63\x64\xfb\xed\x25\x17\x3c\x1f\x58\xd7\x42\xd1\x70\x80\x35\x58\xba\x2d\x4b\x09\x09\x85\x14\x75\x2b\x0c\x63\x9e\x48\xf2\x32\xa0\xf8\xd5\xef\x0d\xbf\x47\xf7\x36\x5e\xd0\xfb\x05\xbd\x07\xbb\x16\x87\xe4\xe1\x9b\x14\xc9\xb8\x4c\x38\x20\xd7\x54\x76\x67\xbe\x1e\xa6\xc3\x75\x6f\xda\x0f\x39\x8d\x69\xed\x02\x96\xfc\x09\xc6\x71\x6c\x17\xa9\xf7\x62\x52\x26\xee\xfc\xf5\xd1\x3b\x29\x9a\x7a\x52\xac\xdb\x84\xf1\x01\xd7\xf3\x92\xb2\x42\x22\xff\x36\xd8\xde\x0b\x3a\xd7\x62\x4d\xcf\xd6\x16\xf8\xb2\x83\x96\xdf\xed\xd9\x70\xc8\x99\x76\xe2\x3e\x0f\x50\xcb\x6f\xf2\x84\xbf\xf9\x37\x7a\xda\xd0\x60\xad\x3e\x3c\x7f\xda\x37\xce\x8e\x16\xea\xac\x53\x2f\x0a\x0b\xe8\xb8\x76\x51\xe0\x05\x0d\xff\x78\x68\xf8\xe4\x12\x72\xff\x15\xb2\xf4\xf2\x31\xe1\x60\x5f\x88\x94\x88\x1c\x1c\x8d\xf0\x1e\x11\xa1\x5a\x59\xc1\x35\x0a\x25\x50\x5e\x3f\xad\x82\x1c\x20\xaf\x42\xa5\xc8\x1d\x6e\x43\xbc\xcd\x35\xe5\x58\x9e\x6d\x57\x17\x07\xbc\xc0\xd4\x16\xbc\xf3\x63\x9c\x87\xa0\x0a\x54\x29\x8e\x7c\xf4\x0a\x51\x94\x56\xc2\x3d\x8e\x17\x0b\xfe\xff\xb3\xe0\xde\x6b\x9c\xbf\x67\xfb\x2d\x1a\x17\x1a\xed\x5b\xa9\x7e\xe3\xd0\xc5\x2c\x2c\x6e\xc5\x13\xd6\x52\xcd\xb9\x14\x47\x60\xc2\xbd\x00\x68\xc8\x1b\xb1\x3d\x97\xe6\x8c\x9b\x5f\xac\xfd\x0f\x62\xed\x67\x72\x1b\x3f\x84\xe3\xb5\x1f\xd7\x3d\x88\x31\x65\x86\x56\x25\x17\xde\x2e\x5a\xf5\x27\x67\xc0\xa9\xcb\x8b\x07\xb8\x7c\xb8\xfc\x6f\x00\x00\x00\xff\xff\x1b\x6e\x4b\x38\x81\x41\x00\x00")

func actionsSchemaJsonBytes() ([]byte, error) {
	return bindataRead(
		_actionsSchemaJson,
		"actions.schema.json",
	)
}

func actionsSchemaJson() (*asset, error) {
	bytes, err := actionsSchemaJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "actions.schema.json", size: 16769, mode: os.FileMode(0644), modTime: time.Unix(1556702820, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x95, 0xc8, 0x79, 0x81, 0xef, 0x5, 0x2e, 0x9e, 0x60, 0xe, 0x47, 0x7c, 0xbe, 0x40, 0xa7, 0xc3, 0x96, 0xde, 0x5a, 0xad, 0xec, 0x61, 0x33, 0xd8, 0xdd, 0xc4, 0xf8, 0xc0, 0x98, 0xa7, 0x38, 0x1f}}
	return a, nil
}

var _layoutSchemaJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x59\x5d\x6f\xdb\x36\x17\xbe\x96\x7f\x05\xa1\x16\xc8\xc5\xeb\xc4\xe9\xdb\x8b\x61\xb9\x2b\xd0\x9b\x6e\xc3\x5a\xac\xc3\x76\x11\x64\x01\x2d\x1e\xc7\x6c\x25\x52\x25\x29\x3b\x6e\xe1\xff\x3e\x1c\x92\x92\x49\x89\x52\x6c\xb7\x1d\x30\x2c\x40\x6b\x59\x3a\x1f\x0f\xcf\xf7\x91\xbf\xcc\xb2\xfc\x39\x67\xf9\x0d\xc9\xd7\xc6\xd4\xfa\x66\xb1\xa0\xc2\xac\xe9\x55\x21\xab\x45\x49\x77\xb2\x31\x97\xba\x58\x43\x45\xf3\x39\x92\xfa\x6b\x4f\x7e\xb3\x58\x7c\xd0\x52\x78\x8a\x2b\xa9\x1e\x16\x4c\xd1\x95\xb9\xbc\xfe\x61\xe1\xee\x3d\xb3\x6c\x0c\x74\xa1\x78\x6d\xb8\x14\xc8\xfa\xd3\xfb\xb7\xbf\x92\xf7\xf6\x39\x59\x49\x45\x18\x14\x1f\x89\x53\x46\x42\x52\x64\x35\xbb\x1a\x90\x47\x2e\x3f\x40\x61\xec\x2d\x05\x9f\x1a\xae\x00\x41\xdf\xe6\x4b\x58\x49\x05\xf9\x9c\xe4\x74\x65\x40\xe1\x85\x80\xed\x3d\x67\x1a\x2f\x37\xa0\x34\x4a\xba\x43\xbe\x5a\xc9\x1a\x94\xe1\xa0\xf3\x1b\xf2\x65\x46\xfc\x5f\x47\x14\xde\xb4\x0f\x7a\xb0\xcd\x1a\x88\xa7\x25\x72\x45\xcc\x9a\x6b\xa2\x0f\xa7\xd8\xd0\x92\x33\xea\x81\x47\x72\x0a\x29\xb4\x41\x09\x2f\xae\x5e\xe4\xdd\xa3\xfd\x7c\x96\x65\x2d\x7e\xd4\x9d\x65\x59\xfe\x5c\xc1\x0a\x29\x9f\x2d\x18\xac\xb8\xe0\x28\x4e\x2f\xd0\x40\xef\x9b\xaa\xa2\x6a\x97\xcf\xb2\xcc\x71\xba\xf3\x9e\xc1\xd8\xda\xa7\x65\x1d\x9a\x38\x1b\xb8\xac\xa2\xb5\x26\x2b\x25\x2b\x82\x56\x70\x94\xe4\xcd\x6b\x4d\xb8\x20\x17\xee\x0c\x17\xc4\x48\x72\x61\x51\x5d\x78\x21\x94\x31\x0b\x84\x96\xef\x62\xdb\x67\xa1\x5e\x6d\x14\x17\x0f\x8e\x65\xa0\xf8\x95\x20\x6f\x5e\x5b\x2d\x5e\x32\x52\xed\x67\xf6\xdf\xde\xc5\x56\x77\x5c\x2f\x3a\x0f\x4f\x7d\xca\x21\x5f\xdb\xaf\x4b\xd0\xf6\x90\x3e\x20\xe5\x8a\x50\x52\xf2\x4f\x0d\x67\x6b\x2a\x58\xc9\xc5\x83\x0b\x58\x6a\x08\x25\x0f\x7c\x03\x82\xd4\x92\x0b\x83\x20\x0d\xaf\xc0\xcb\x8e\xa2\xb4\x96\xda\x43\xbc\x73\x4f\xeb\x84\x3d\x0e\x44\xed\xad\x34\xec\x21\xf0\x9f\x61\xa7\x09\x55\x60\x71\x0b\x5a\x81\x46\xd8\x9d\x3c\x22\x85\x7d\x62\x61\xdb\xc8\x85\xd4\x89\x3a\xe9\x93\x6e\x9b\x8c\xb5\x77\x5e\x65\xee\x48\xf7\xb3\xee\xff\x7d\x17\x7e\x11\xdd\x29\xee\xa1\xdd\x89\xda\x03\x8d\x1e\x22\x6d\x7d\xac\x08\x9a\x7f\x86\x63\x9c\x10\xf8\x20\x7d\xdc\x42\x4a\xc5\xb8\xa0\x06\xf4\x4b\x77\xdc\xbd\x8f\x60\xab\xe2\x04\xf6\xff\xc7\xec\xdc\x40\xf5\x74\x04\x48\x01\x6f\x51\xec\x6d\xeb\x95\x31\x45\x75\x49\x0d\x78\x7f\xb4\x3a\xb2\xec\x29\x2e\xc3\xeb\xa5\x7c\x3c\x87\x6d\x4b\x75\xa0\xcf\x7d\xde\x05\x71\x30\x59\x19\x56\xb4\xd4\xd0\x45\x4a\x64\xa3\x53\x22\xe5\xf7\xad\x24\x8c\x57\x20\xb4\x55\x41\x02\x41\xa9\xf8\x78\xbc\xaf\x2a\x8c\x8d\x1d\x7e\x4e\xc4\x86\xa5\x1b\x7a\x46\x34\xd5\x12\x3b\x4f\x3a\x37\x1f\x03\xed\x58\x22\xaa\x2a\x76\xf7\xee\x1c\xa1\xbb\x31\xa1\x5f\x61\xe2\x97\xa7\x99\x78\xad\x00\xfe\xa3\x46\xf6\x42\x3f\x9f\x23\xf4\xf3\xb7\xf2\x5c\x58\x26\xd2\x1e\x8b\xcc\xcf\x99\x1d\x8d\x68\x65\x67\x25\x4b\x3f\xc7\xb6\x2e\x9a\x15\x2d\x4c\xa3\xdc\xe8\x14\x4d\x5f\x24\x57\x72\x6b\xc7\xa8\x42\x96\x4d\x25\xec\x65\xe7\x71\xfb\x6d\x0b\x65\x79\x9f\xb8\xa5\x0d\x55\xa6\xfb\x66\xb5\x4d\xb8\xdc\xce\x9f\x7d\x33\x46\x53\x41\x22\xfc\x80\x34\x82\x7f\x6a\x00\xc7\x03\xdf\xd3\xfc\xe9\x23\x27\xd9\x13\x9f\x23\x1d\x19\xbb\x31\xaf\xa6\xca\xf0\xa2\x29\xa9\x4a\x2a\xf1\x52\x4f\x57\x82\x94\xb1\x26\x27\x7e\x4e\xe0\xea\xe1\x8a\x5c\xd4\x85\xb2\xf5\xfb\x22\xd6\x17\xb9\xed\xeb\x0e\x07\x24\x14\xd6\x83\x11\x6b\x8d\xe5\x9c\xa6\xf4\x15\xc1\x2e\xbd\x0b\xc7\xfb\x49\x5d\x36\xf2\x4e\xcd\x2d\x7b\x32\x4b\x81\xb2\x51\x04\x7e\x62\x08\x6a\x22\x15\x31\xbc\xb6\xf3\x6a\xa8\xb4\x95\x54\x35\xa5\xe1\x75\xe9\x7a\xea\x8b\xab\xeb\xee\x3e\x17\xbc\x6a\x30\xd1\xae\xaf\xae\x23\x84\x6d\x4e\x7c\x1d\x48\x2f\xe5\x3b\xe2\x0c\xb2\xf3\x00\x35\x81\x09\xa7\x96\xbe\x4b\xe6\x67\x4e\x40\xfd\xb2\x70\xa4\x62\x38\xce\x06\x67\xa1\x71\x15\x69\x12\xc8\x61\xba\x74\x60\x0a\x10\xa6\x4d\x09\x20\x4d\x5d\x83\x22\x25\xac\xcc\x65\x25\xb5\xb1\x50\x3d\xd2\x6f\x0a\xb4\x57\x4b\x46\x8a\xc6\x1c\x55\x2b\x6a\xd6\xa0\x2c\x3a\xbd\xa6\xf6\x66\x68\xc6\x16\x95\x9b\x00\xe7\xd3\xf9\x0a\xc2\x46\xcf\x6d\x5e\xec\x4a\x2e\x98\x6b\x07\x45\xb3\x94\x9c\xe5\x77\xa9\x31\xde\x49\x6d\x1b\x50\x0f\xe6\x6f\x50\x2b\xd0\x20\x70\x45\xb2\x84\xe1\x0a\xd2\x2e\x87\x65\xe9\x26\xd8\x68\x59\xc6\xbf\x2f\x64\xc4\x74\xb6\xe7\x11\x6f\xae\xd6\x42\x89\x8e\x92\x65\xf9\x47\x2e\x82\xae\x32\x44\xd8\x68\x60\xb8\xb0\x6e\x40\xf1\xd5\xce\x62\x73\x4d\xb1\x63\xe8\xb6\xf6\xfe\x04\xdd\x87\xdb\xed\xf9\x0a\x34\x67\x0d\x2d\xef\x37\x98\xd5\x30\x78\xa1\x30\x60\xe8\x61\x72\x6c\x64\xbb\xe6\xc5\x9a\x14\x54\x08\x69\xc8\x12\x88\x82\x4a\x6e\x80\x1d\xd6\x6f\x9f\x25\xab\xd6\xb8\xf9\x38\x26\xab\x26\x6d\xcc\x0a\xa8\x6e\x14\x54\x20\x4c\x3e\xca\xdf\x4d\xfe\x68\x0e\x03\xc2\xe8\xc8\xa8\x35\x35\x06\x94\x48\xae\x89\x59\x96\xff\x75\x7b\x7d\xf9\xe3\xdd\xff\x9e\x47\x77\x13\x7b\x9d\xab\x82\x5d\x11\x3c\x38\x61\x74\xff\x39\x42\xfb\xa8\xfe\x21\x02\x9b\xed\xfe\x7c\x61\x0e\xe1\x8e\x8f\xd7\xed\x96\xef\xeb\x03\x77\xb1\xec\x50\x87\x88\x46\x73\xde\x2d\xaa\x79\x40\xb9\x3f\x5c\xef\xc3\x33\x3d\x31\xfb\xf5\x98\x0f\xac\xc7\x30\x7a\xb6\x96\x29\x1e\x10\x6d\xc2\xcc\x87\x41\x4c\x82\xec\xbf\xeb\xb2\xdf\xef\x87\x47\xa4\xbf\xa3\x9c\xca\xff\xe9\x45\xd9\xe6\x7c\x58\x26\xff\xb1\xbc\x1f\xd9\x81\xf3\x8a\x6b\xcd\xc5\xc3\x3d\xf6\xa8\x29\x55\xbf\x70\x6d\xdc\xab\xa4\xc3\x9b\x18\x6a\x7c\x7a\xdb\x06\x47\x15\x10\x4c\x72\x6f\xad\x00\x44\x1b\xf5\x54\x29\xba\x0b\xee\xa3\x35\x7a\x49\xd6\xd3\xfa\x27\xc6\x6d\xd0\x63\xba\x31\xcf\xa1\x46\xc5\x41\xc8\x8e\xa7\x57\x1c\x1d\x4a\x6e\xfd\x22\xe0\x07\xf9\x09\x1f\xf8\xf9\xad\x9f\x8c\xe9\xf1\xc8\x9b\x34\x3d\xd6\xb4\xf6\xee\xcf\x36\xc3\xb4\x41\x60\xdf\x5b\xdf\x2c\xa1\xf9\xa8\x7c\xdd\x1f\x9d\x7f\x63\xb9\xe6\x5e\xaa\x1c\x97\x6d\x96\xf6\x5f\x9a\x6f\xf1\xcb\xa3\x2e\xe3\x18\xd7\x05\x55\x0c\xd8\x30\xe7\xc6\xbc\x3c\x3d\x76\x77\x02\x83\x41\x13\x88\xd3\x7e\x10\x31\x1a\x25\xc9\x10\x39\xd3\xbb\xbe\x35\x8c\xf8\xb6\x7d\x25\xdd\xbd\x83\x46\xac\xed\x4c\x85\x2d\xca\x3b\xf7\xe9\x17\x00\xdd\xea\x2f\xcd\xa1\xb8\x4f\x2c\xe5\xbd\xc5\x39\x61\x4e\xc6\x75\x5d\xd2\x5d\xb4\x4c\xfa\xd3\x8c\x0c\x9b\xd1\xb4\x1b\x01\x99\x1e\x78\x91\x92\xf8\xe9\xa8\x8d\xf5\xb4\xc2\x23\x46\x9d\x10\x83\x96\x8d\x2a\x60\x62\x41\x61\xf1\x4f\x02\x56\x99\xf6\x25\x7c\x0b\x0a\x87\x80\x6a\xc9\x85\x8b\xf1\x95\x54\x95\x1b\xb9\x03\x5f\xa9\x25\x37\x8a\xaa\x1d\x91\x8a\x05\x5b\x61\xb2\xc0\x0f\xca\xfb\x70\x8d\x76\x88\xbd\x31\xba\x58\x1d\x2b\xe3\xe9\x10\x68\x99\x49\xde\x08\x6e\x0e\x05\x7d\xa4\x9c\xf7\x42\x61\x64\x7e\x0a\xa3\xc0\x81\x9c\x7b\x3b\xe9\xb5\x6c\x4a\x46\xd6\x74\x03\xa4\x02\x2a\x6c\x17\x92\x6e\xa9\xd2\x51\x65\x4e\xc6\x4b\x3c\xe6\xf4\xe3\x65\x04\x8d\x0f\x96\xee\x67\x3a\x67\x35\xb3\xc6\xfe\x4b\x35\xa1\x8c\x01\x4b\x29\x1e\x34\x8b\x74\x37\x08\x00\x59\x13\x3e\x05\x07\x89\x6c\x99\xf1\xcd\xff\x80\x90\x6b\x02\x8f\x18\xd3\x3a\x0d\xa8\x6f\x89\x41\x6d\x7c\xba\x01\xed\x9f\xf8\x5d\xe7\x40\x7a\xf2\x4b\xc8\x30\xb1\x46\x2a\xd8\x2b\x12\x10\xcd\xb1\xf8\xf3\x82\x96\xe5\xce\xed\x2f\x7f\x04\x71\x7c\x44\x11\xdb\xd0\xb2\xe9\x07\x6e\xb2\x7a\x39\xc2\xc9\xaa\x62\x49\xba\xe9\x28\x38\x48\x3f\x45\x7d\x4c\x44\x85\x23\xf6\xfa\xc8\xdb\xc9\xbe\xcb\x03\x25\x49\xbf\xa7\xcb\xe5\x09\x3e\x99\x65\xfb\xd9\x7e\xf6\x77\x00\x00\x00\xff\xff\x5d\x73\xfc\x54\xd4\x1f\x00\x00")

func layoutSchemaJsonBytes() ([]byte, error) {
	return bindataRead(
		_layoutSchemaJson,
		"layout.schema.json",
	)
}

func layoutSchemaJson() (*asset, error) {
	bytes, err := layoutSchemaJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "layout.schema.json", size: 8148, mode: os.FileMode(0644), modTime: time.Unix(1556696510, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xae, 0xfe, 0x53, 0x53, 0xbb, 0x9b, 0x88, 0x5e, 0x5d, 0x3b, 0xa4, 0xb5, 0x83, 0x6b, 0x32, 0xd2, 0x74, 0x77, 0x3a, 0x6, 0x4e, 0xb2, 0x74, 0x42, 0x90, 0xce, 0xb8, 0x10, 0x6, 0x42, 0x4d, 0xb9}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"actions.schema.json": actionsSchemaJson,

	"layout.schema.json": layoutSchemaJson,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"actions.schema.json": &bintree{actionsSchemaJson, map[string]*bintree{}},
	"layout.schema.json":  &bintree{layoutSchemaJson, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory.
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
