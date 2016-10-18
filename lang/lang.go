package lang

import (
	"net/http"
	"strings"
)

type Locale string

const (
	EN Locale = "en"
	CY Locale = "cy"
)

func Get(req *http.Request) Locale {
	if strings.HasPrefix(strings.ToLower(req.URL.Host), "cy.") {
		return CY
	}
	return EN
}
