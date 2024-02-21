package assets

import (
	_ "embed"
	"errors"
)

//go:embed redirects/redirects.csv
var redirectsRedirectsCSV []byte

func Asset(s string) ([]byte, error) {
	if s == "redirects/redirects.csv" {
		return redirectsRedirectsCSV, nil
	}
	return nil, errors.New("no such asset found")
}
