// +build production
// Code generated by go-bindata.
// sources:
// templates/error.tmpl
// templates/main.tmpl
// templates/partials/footer.tmpl
// templates/partials/header.tmpl
// redirects/redirects.csv
// DO NOT EDIT!

package assets

import (
	"bytes"
	"compress/gzip"
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
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
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

var _templatesErrorTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x52\xc1\x8e\xd3\x30\x10\xbd\xe7\x2b\x46\xbe\x3b\xd5\x5e\x51\x1a\x21\x01\x12\x48\x70\x5b\xc4\x71\x35\xb1\x27\xb5\x55\xc7\x63\x79\x26\xe9\x56\x4b\xfe\x1d\xa5\xb0\x2c\xad\x7a\xe2\xf4\xc6\xcf\x6f\x34\xe3\xf7\xdc\xf9\xb8\x80\x4b\x28\xb2\x37\x05\x0f\x64\x03\xa1\xa7\x0a\x92\x50\xc9\xf4\xcd\xbf\xf7\xa7\x8a\xa5\x50\xdd\xd8\xf0\x70\xa7\xe9\xe9\x49\xa3\x26\x82\x31\xcd\x12\xac\x0d\x98\x46\x6b\x07\x56\xe5\xc9\xf4\x2f\x2f\xd0\x7e\xaa\x95\x6b\xfb\x78\x11\xad\x6b\xb7\x0b\x0f\x7d\xd3\xed\x7c\x5c\xfa\xe6\x15\xef\xcc\x83\x82\x99\x92\xb5\x13\xd6\x9b\x85\x0e\x35\x7a\xbb\xa9\xee\xf1\x8e\x13\x78\x92\xa3\x72\xb1\x17\x62\x9c\x53\xb2\xa7\xe8\x35\x6c\x72\xac\x1a\x5d\xa2\xeb\xc6\x81\x9f\x61\xe0\x67\x6b\x0b\x7a\x4f\xfe\x77\x5d\xc9\x5b\x8f\xf5\x78\x7d\xb2\x56\xa8\x60\x45\x25\x6f\x13\x8d\x0a\x42\x4e\x23\x67\xf2\x66\x7b\xcd\x7f\x6f\xa9\x27\xb6\x1a\x62\xf5\x62\xfa\xe6\xcd\xb4\x8f\x24\xae\xc6\xb2\x4d\x80\x9f\x20\x38\xd2\xe7\xc7\x6f\x5f\x61\x5d\x9b\xae\xf4\x5f\x46\x38\xf3\x0c\xa2\x31\x25\xa0\xec\x78\xce\xba\xf9\x56\x79\x48\x34\x09\x94\x44\x28\x04\x1d\x42\xa8\x34\xee\xcd\x84\x31\x29\xbf\x3b\xd1\xd0\x3a\x9e\x26\xca\x2a\xef\x39\x4b\x7b\xe0\xa5\x9d\x8f\x06\x2e\x31\xee\xcd\x07\xce\x8a\x4e\xe1\xbb\x98\xde\xfd\xa9\x67\xe9\x76\xd8\xb7\xf0\x83\x00\x0b\x27\x3e\x44\x21\x18\xb9\x02\xe6\x33\xc4\xec\x38\x2f\x94\x23\x65\x47\xa0\x21\x0a\x4c\x78\x86\x80\x0b\x81\xc3\x59\xc8\xb7\xdd\xae\xfc\xcd\xfc\x16\xde\x22\xb9\xfe\x13\xaf\xf8\x2b\x00\x00\xff\xff\x38\x7a\x05\x98\xac\x02\x00\x00")

func templatesErrorTmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesErrorTmpl,
		"templates/error.tmpl",
	)
}

func templatesErrorTmpl() (*asset, error) {
	bytes, err := templatesErrorTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/error.tmpl", size: 684, mode: os.FileMode(420), modTime: time.Unix(1578059292, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesMainTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x57\x6f\x4f\xdb\x4c\x12\x7f\xfd\xf0\x29\xe6\x71\xa5\x73\x52\xbc\x76\x5c\x8e\x5c\x20\x31\x12\x6d\x51\x0f\xa9\x2d\xe8\x4a\x75\x77\x42\x48\xdd\xd8\x63\x7b\xdb\xf5\xae\xbb\x3b\x4e\x1a\x05\x7f\xf7\xd3\x3a\x09\x04\x4a\x39\x1e\xa1\xbe\xb1\xbd\xde\x99\xdf\x9f\xd9\xf1\x66\x33\xf9\xf3\xed\xd9\x9b\x8b\xff\x9e\x9f\x40\x49\x95\x3c\xda\x99\x6c\x6e\xc8\xb3\xa3\x1d\x00\x80\x09\x09\x92\x78\x74\xaa\x08\x8d\xe2\x12\x2c\x9a\x19\x1a\x40\x63\xb4\x01\x06\x67\x79\x2e\x52\x84\x5c\x1b\xf8\xc8\x49\x68\x17\xf2\x89\x38\x09\x4b\x22\xb5\x93\x68\x95\xbd\x42\xaa\x90\x38\x94\x44\x35\xc3\xef\x8d\x98\x25\xde\x7f\xd8\xe7\x63\xf6\x46\x57\x35\x27\x31\x95\xe8\x41\xaa\x15\xa1\xa2\xc4\x3b\x3d\x49\x30\x2b\x30\x48\x4b\xa3\x2b\x4c\x62\x0f\xa2\x6d\x90\xb4\xe4\xc6\x22\x25\x5e\x43\x39\x1b\x79\x77\xe7\x36\x18\x73\x91\x51\x99\x64\x38\x13\x29\xb2\x6e\x10\x08\x25\x48\x70\xc9\x6c\xca\x25\x26\x71\x38\x08\x1a\x8b\xa6\x1b\xf2\xa9\xec\x78\x14\xaf\x30\xf1\x66\x02\xe7\xb5\x36\xe4\x6d\x23\xaf\xa6\x72\x6d\x2a\x4e\x2c\x43\xc2\xd4\x19\xde\x52\x4d\x28\xb1\x2e\xb5\xc2\x44\xe9\x07\x32\xa9\xc4\x0a\x59\xaa\xa5\x36\x5b\x49\x2f\xf6\x47\xfb\x07\xfb\xaf\x1f\x88\xe7\x75\x2d\x91\x55\x7a\x2a\x24\xb2\x39\x4e\x19\xaf\x6b\x66\x89\x53\x63\xd9\x94\x1b\x66\x69\x71\xa7\x68\xbf\x46\xd2\xca\x1e\xba\x85\x13\xe9\x76\x7c\x56\xb3\xdc\x74\x83\x8c\x19\xdd\x10\x9a\x4d\xaa\x14\xea\x1b\x18\x94\x89\x67\x4b\x6d\x28\x6d\x08\x44\xea\xac\x96\x06\xf3\xc4\x73\x6b\x68\x0f\xa3\x28\xcd\x54\xa8\x95\x0d\x0b\x3d\x0b\x9b\x6f\x11\xb7\x16\xc9\x46\xa2\xe2\x05\xda\x28\xe7\x33\x97\x13\x8a\x54\x7b\x40\x8b\x1a\x13\xaf\x9b\x89\x7e\xb0\x0e\xeb\x68\x67\xe7\x8f\x3f\x96\xcb\x2f\x93\x3f\x19\xbb\x14\x39\x48\x42\x38\x3d\x81\xd1\xd5\xd1\x17\xb8\x06\xcb\x73\xfc\xe7\xc5\x87\xf7\x6d\xdb\x29\xba\xab\xc9\xf9\xb6\x25\x22\x6d\x04\x2d\x97\xe1\x39\x27\xd7\xa1\xef\xc5\xd4\x70\xb3\x38\xee\xa4\x9c\x73\x2a\xdb\x36\x4a\xad\x8d\xb4\xcc\x98\xc0\x30\xb5\xd6\x3b\xda\xf0\x5e\xa2\xca\x44\x7e\xc5\xd8\x43\x8c\x5b\xca\x4e\x4f\xe0\xe0\xf7\xa8\x12\xc8\x0e\xd6\x9a\x6e\x38\x7f\xa9\xea\x4e\xb5\x0a\x5a\xcb\x72\x2f\x7e\x8b\xb6\x8a\x0b\x75\x5f\x1b\x63\x8f\xe8\x9b\x44\xab\x5d\x63\x32\xd5\xd9\x02\x52\xc9\xad\x75\x3c\x10\x5e\x2c\x6a\x84\xb6\x75\x2b\xbe\x42\x02\xc2\xaa\x96\x9c\x10\xbc\x9a\x1b\xf7\x45\xda\x2e\x17\x8d\x07\x21\xb4\xed\xce\xba\x81\xb9\x50\x20\xb2\xc4\x73\x0f\x1e\x18\x2d\x71\xf3\x4c\x7c\x2a\x54\x86\x3f\x12\x8f\xc5\x6b\x81\x2b\xe8\x85\x40\x99\xad\x4b\x30\xe9\x3c\x3c\xca\x9a\x6b\x4d\xf7\x58\x6d\x6a\x44\x4d\x1b\xcc\x5e\xde\xa8\xee\x33\xef\x89\xc0\x06\x3a\x28\x02\x13\xf0\xa0\xea\x2f\xc5\xa5\xff\x4e\xeb\x42\xe2\xb1\xe2\x72\xe1\x36\xbb\xb3\xe9\x57\x4c\xc9\xbf\x4a\xcc\x58\x5c\x9a\xab\xc4\x5d\xae\xaf\x6f\xf2\xfb\xcb\x0d\xa4\x9b\x08\xbf\x27\xab\xdb\xf5\xf5\xe5\x55\x3f\xac\x1b\x5b\xf6\xb8\x29\x9a\x0a\x15\xd9\x7e\x1b\x74\x93\x32\x89\x5f\x2a\x9c\xc3\x5b\x4e\xd8\xeb\x8f\x79\x62\xc3\xd4\x20\x27\x3c\x91\xe8\x02\x7b\xba\x1f\xac\x41\xab\xc4\x86\x05\xd2\x7a\xc2\xbe\x5e\x5c\xf0\xe2\x23\xaf\xb0\xa7\xfb\x97\x83\xab\x31\x0f\xb9\x5d\xa8\x34\x89\xc7\x3c\xb4\x26\x4d\x8a\x71\x15\xd6\xdc\xa0\xa2\x8f\x3a\xc3\x50\x28\x8b\x86\x5e\x63\xae\x0d\xf6\x9c\xbd\x35\x6a\xdb\xef\xcd\x85\xca\xf4\x3c\xc8\x74\xda\x69\x0b\xfc\x55\x7d\xfc\xc0\x8f\xa2\xf9\x7c\x1e\x16\x5d\x11\x18\xdf\x54\x21\x4c\x75\x15\xdd\x8e\xbe\x5a\x3f\xf0\x0b\xee\xf7\xc7\x6b\xc8\x82\xf7\xfc\x95\x09\x3f\x00\xff\xf3\x31\xdb\x1f\x8e\x0e\x5e\x0d\xf6\xfe\xc1\xf6\xfc\x00\x96\x3e\x97\x52\xcf\x8f\x55\x5a\x6a\xe3\x1f\x02\x99\x06\xdb\x3b\xb9\x16\x55\xe6\x32\x6b\x5e\xa0\xdb\xa5\xbb\x24\x37\xf0\x0f\x41\xea\xb4\xfb\x0d\x0a\x6b\x4e\xa5\xdb\xf9\x60\x17\x0a\xa4\x4f\xc8\x4d\x5a\xf6\xfa\xb0\x7b\x1b\x51\x72\x5b\xde\x02\x6f\x16\x69\x3b\x7a\x79\xd3\x55\xae\x79\xa2\x97\x70\x71\xf6\xf6\x0c\x18\xfc\xbb\x44\x05\xb6\x0b\x02\x61\xa1\xd2\x33\xcc\x80\x34\x18\x54\x19\x1a\x34\x30\x47\x5f\x4a\x50\xb8\x7a\xcd\xb3\x6c\x13\x4d\x68\x2a\x10\x8a\x34\x38\xbd\xf0\xee\x18\x0c\xda\x5a\x2b\x8b\x5b\x54\x51\x04\x22\xef\xfd\xec\x24\x49\x12\xf0\xa3\x15\x92\x7f\x47\x5c\x14\x75\xb7\x19\x37\xa0\x9a\x6a\x8a\xe6\x2c\xff\x17\xda\x46\x92\x85\x04\x96\xcb\x17\x22\x77\x3c\x8d\xa4\xf0\xde\x74\xdb\xc2\x72\xf9\xc8\x14\x4a\x8b\x6d\x0b\x03\x67\x5f\xe4\x6d\x3b\xfe\x99\xd4\x20\x35\x46\xdd\x96\x75\xed\x74\x17\xfc\xbf\xdd\x43\x4c\x7c\xd8\xbd\xaf\xef\x1e\x60\x0b\x8e\x11\x1e\xf0\xf6\x30\xcd\x4f\xe9\x2f\xa3\x9b\xcd\xef\x36\x6b\xd3\xca\x3b\x3b\x37\x91\x1f\x90\xdb\xc6\x20\x90\xa8\x10\xb4\x5a\x2d\x07\x83\xd4\x60\x26\xc8\x2d\x9a\xfb\x81\x3b\x8c\xa2\x12\x65\x1d\xde\xf4\xb2\x3b\x90\x74\xdd\xbd\x6a\xf9\x9b\xf7\x51\xd5\xc1\x09\x55\x30\x07\xc8\xb4\x62\x53\xdd\xa8\x14\x99\xc3\x8d\xee\xb7\x98\x8b\x31\x71\xdc\xeb\x2f\xb7\xdb\x19\x67\xa8\xc8\x3d\x5c\x88\x0a\xcf\xd4\xb9\x6b\xe8\x00\xfc\xb8\xbb\xc4\x6c\x6f\x00\x16\x53\xad\x32\xeb\xfa\x1d\x7c\xa5\x55\x77\x22\xe3\x1d\xa6\x7f\x08\x31\xb4\xfd\x71\xfb\x20\xd7\xde\x13\xb9\x5e\xb9\xcb\x5e\xcc\x86\xcf\xe0\x1a\x3e\x91\xcb\x7d\xeb\xfe\x30\x66\xf1\xe8\x19\x64\xf1\xe8\x89\x6c\x7f\xef\xaa\x38\x72\xd6\x9e\xe3\x6d\xf0\x44\xba\xfd\xce\xdc\xa0\x73\xf7\x2c\x7b\x4f\x25\x1c\xae\xfc\x0d\xe2\xdd\xbf\xc2\xf6\x04\xe0\x41\x77\x61\xf1\x93\x5c\xac\x61\x2d\x92\x03\xd1\x0d\xf5\xd6\xad\x1e\xc4\xf1\x60\x30\xf8\x65\xc0\x5e\x1c\xec\x3d\x1a\x30\x8c\x83\xe1\xa3\x01\xf1\x28\x0e\xe2\xd1\xe3\x18\x83\x38\x18\x0e\xfe\x0f\xca\xc0\xc1\x6c\x07\x4d\xa2\xcd\x69\x60\x12\xb9\x53\x8d\xbb\xaf\xfe\x21\xfd\x2f\x00\x00\xff\xff\xee\x50\x5e\x8c\x39\x0d\x00\x00")

func templatesMainTmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesMainTmpl,
		"templates/main.tmpl",
	)
}

func templatesMainTmpl() (*asset, error) {
	bytes, err := templatesMainTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/main.tmpl", size: 3385, mode: os.FileMode(420), modTime: time.Unix(1578059292, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesPartialsFooterTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xdc\x57\x4d\x6f\xe3\x36\x10\xbd\xef\xaf\x20\x78\xae\xcc\x60\x03\xf4\x50\xc8\x02\x82\xb4\xd9\x2e\x1a\x24\x87\xa4\xe8\x31\x18\x93\x23\x71\x6a\x8a\x14\xc8\x91\x14\x77\xb1\xff\xbd\xa0\x15\xc7\x8e\xed\x60\x3f\xd0\x36\x46\x6f\xe4\xcc\xbc\xe1\x1b\xf3\x71\x34\x2e\xeb\x10\x18\xa3\xd0\x0e\x52\x9a\xcb\x2e\x92\xe7\xa2\xb0\x64\x50\x56\xef\x84\x10\xa5\x7d\xbf\xf1\x0d\x94\x7a\x70\x6e\x65\xc9\x18\xf4\xb2\xba\x9a\x90\x8e\xfc\x32\x95\xca\xbe\x9f\xe2\x0d\x0d\x1b\xc0\x94\x7a\xca\xb3\xe7\x1a\x23\x74\xdd\x8e\x2f\xbb\x3d\x0c\x3b\xdb\xa3\xb9\x0a\x0f\x83\xd0\xc1\x15\x19\x2f\x5f\x46\xef\x01\x74\x70\xeb\xc8\xc2\x35\x45\xf0\x58\xb0\xa5\x68\x26\x4b\x6b\xb6\x96\xc3\x24\xeb\xa2\xcf\x0f\xcf\x7d\x78\xb0\x08\x86\x7c\x23\xab\x5f\xd1\x75\xa5\xb2\xe7\x47\xb1\xbd\x3b\x86\x75\x94\xf8\xe8\x59\x19\xe2\xe8\x18\x84\x18\xdb\xd7\x20\x19\x05\xc2\x46\xac\xe7\xf2\xd3\xa7\xd9\x3d\x3c\x06\x1f\xda\xd5\xcf\xa1\x05\xf2\x9f\x3f\x2b\x8b\xae\x53\xa0\x35\xa6\x44\x0b\x72\xc4\x2b\x59\x5d\xec\x6e\x4b\x05\xaf\xb1\x51\x8e\xfe\x63\xa2\x3a\x84\x25\x61\x02\x6f\xba\x48\x03\xe8\x95\xac\x2e\x27\x93\x00\x6f\xc4\x93\xf1\xa4\x28\x33\xc6\x36\x13\xd6\xc1\x1b\x62\x0a\x3e\xc9\xea\x3e\xdb\xd6\x8c\xb7\xd6\x6f\x27\x5d\xaa\xde\x1d\x2a\x5b\x19\x1a\xde\x4c\xef\x17\x8b\xd0\xb3\xb8\xbd\xb9\x3b\x6d\xd1\x43\xa6\xd9\x27\x35\x5a\xe0\x11\x4d\x90\xd5\x1f\x16\x58\x8c\x28\x4c\x38\x15\xf1\x6c\x38\x6a\x88\x88\x31\xc9\xea\x72\x5a\x9c\x1c\xbf\xe0\x19\x34\xf7\x99\xe1\xb4\x14\xfd\xc9\x90\xf4\x38\x26\x59\xdd\xe0\x78\x32\x8c\x36\x3f\x1b\x47\xf0\xa9\x83\x88\x5e\xaf\xc0\x9b\x26\x0c\x18\x3d\x78\x8d\xaa\x8e\x88\x26\xb4\xa1\x26\x5f\x87\xd8\x42\xee\x0e\x75\x20\x59\x5d\x4d\x0e\x11\x6a\xf1\x71\xeb\xfa\x5f\xf4\x8d\xcb\xe0\x3d\x6a\x16\x23\xb1\x5d\xcb\xe7\x54\xba\x87\x65\xee\xd2\x4f\x4a\xf1\x48\xcc\x18\x67\x3a\xb4\xea\xf6\xe6\x4e\x6e\x72\x92\x0e\xfe\x69\x0a\x12\x0c\xb1\x41\x9e\xcb\x87\x85\x03\xbf\x94\xd5\xfd\x84\x79\x0b\xe5\x6d\x68\x8f\xe3\x38\xab\x41\xe3\x22\x84\xe5\x37\x71\xbf\x7a\x02\xbd\x35\xf9\x3c\x33\xa2\x21\xbf\x26\xaf\x43\xdb\x81\x5f\xa9\x50\xd7\xa4\xb1\xa8\x43\x4e\x9c\xdf\x00\xb8\x22\x31\x30\x25\x26\x9d\xbe\xaa\xbe\xeb\x75\xde\x8f\xdf\xf1\x7a\xfe\xb9\xfa\xba\x7e\xe1\x48\xcf\x9a\x30\x18\x74\x34\x60\x5c\xad\xab\x04\xad\x43\xef\x39\xa9\xdf\x7f\xbb\xbd\xb9\x53\xa9\x5f\x24\x1d\x69\x81\x31\xe5\x76\xf6\x55\xd5\xfd\xd2\x02\x39\x01\x0e\x23\xff\x8b\x73\xc5\xbe\xa9\x54\x3b\x83\xf9\x0b\xe7\x97\x06\xfa\xc3\xf1\xdd\x91\x46\x9f\x70\xef\xe7\x2c\xa9\x6d\x8e\x07\x3e\x3c\x50\xdb\x48\x01\x8e\xe7\xf2\xf6\xc3\xb5\x14\x23\x19\xb6\x73\xf9\xe3\x99\x14\x29\xea\xb9\x54\xd4\x36\x2a\x34\x6e\xd6\xe5\x7e\xf3\x32\x6b\xf7\x5a\x4e\xc6\x47\x16\x2d\xc4\x86\x7c\xe1\xb0\xe6\x22\xb5\x45\x71\x76\x78\xc7\x17\x2e\xb7\x46\xcf\xe8\x59\x50\x12\x30\x00\x39\x58\x38\x14\xbd\x37\x18\x05\x5b\xcc\xb7\x7f\xe4\xe2\xb6\x82\x78\xd2\xfb\x46\xcf\x10\xb5\xa5\x01\x53\x56\xc7\xac\x5f\x2a\x13\xb4\x0a\x1d\xfa\x62\xfa\x50\xb4\xe8\x79\xa2\xa9\x51\x0d\x18\x13\x05\xaf\xce\xd5\xa1\x12\x6e\x3b\xf4\xe2\xc3\x33\x46\x5c\x4f\x18\x31\x9c\xcf\xce\xb2\x32\x44\x99\x3a\xf0\xbb\xd4\xc4\x9a\x1f\x3e\x72\xfe\x1e\xe5\x6e\x4f\x8d\xcd\x85\x83\x73\xb2\x2a\x55\x0e\xaf\x7e\x10\xf8\xa8\xb1\x63\x31\x5a\x8c\x28\x02\x5b\x8c\x23\x25\x14\xf9\x15\xa2\xd9\x53\x49\xf7\x42\x23\xbb\xb2\x78\xde\xec\xac\xb2\x16\xc8\xe4\x3f\x91\x38\x76\x21\xe6\xa3\x9f\x35\xff\x67\x2a\xb6\x66\xfa\x0b\x33\xa1\xd7\x81\xad\xf9\x4e\xa0\x6b\xbe\x0c\x2c\xd5\x24\x96\xea\xdd\xdf\x01\x00\x00\xff\xff\x67\x15\x8a\xeb\x19\x0f\x00\x00")

func templatesPartialsFooterTmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesPartialsFooterTmpl,
		"templates/partials/footer.tmpl",
	)
}

func templatesPartialsFooterTmpl() (*asset, error) {
	bytes, err := templatesPartialsFooterTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/partials/footer.tmpl", size: 3865, mode: os.FileMode(420), modTime: time.Unix(1577975052, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesPartialsHeaderTmpl = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x57\x4d\x6f\xdb\x46\x13\xbe\xe7\x57\xcc\xbb\xef\x25\x3d\xac\x58\xf4\xd2\x26\xa0\x08\xb4\xae\x81\x18\x48\xd2\xc2\x6e\x0e\x45\x60\x38\x63\xee\x90\xdc\x6a\xb9\xcb\xee\xae\x64\x0b\x2a\xff\x7b\xb1\xe2\xa7\x25\x52\xb6\xd1\xa0\x4d\x74\x22\x87\xf3\xf1\xcc\xcc\xb3\xb3\xa3\xb8\x20\x14\x64\x93\x17\x00\x00\x31\x42\xaa\xd0\xb9\x25\x73\x2b\x59\x29\xa9\x57\x0c\x0a\x4b\xd9\x92\xfd\xbf\x44\xa9\x19\x78\xbc\x95\x5a\xd0\xfd\x92\x7d\xcb\x92\xab\x95\xac\xc0\x1b\x08\x9f\x20\x35\xda\x93\xf6\x71\x84\xad\x2b\x21\x37\x20\xc5\x92\x55\x98\xd3\xaf\xe8\x0b\xd6\xb9\x2e\xa4\x20\x96\xec\x76\x8b\x0f\x97\x17\x75\x1d\x47\x42\x6e\x92\x17\x83\x4d\xab\x75\x67\xb1\xaa\xc8\xb2\xc6\xdb\xe1\xd7\x06\x34\xa4\x46\xf1\xa0\x38\xd2\x3a\xd4\x4c\x8d\xda\xab\x71\x95\x73\xa3\x89\xfb\x42\x5a\xd1\x48\x4a\x31\x48\x0e\x3c\xb4\xc5\x08\xf8\x95\xc9\x0d\x1f\x97\x22\x9a\xd0\x0d\xbf\xdd\xee\x53\xfc\x3f\xce\x3f\xca\x0c\x94\x27\xb8\x38\x87\x1f\xae\x93\x4f\xf0\x17\x38\xcc\xe8\xcd\x6f\xef\xde\xd6\xf5\xa4\xdd\x3e\x96\x2c\xf3\x0e\x71\x08\xc8\xc0\xd9\x74\xc9\x0a\xef\x2b\xf7\x3a\x8a\x52\xa1\x17\x46\xbb\x45\x6e\x36\x8b\xf5\x2a\x42\xe7\xc8\xbb\x48\x96\x98\x93\x8b\x8c\x76\x3c\xd8\x44\x9b\xef\xfa\xe7\x45\xa5\x73\x06\xa8\xfc\x92\xfd\x92\x65\x32\x25\xc8\x8c\x85\xf7\xe8\xa5\xd1\xa8\xe0\xca\xa3\x97\xce\xcb\xd4\x9d\xca\xe5\x23\x69\x21\xb3\x6b\xce\x9f\x92\xc5\x28\xfb\xbc\xc9\xfe\xd5\x75\x12\x24\xff\x59\x09\xdc\xe6\x5f\x2c\x41\xcf\xfb\x41\xb0\xe7\xf5\x03\xd1\x0c\x2d\xfd\x9d\x69\x48\xe8\x7a\x5e\x8e\x44\xe1\xb8\x70\xee\x4a\xa8\xac\xd4\x9e\xf3\xe6\xf8\x1c\x03\x10\xaa\xa1\x2b\xea\x7c\x8d\x39\xf5\xc7\xad\x17\x4c\xa7\x19\x0b\x7f\xa8\x79\x73\xe3\xa5\x57\xc4\x92\xb3\x02\x75\x4e\xd0\xc9\x5f\xc7\x91\xf0\x73\x5e\xc4\xb1\x17\xe9\xa9\x9c\x89\xda\x14\x58\x66\x40\x7f\xc2\xe2\x6d\x6b\x01\x8c\x34\xab\xeb\x61\x08\x0d\xae\xc2\xf9\x03\x99\x1a\xdd\xe6\xdf\x1d\xc5\x28\xdd\x2e\x76\xbb\xc5\x95\xf4\xf4\xb3\x09\x73\xa8\xae\xbb\xc1\xc2\x92\xb3\x6d\x69\x91\x72\x78\x79\xf6\xfb\x37\xa1\x3f\xbb\x1d\x69\x71\x82\x81\x13\x80\xd2\xed\x73\x01\xcd\xa2\x39\xd7\xb9\x92\xae\x80\x97\xe7\xef\x1f\x45\x13\x47\x42\x4c\xb4\x38\x12\xea\x59\x24\x73\x94\x1a\x2d\xd0\x6e\xb9\xc6\x0d\x3c\x9d\x72\x8f\x30\x6d\xdd\x30\x4d\xe3\x86\xf7\x11\xd8\x64\xcc\x50\x27\xe7\xe1\x0f\x17\x5e\x78\xaa\x8c\xa6\x46\x34\x47\x46\x25\x67\xfc\x3c\xc2\xa5\xd1\xc5\x75\x18\x5f\xaf\xe0\x08\xc0\x68\x98\x5b\x52\x84\x8e\x52\x54\x14\xcc\x58\x72\xd9\x08\xa0\x93\x1c\x9d\xec\xa1\xf8\x4a\x7e\x31\x69\x94\xe4\x0b\x23\x8c\x32\xf9\x96\x25\xef\x86\x97\xaf\x03\xbd\xa6\x3b\x17\x60\x0b\x89\x5f\x07\x60\xbc\x35\x6b\xbf\x76\x2c\xf9\x31\x3c\x7c\x39\x98\x87\x89\x34\x89\xde\xa3\xcd\xc9\x2f\xd9\xcd\xad\xc2\xf0\x6e\x49\x2d\x99\x36\xa6\x22\x4d\x16\xb4\xb1\x94\x91\xb5\x64\xbb\x3c\xbb\xcb\xf7\x56\x99\x7c\x7c\xfb\xb2\xe4\x27\x65\xf2\x67\xa6\x1d\x47\xeb\xd3\xe3\x6b\xf4\x3a\x7e\x1c\x0d\xb4\xca\xca\xb2\x1b\x67\xd3\x43\x2a\xd6\x78\x38\x10\xd7\xaa\x33\x0f\xf5\xe0\x61\x51\xb5\x46\x4d\x5d\xfb\xa3\xa6\x3c\x50\x3d\xd9\x94\x18\xbb\x15\x39\xd8\xb4\x08\xd9\x7e\x42\x96\xa4\xd7\xdc\x9b\x3c\x57\xc4\x00\xad\xc4\xde\x63\x13\xa0\x57\x9e\x0e\x1a\xcc\x4f\x31\xc1\x55\xa8\x67\x4c\x3d\xdd\xfb\x70\xa0\xf4\x3a\x8e\x82\xda\x5c\x9f\x26\x1a\x38\xd3\xbc\x93\xa5\x81\x27\xd5\xc6\x11\xda\xb4\x68\x4a\xd3\x3c\x9f\x2a\x4e\xa7\x3d\x1d\xb5\xfd\xfa\x4f\xaa\x73\xb5\x77\xf1\x19\xea\x33\x41\xec\xe3\x3f\x32\x10\x30\x84\x9d\x60\x4f\x59\x41\x9a\xf5\x97\x68\x4f\x83\x7d\x15\xe8\xbe\x42\x2d\x48\x2c\x59\x86\xca\xcd\x5d\xc0\xc7\x07\xe2\xa9\x17\xeb\x03\x93\x7d\xf3\xfa\xbd\xe7\xc3\xe5\x05\xb0\x88\xd5\xf5\x91\x0e\xe7\x98\x7a\xb9\xa1\x76\x69\x69\x87\xcb\xb0\x9b\x1a\x25\xb8\x24\xce\x85\x74\x95\xc2\x2d\xbf\x55\x26\x5d\xcd\x76\x67\x98\x61\x07\xf8\xf5\x6a\xd8\x52\x4a\xc1\xbf\xef\xf7\x95\x57\xa3\x3f\x5f\x6f\x4c\x49\x9f\x7f\xf4\xf4\x73\xa3\xfd\x12\x47\xdd\x1f\xe3\xbf\x03\x00\x00\xff\xff\x92\xea\xe8\xf4\x22\x0f\x00\x00")

func templatesPartialsHeaderTmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesPartialsHeaderTmpl,
		"templates/partials/header.tmpl",
	)
}

func templatesPartialsHeaderTmpl() (*asset, error) {
	bytes, err := templatesPartialsHeaderTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/partials/header.tmpl", size: 3874, mode: os.FileMode(420), modTime: time.Unix(1577975052, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _redirectsRedirectsCsv = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd4\x5c\x4b\x93\xdb\x38\x92\xbe\xeb\xb7\x94\x8a\x6d\xcf\x74\x6f\xc4\x46\x6c\x6c\xec\x76\xec\x61\x2e\x33\x07\xff\x00\x47\x0a\x48\x92\x69\x81\x48\x16\x1e\x92\x35\xbf\x7e\x03\x00\x9f\x12\x25\x91\x12\xcb\x5d\x75\xb1\x45\x3c\xf2\xfb\x12\x8f\x44\x22\x01\x54\x86\x82\x35\x57\xa7\xac\x30\x6c\xad\xe4\x0a\xad\x23\x51\x1b\x96\x5e\xb8\x42\xd6\xd9\x5b\x45\xf6\x5a\xe6\x5b\x45\x2f\x77\x05\x54\xe8\x4a\x96\xac\xb8\x20\xbc\x25\x69\x93\xd5\xc8\xb5\xc2\x9a\x6b\xaf\xc0\x11\x6b\xd0\x52\x70\x55\x79\x4d\xee\x94\x95\xec\x2d\xe9\x22\xf1\x09\x1f\x58\x1b\x12\x68\x1d\x38\x0a\xb2\x6c\xce\xc6\x56\xa0\x14\x18\x04\x1b\x99\xcd\x92\x37\xa6\x37\x47\xf0\xa6\x53\x59\x27\xb1\x0a\x84\x60\xaf\x9d\xcd\xfc\xde\xa2\x70\x6c\xba\x84\xc8\x96\xb4\xc4\x9f\x9c\x5b\x34\x87\x20\x78\xd4\x6a\xf7\x45\x8c\x09\x5e\xca\xda\x64\x58\xd5\x8a\x4f\x15\x6a\x07\x5a\x2a\xd8\xb1\x37\x15\x98\x3d\xba\xa6\x05\x48\x1f\xd9\xec\xb3\x94\xd3\xb4\x39\x1d\x42\x23\x44\x7a\x97\xe9\x89\xe1\xc3\x52\xc7\x8c\x27\xc5\xf7\x6d\x58\xf0\x01\x8d\x0e\x30\xb5\xdf\x29\x12\x8d\xf6\x5a\x3a\xf8\x89\x36\x1b\x26\xe6\xa4\x41\x0b\x4c\xac\x27\x32\x2c\x68\xd9\x8b\x93\x98\x93\xa0\x40\x5e\xe2\xce\x79\x2d\xd1\xb8\x12\x2b\x00\xeb\x0c\x89\xd2\x39\x83\xd0\xaa\xfa\x38\x97\xb1\xae\xeb\x91\xea\x1b\x88\x74\xde\x0d\xdf\x38\x32\x49\xcb\xd0\xf5\x4d\x33\xc4\x76\x45\x33\xcc\x19\xe9\x74\xb5\xfa\x19\xf3\x69\x39\x9b\x6c\x17\xa6\x09\x5a\x4b\x5a\x7a\xeb\xcc\x29\xb4\x86\x01\x89\x99\x41\x07\xa4\xda\xe4\xc4\x26\xa5\x59\x50\x68\xe3\x40\x8d\x4c\xe6\x4a\x18\x13\x3a\x17\x65\x2c\x2d\x68\x96\x7e\x1e\xc7\xca\x65\x4d\x8f\x34\x8a\xdf\x4f\xca\x79\xd4\x00\xb4\x73\x80\xf5\x53\x26\xc0\xef\xa7\xe4\xdd\x20\xb5\x03\x15\x46\x22\xe7\x35\xc4\xf9\xdc\xb0\xf2\xfb\xd8\x0d\xb7\xb9\x5c\xd6\x3d\x27\xd3\x0a\xb9\x63\xbf\x47\xc9\x15\x15\x26\xfe\x1e\x24\x87\x15\xa1\x02\xd7\xf6\xe0\xd7\xdf\xbe\x7c\x11\xa8\xad\xb7\x23\x4b\x8c\xba\x50\xa0\x25\x68\x79\x0c\x43\xa3\x02\x23\xca\x50\xf4\xbe\xc1\x5f\x40\x60\xac\xe2\x72\x26\xf7\x96\x32\x04\xe5\x4a\xd0\xd2\xb2\x20\x50\x02\x0c\x66\x02\xbc\x45\xcb\xb9\x44\x70\x65\x6a\x01\x50\x82\x4b\x56\x06\x15\x38\x94\x31\xc3\x92\x76\x25\xfa\xfd\x8c\xe5\xed\x2e\xc6\x58\xc9\xdb\x60\x37\xcc\x40\x9b\xd1\xfd\x68\x97\xa6\x46\x09\xad\x3d\xa8\x2e\xcf\x9b\x03\x9e\x6e\x5b\x85\xeb\x02\xcf\x18\x4f\x4b\x5e\x79\x18\x26\x94\x8a\xe4\x09\xc1\x4c\x94\x7a\xc7\x71\x37\x03\x7a\xf6\xda\x3f\x2a\x96\x3e\x10\xdd\xa9\x1e\xab\xd9\x83\xa4\xe6\x84\xda\x2e\x71\x05\x6e\x80\x4c\x69\x36\x8d\xb6\xd2\xdc\x39\x30\x49\xd8\x29\xac\xd8\x38\x50\xe4\x4e\xa4\xcf\x26\xec\x3b\xcc\xa1\x59\xa0\x77\x14\xdc\x91\x71\xa5\x4d\x13\x30\x0c\x18\x30\x86\xa0\x40\x9b\xa5\xa4\xa4\x5e\x2c\x14\x26\xac\x71\x2e\x00\x06\x7b\x44\x3a\x87\xe0\x60\x84\x52\xf7\x75\xbb\x07\x33\xd6\x6c\x06\xde\x1c\x13\x01\x8d\x03\x68\xe9\xdf\x18\x46\x12\x8b\x34\x17\x92\x4e\x4d\x29\x89\x15\x17\x06\xea\x72\xae\x9d\xb8\x26\xf5\x4c\x85\x29\xf1\x43\x77\xe2\x80\xd6\xc5\x15\xae\x46\x6d\xc3\xa0\x8c\x48\xde\xb6\xcb\x65\x4f\xa4\x2d\x79\xe6\x52\xdc\x10\x30\x4d\x65\x24\x69\xb5\xa9\xdc\x4a\x37\x58\x90\x75\x68\xfa\xf2\x69\x9a\xed\x0c\xbe\xcf\xac\x9e\x09\xfc\xc8\x04\x67\x2d\xc9\x35\x4d\x2a\xc9\x22\xd8\x56\x5b\x11\xbc\x13\x93\x20\x93\x65\xed\x17\xe9\x07\x67\xf8\x24\xd8\x58\xd5\x3b\xa8\xef\xa1\x62\x68\x43\x3a\x80\xea\x80\xc2\xda\xa7\x14\x3a\xd2\xef\xae\xe8\x1d\xec\xc5\xfe\xb0\x80\x9a\x5c\x10\xc7\x62\x1f\xc0\x9b\x6f\x11\x56\x82\xaa\x7e\xda\x41\x9e\x25\x7e\x0d\x33\x2c\x4a\x52\xb2\x33\xf6\x4b\x06\xde\x32\xf3\x7b\x03\x67\xb6\xd9\x18\x6e\x4d\x6b\x34\x96\xb5\x46\xd5\xa8\x41\x07\x52\xad\x97\x35\x56\xe2\x29\xe1\x67\x4a\x5c\x43\x59\xcd\xf2\x09\x05\x54\x81\x76\x71\x50\xbc\x8b\x89\x3b\x47\x78\x74\x10\x09\xd6\x02\xeb\xa6\x7c\x8e\xc6\x51\xe8\x58\xd3\xbb\x9d\x7d\x81\x55\x46\xd5\x2d\xbc\x33\x0d\xa7\x81\x6f\x4c\x71\x0b\x0e\x95\x22\x87\x67\x73\x3c\xce\x36\x34\xce\xa0\x96\x77\x62\x5e\x97\x22\x2e\x48\x8d\x65\xad\x31\x79\xdb\x3d\xce\x81\xd5\x81\x74\x21\x14\x5b\x67\x48\x92\xaf\x24\xe5\x39\x09\x52\x38\xd3\x4f\x5c\x36\x99\x17\xe0\xbe\x83\x9a\x95\xb1\xf0\x2b\xd4\x1a\xe2\xac\xa7\x46\xb3\x3d\x75\x2c\x8d\x2f\x6a\x26\xcb\x9a\x74\xf1\x80\x6b\xff\x88\x6e\xf3\xc1\x1f\x56\x98\x0e\x6c\xda\x00\x67\xf3\x61\x57\xd5\xae\x01\x38\x53\xef\x06\xd2\x0d\xb7\x3e\x86\x31\x49\x3b\x34\x3a\x78\x02\xc3\x38\x20\x06\x2e\x68\x04\xde\x76\xe1\xa7\x25\x8c\xc9\x0d\x45\xf5\x56\x08\xf5\x81\x0c\xc7\x90\xea\xc0\x8e\x24\xec\xa9\x2c\xd6\x40\x06\x2b\xb2\xd1\x31\x1f\x59\xa3\x69\x51\x67\x24\x66\xc9\x7c\x86\x1e\x6a\x34\xc5\x69\x2d\x62\x9d\xb4\xe7\x28\x0d\x92\x63\x04\x7c\x3d\x7a\x13\x92\xd7\xb0\x12\xf8\x53\xa0\xb5\xc7\x38\xa4\x9e\xd9\xf9\x2f\x33\x0f\xb3\x51\x97\x84\xd1\x48\x6b\x3e\x0c\xb6\xc7\x39\x1b\xa4\x42\x4b\x32\x28\x5c\xbf\x73\xcc\x25\x2d\x0c\xa7\x0d\x04\x8f\xd5\xb8\x89\xb0\x4e\x58\xa6\x48\x9e\x1f\x6b\x5b\xa3\xa0\xb0\xd8\xc5\x20\x42\xd7\x68\xeb\xc7\x64\xee\x21\x4e\x77\xd3\x5c\xbf\x11\x4c\x58\x02\xc2\x10\x09\xdf\xa4\x8b\x92\xbd\x19\x07\xfd\x53\x6d\xc1\xd6\xd9\x1a\x4d\xc8\x5f\xe4\x99\x5e\x41\x98\x3c\x14\x9c\x84\x9a\x7b\x62\xd2\x9d\x29\xd0\x01\x1c\x36\xa7\xa3\x26\xce\xd0\x74\x1e\xfa\xc8\x01\xca\x1c\xa1\x4b\xce\x7d\xc3\x14\x2b\x82\xdb\x2a\xe3\xc1\x4c\xc9\x4a\xa6\xd8\x02\xcd\x89\x79\x4d\x1f\xf9\x5e\x95\xf9\x70\x30\xb9\xfb\x75\x24\x57\x36\x41\xf4\x9e\xbf\xd1\xa0\xba\x12\xcf\xc7\x91\x27\xc1\x2e\x35\xbc\x8a\x3a\x58\x21\xe2\xff\x24\xd8\xbb\xda\xbb\xd8\xb1\x83\x43\xdd\xe1\x47\x85\x60\xbd\xe9\x47\x4e\x14\x9f\x7c\x7a\xc1\x55\x0d\x26\xf8\x46\xb6\x3f\xa0\xea\xcf\x94\x1f\x87\x9a\x54\xe9\x3e\xe6\x9d\x1e\x54\x48\x41\x7c\xb0\x96\xec\x0d\xd9\x6a\x42\xa5\x1a\xac\x45\x5d\xa4\x30\x08\x9e\xa8\x9e\xd1\x59\x97\x72\x6f\xf0\x3f\x03\xb8\xe3\x75\x0d\x6b\xa6\xa4\x4b\xca\x31\x9d\xf4\xe8\xc6\xc1\x12\x89\x37\xc8\x5e\x8a\x9e\x31\x81\xe3\xac\x12\x25\x18\x10\x0e\x4d\xda\x55\x66\x25\x57\xd8\x3a\x7f\x9d\x41\xaf\x50\x12\x78\x0b\xc5\x48\x2d\x74\x20\xc2\x1a\xdb\xcf\xd0\xe8\x38\x4a\x3a\x90\xf4\xa0\x66\xce\xfd\x87\x58\x4c\x35\xc5\x0c\x3a\xb3\xd7\x90\xf0\x4f\xad\x40\xa0\x24\x5b\x7b\x87\x03\x53\xdf\x07\xe8\x86\x17\x35\xda\x72\x83\xb3\xb4\x55\xa1\xa6\x6e\x6f\x5c\x60\xae\x16\xb3\x49\x35\xf3\xb8\xff\x88\xc2\x55\xfe\x3e\xd1\xe9\x2b\x40\x2b\x2b\x92\x6a\xbe\xa3\x02\x1d\xc0\x03\x0e\x99\x24\x0b\xbb\x18\xf9\x69\x28\x53\x8e\x5c\xd7\x6c\x5c\xa8\x43\xa3\x21\xf5\x8c\xf0\x33\xde\x57\x51\xee\xad\xae\x31\x90\x08\xc1\x39\xeb\x26\x5a\x7b\xa9\x26\x23\x2d\xb8\x0a\xf6\xf5\x18\xb9\xb4\x0a\xc5\x88\x4a\x70\x80\x40\xcb\x9c\x59\xce\xd5\x68\x11\xd6\xb9\x7e\x57\x41\x9f\x3d\x8b\xae\x0d\xff\x40\x31\xb0\x00\x15\x98\x14\xd9\x06\xe7\xed\x64\xb9\x67\x0f\xa3\x87\x90\x63\x3d\xef\x63\xdf\x5e\xb8\x44\x09\x3a\x78\xe3\xdc\xed\x48\x2a\x34\x61\xd1\x0b\xce\xb7\x78\xf3\x64\x87\xb6\x6e\x3a\xaf\x82\xdb\x4b\xd9\x6c\x8c\x33\xd5\xae\x83\x3d\x7e\xf3\xa1\x62\xed\x4a\x75\x7a\x8f\xab\x0f\xd7\x44\xaf\x66\xcc\x42\xa1\x68\x2a\x7f\xf0\xee\x7d\xcc\xf1\x39\xc2\x63\x53\x65\xe4\x99\xf4\xc9\x51\x07\x1b\x4f\xab\xd1\x54\xd3\x85\x3a\x1f\x38\x67\xa3\x58\x80\x02\xef\x4a\x36\xd1\x42\x3d\xec\x8c\x5f\xe1\x33\x56\x7d\x01\xb1\xc5\xa1\xc5\xf7\x6d\xa7\xc5\x71\x94\xbf\xb0\x99\x9e\x3f\x16\x1b\x66\xf5\xd2\xd6\x3e\x16\xbb\x8a\xf2\xc9\x2c\xcf\x62\x87\xa1\x4b\x0b\x0e\x02\xfe\xac\x51\x38\xd0\x82\x5a\x25\xb8\x26\xdd\x9c\x4d\x87\x02\xd6\x9d\x14\x3e\xe3\x9f\x5c\x87\x1b\x6b\x78\x13\xf7\x61\x35\xc3\x4f\x7b\xb2\x0e\x9b\x0d\x66\x3b\x9c\x77\x68\x10\x0e\xd8\x5f\x92\x38\x70\xb7\x61\x7b\x1e\x6a\xac\xd9\x7d\xcc\xd5\xd6\x8f\x03\x08\xd0\xe2\xb4\x7c\x8b\x32\x7b\xfd\x38\x47\x98\x1d\xc6\x68\xa6\x5b\x1a\xdd\x97\x97\xe1\x47\xf9\xc3\xec\xce\xde\xa0\xf4\xe9\x82\xd1\x92\x88\xc6\x75\xd4\x69\x6b\x30\x0f\x7e\x6e\x6c\x6f\x0c\xfc\xfc\x65\xf1\x2b\xf2\x9a\xba\xbf\xb4\x37\x1c\x3b\x50\x63\x31\x7f\x4d\xbf\x4c\x13\x59\x7d\x8b\xe3\x4a\xc4\x3c\x47\xe1\x2c\xe7\xf1\xe4\x04\xb4\xdc\xa1\xc6\x9c\x9c\x65\xdd\x09\x49\x75\xdf\x75\xdf\x73\x9b\x49\x73\x75\xbe\xe7\x72\xa7\x25\x50\xb5\xbb\x8e\xf8\x8b\x0d\xa8\xe1\xf5\xaa\xf6\xc6\x7a\x97\xb9\xe4\x5a\xc4\x3d\xd9\xe7\x37\xda\xa7\x41\x96\xdc\x15\xac\xde\x7e\xef\x0b\xed\x4e\xa4\xad\x37\xa1\x4d\x63\xa4\x52\x13\x76\xf5\x72\xaf\x65\x5f\xf9\xd1\x0b\x85\x8f\xc3\xf5\x4a\x85\x16\x19\x5d\xca\x48\x2f\xb5\xc8\xd6\x6c\x61\xa7\xf0\xac\x3b\xdb\xf7\x1e\xa9\xd2\xed\xb2\x23\xb5\x96\xe2\x9c\xbf\x0a\x99\x09\xb8\xa4\xb7\xf8\xa8\xd1\xd8\x92\x6a\xce\xfd\xfe\xcd\xb3\x43\x69\x4b\x30\x17\xb6\x71\x76\x7f\xdc\x12\xf8\xf4\x3d\x59\xbf\x6f\xcb\xad\x7a\x3f\xb6\x17\x7b\xa5\xfc\x93\xe7\x97\x7e\x3f\x7a\xfa\x17\x3c\x30\x2d\xc9\x79\x83\xac\x0d\x5a\x04\x13\x36\xef\x52\xe2\x01\x15\xd7\x9d\xa7\xbd\xc2\xb9\xe6\x43\xc8\x9b\xf1\xbb\xc6\x03\x28\x8f\x20\x25\xca\xe2\x00\x13\x63\xbf\xcf\xef\x1f\xe5\x40\x5d\x1b\x06\x51\x5e\xbe\x92\x1c\x4b\xbb\x31\xc2\xef\x88\xdd\xcc\x7b\xbd\xf9\xe6\xc1\x38\x34\xea\x34\x8c\x6e\xa4\x3b\x9f\x61\xfe\x28\xcb\x79\x73\x45\x12\xac\x45\x37\x8a\x2c\x2f\x7a\xdc\xf9\x30\xd0\x26\x4b\x3f\x6d\x56\xa0\xdb\x5a\x17\xc4\xc8\x97\xd2\xb9\xda\xfe\x67\xd6\xe6\xbd\xb2\xb6\xaf\x05\x1f\x5e\xfd\x7e\x93\x1d\xc8\x7a\x50\x64\xd3\xd9\x6e\x16\x06\x80\x60\xed\x50\xbb\x4c\xe2\xe1\x6f\xd9\xbf\xfe\xf9\xed\xfb\xb7\x7f\xfd\xf9\x9d\x85\xf0\x75\x2c\xf4\x5d\xb0\x24\x5d\x7c\x77\xcc\xea\xb5\x74\x95\xea\xe4\xb3\xb6\x92\x8a\x40\xeb\xb5\x20\x57\xfa\xdd\x2b\x71\x26\xeb\xad\x08\x34\x29\xa7\x34\xfe\xb7\xa1\xa2\xcd\x6c\x70\x22\xc0\xc8\x6d\x2f\x19\xd4\x59\xd1\x39\xe8\xf3\x54\xf8\xe7\xb7\x6f\xff\xf7\xe7\x77\x49\x56\xf0\x01\xcd\xe9\x97\x91\xbf\x8a\xfb\xe4\x2d\x8d\xb3\x99\x69\xd1\x14\x9e\x24\x3a\x9e\xb8\xdd\xba\xea\x1d\x8d\xdb\x50\x3f\xbc\x3a\x7d\xfd\xed\xcb\x7f\xdc\xec\x94\x6f\xff\xf8\x33\x75\xec\x3f\xfe\xfc\x5e\x12\x9a\x60\x37\x4e\xdf\x0f\x84\xc7\x35\x7a\xa4\x31\x70\x74\x65\x30\x4d\x63\x6e\x32\x83\xaa\xb9\xb7\x9d\x76\xe6\x39\x15\xde\xa0\x65\xdd\xde\x59\x6b\xee\xe2\xcb\xdd\x09\x0c\x02\xe7\x3e\xe8\x67\xd0\x92\x44\x2d\xf0\x2c\x72\xf2\xf5\xb7\x2f\xbf\xa7\xc7\x39\x2f\x2b\x4b\xfe\xe1\x35\x7e\xfd\xed\xcb\x1f\xb5\xe1\x03\xd9\x38\xee\x36\x59\x7a\x55\x97\x95\xa8\xea\xcc\xa1\xa9\xe2\x65\xed\xee\x24\xea\xa5\xcd\x4f\xff\x39\x03\xda\xe6\x6c\xaa\x36\xec\x5c\x18\xa8\x2a\x9c\xa8\x38\x16\x2c\x98\xf7\x14\xbd\xd3\x78\x0d\x41\x9c\x66\xca\x6d\x4a\x87\x41\x12\xf7\xca\x23\xa9\x33\x65\x84\x31\xd5\x94\x40\xeb\x36\xdd\x6e\x44\xb0\xce\x5f\xb2\xb0\x0b\x76\xde\x66\xc7\x12\xdc\x11\x25\x67\x78\x88\x2f\x2c\x87\xa5\xd0\x84\xc6\xdc\x64\xda\x96\xc7\x97\xce\x4a\x92\xee\xc0\x72\x36\xd7\x8e\x21\xfb\x43\xcf\x71\x7a\x23\xa4\x8d\x07\x58\xe7\xe5\x89\xf3\x2e\xa2\x70\x44\xa5\x76\x48\xba\x88\xd7\xce\x0d\x86\x69\x77\x62\xaf\x8b\x34\x1d\x1f\xad\x97\x74\x28\xcd\x47\xd4\x22\xc0\x5f\x54\x1c\xd3\x1b\xd4\x4d\x27\x3e\xa4\x4b\x50\x61\x49\xeb\x86\xfc\x26\x4c\xfd\x30\x6f\xfa\x4b\x19\x5f\xfa\x9f\xdb\xde\xdc\x6c\xbb\x85\x72\x6b\xb0\x66\xe3\xb2\x1c\x77\xc6\x83\x39\x6d\xc3\x1c\xc9\xac\xdb\x6d\x2b\xfb\x66\xb6\xa3\xe4\x64\x65\xd6\x0c\xb8\x76\x8f\x38\x7a\x96\x83\x2d\x4f\xcb\xf1\x8c\x62\xa0\xb2\x9a\xa6\xf1\x82\xd0\x87\x54\x4c\x55\x36\x8b\x7d\xbe\x4d\xbc\xb6\x9a\xdd\x96\xf4\xb6\x8b\xbf\x6c\xb7\x7d\xb0\x6a\xcb\x66\xeb\x0c\x90\x26\x5d\x6c\xb7\x1a\xd1\xd9\xed\x65\x9f\x36\x7a\xde\x0b\x85\x69\x76\x4d\x34\xcc\xeb\xbe\xec\x40\xa7\xc1\x50\x8c\x65\x3b\x4a\x7d\x69\x36\x2d\x9d\x40\xe6\x8a\x8a\x1a\xe0\x6b\x66\x83\xb5\x91\xdb\x36\x8e\xb1\xe5\x7c\x1b\x1d\xc7\xbf\x07\xda\xbf\xc7\xa1\xf8\x35\xe5\x87\xf4\x6d\x93\xde\xaa\x72\xc7\x3d\xec\x39\x27\x98\x16\x85\xf3\x88\x91\x3a\xe2\xef\x2c\x9c\x63\x89\x22\x08\xee\xd9\xed\xbc\xdd\xa6\x3d\x57\xe7\xe2\x6f\xfb\x3d\x58\x4b\x70\x3b\x58\x50\xb6\x06\xad\x57\xce\x8e\x46\xd5\x7c\x8a\x97\xcf\xf2\x26\x19\x0e\x00\x1b\xbc\xf1\xa0\x49\x3d\xba\x4d\x5d\x3a\x98\x0a\x37\x46\xfe\xf3\xc1\xd1\x5e\x0b\xbf\x1f\x09\xb8\xec\x79\xa8\xa9\x09\x50\x65\x47\xdc\x0d\x3f\x4b\xae\xf0\x25\x2d\x9a\xf1\x44\xac\xe7\x9e\x6a\x0e\x29\xcf\x0d\x3d\xb6\x4f\x37\x2e\x42\x84\x8b\x83\x8e\xd3\x92\x48\x0b\xe5\x25\x5a\x50\xea\x6f\x4d\x3d\x51\x53\x29\x6a\x02\x2d\xcd\xe8\xcf\x35\x90\x68\xe7\xc4\xc4\xd2\xdb\xad\xd9\xc9\x51\x88\x67\xda\xfd\x4a\xdc\xb7\xc4\x78\xa5\xcf\x50\x17\x50\x90\x2e\x8e\xe4\xca\xcb\xc2\xe7\x0b\xfa\x48\x4c\x64\xd2\x73\x8b\x92\xa2\xa3\xf1\xeb\xd8\x5d\x24\x25\xbe\x03\x52\xf1\xfb\x03\x10\xea\x8e\x62\x66\x73\x61\xad\x48\xe3\xe4\x69\x4e\x58\xf1\x37\xeb\x7a\x21\x53\xb7\x8b\x56\x76\x74\x2e\x20\x56\x56\xe1\xce\x9d\xc7\x95\xb5\xb9\x85\xb6\x76\xdf\x5c\xb9\xa0\xa3\x44\xbe\x76\x17\x5d\x41\x5a\x59\xa1\x3e\x77\x10\x45\x29\x61\xed\x2e\x9a\x84\x59\x59\x95\x1c\x2a\x0a\x8e\x98\x65\x1f\xc6\x75\x4a\xce\xcd\xda\xaa\x4c\xc2\xac\xac\xca\xf5\x83\x63\xae\xf5\xca\xfa\x5c\xc7\x5a\x59\x29\x61\xa8\x6a\x24\xf7\xf7\x2c\xe2\x86\x7e\x3a\xab\xcd\xcd\xe1\x6d\xed\x3e\xbc\x89\xb7\xc9\xaa\x53\xb4\xeb\xb3\xd7\x87\xa6\xfc\x26\x73\xb0\xc7\x1a\xcc\xfc\x65\x37\xad\x2c\xbd\x16\x49\x0c\xbe\x79\x50\xf1\xf6\xd0\x6c\x41\x02\xb5\x33\x18\x94\xe9\xea\x46\xcd\x85\xf2\xc1\xbf\xdc\x64\x60\x4b\xbc\xd6\x8a\xad\xb7\x3a\xf8\xc3\x0b\x6d\xc1\x14\xc1\x69\xc6\x5e\x1e\x1f\x8c\x04\x87\xb1\x7d\x4a\x62\x4b\xdc\x74\x14\xe3\x02\x5d\x83\x41\x2d\x4e\xdd\x5f\x2a\x8b\x7f\xd9\x4c\x82\x83\xda\xb0\x4b\x67\x67\x6d\x38\xe4\xac\x2b\x47\x97\x2a\xe2\x5a\x1d\xba\x2d\xb0\x7f\x59\x09\x22\xb0\x0f\x5d\x14\x9a\x1a\x8d\x9d\xc3\x9c\xb5\x8d\x92\x59\x91\x18\xf5\xc7\xf5\x2a\x8a\x79\x4f\xba\x80\xdc\xa1\x01\x2d\xe3\x3b\x91\x20\x23\x67\x93\xce\x74\x9b\x93\xcd\xac\x15\xfa\x00\x8f\xac\xe4\xe3\x11\x05\x2b\x85\xc2\x85\x8c\x8f\x45\x6c\x8f\x58\x87\x54\x8b\xc2\x1b\xfc\x58\xdc\xbc\xc5\x8f\xd7\x60\x51\x9c\x2d\x87\x01\xeb\x0f\x45\x0f\x4d\x45\xcd\x43\x05\xc7\x0d\x46\x7b\xde\xf5\x61\xa8\x2a\x2c\xc8\xa6\xad\x1c\xeb\x77\xec\xe3\xff\x25\x5d\x7c\x23\x87\xff\xe3\x5d\xf9\xfa\x73\x10\xad\x17\x52\x0f\x0e\x96\xb2\xe4\xeb\x64\x07\x50\x24\xd3\x9e\xe5\xbc\xe6\x26\xcb\xed\x75\x77\xf8\x86\x59\x4e\xf7\x2a\xa8\xb7\xcc\xed\x31\xb2\x28\xb1\x0a\x1c\x49\xce\x5f\x84\x46\xfe\xf3\xe0\xe8\xd2\x21\x54\x9b\x6c\x87\xf1\x70\x1d\xe5\x6c\x79\x5d\x8d\x4d\x66\x50\x03\x59\x1b\x5a\xba\x36\xf4\xef\xf9\x5b\xae\xbb\xdb\x3f\x57\x22\x6b\x7b\x2e\x7f\xf3\x48\x2c\x7d\xf6\xe1\x80\x75\xa4\x8b\xb0\xdb\x4c\xc5\x83\x94\x14\x8b\x7f\x00\xf4\xec\x33\xc5\x0f\x9f\x65\x32\x29\x74\x93\xe5\xa4\xf0\xbf\xbd\xa1\xff\x5a\x8d\x68\x66\x8c\xb2\x0c\xe3\x2c\x51\xbb\x2f\x7f\xbc\xfe\x54\xf6\x65\x29\xe2\x32\x85\x6e\x83\xbf\x8b\xba\xd6\x57\x15\x98\x13\xe7\xf1\x0f\xf7\xda\x78\xab\x07\x65\xdc\x40\x87\x92\x2d\xdd\xd7\x5a\xe6\xef\xad\xfc\x02\x2a\xf3\x07\xe6\x1f\x3d\x4c\xb0\x06\x2a\x67\x13\xfc\x61\x57\xe2\xd1\xe0\x9e\xf4\x13\x03\xf3\x9e\xe8\xf5\x48\x26\xdb\xe2\x4a\x1c\x97\xfd\xcc\xdc\xbb\xe4\x94\xf4\xa9\x55\x39\x96\x27\xb2\xae\xc4\x78\xc1\x31\x45\xd4\xb4\x9a\x7b\x82\xfa\x61\x55\xfa\xe4\x53\x46\xb0\xce\x49\xa2\x76\x14\xaf\x31\x80\x96\xd1\x67\x27\xf7\x49\x3a\xa6\x40\xb7\xe0\x0c\xfd\x63\x90\xcd\x92\x23\x4d\xe9\x89\xa3\x02\x5d\x78\x28\xf0\x93\xa9\x20\xa5\x09\x3a\x58\xeb\xf1\x93\x98\xa5\x96\xba\x43\x51\x6a\x12\xa0\x3e\xdf\xc0\x71\xa7\x1a\xe3\x9f\x25\xa8\xd1\x7c\x96\x19\x1a\xdf\xbf\xa3\x7e\xf3\x64\xe8\xb3\x0c\x15\x07\x61\x4f\x98\x22\x79\x9f\x8b\x6f\x26\x99\x4a\x38\xa0\x63\xc9\xcd\x5a\xfb\xe9\x54\x08\x43\x9d\x74\x11\xf8\x9b\xf8\x16\xa2\x0d\xaa\x7e\x32\x3d\xc2\x16\x57\x32\x69\x44\xf9\x4b\xb4\xf8\xff\x00\x00\x00\xff\xff\x55\x07\xdf\x8e\x80\x65\x00\x00")

func redirectsRedirectsCsvBytes() ([]byte, error) {
	return bindataRead(
		_redirectsRedirectsCsv,
		"redirects/redirects.csv",
	)
}

func redirectsRedirectsCsv() (*asset, error) {
	bytes, err := redirectsRedirectsCsvBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "redirects/redirects.csv", size: 25984, mode: os.FileMode(420), modTime: time.Unix(1578062689, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
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

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
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
	"templates/error.tmpl": templatesErrorTmpl,
	"templates/main.tmpl": templatesMainTmpl,
	"templates/partials/footer.tmpl": templatesPartialsFooterTmpl,
	"templates/partials/header.tmpl": templatesPartialsHeaderTmpl,
	"redirects/redirects.csv": redirectsRedirectsCsv,
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
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
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
	"redirects": &bintree{nil, map[string]*bintree{
		"redirects.csv": &bintree{redirectsRedirectsCsv, map[string]*bintree{}},
	}},
	"templates": &bintree{nil, map[string]*bintree{
		"error.tmpl": &bintree{templatesErrorTmpl, map[string]*bintree{}},
		"main.tmpl": &bintree{templatesMainTmpl, map[string]*bintree{}},
		"partials": &bintree{nil, map[string]*bintree{
			"footer.tmpl": &bintree{templatesPartialsFooterTmpl, map[string]*bintree{}},
			"header.tmpl": &bintree{templatesPartialsHeaderTmpl, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
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
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
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
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

