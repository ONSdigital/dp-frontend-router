package assets

import (
	_ "embed"
	"errors"
)

//go:embed redirects/redirects.csv
var redirectsRedirectsCSV []byte

const assetPath = "redirects/redirects.csv"

func Asset(s string) ([]byte, error) {
	if s == assetPath {
		return redirectsRedirectsCSV, nil
	}
	return nil, errors.New("no such asset found")
}
