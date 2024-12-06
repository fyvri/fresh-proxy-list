package utils

import (
	"net/url"
)

type URLParserUtil struct{}

type URLParserUtilInterface interface {
	Parse(rawURL string) (*url.URL, error)
}

func NewURLParser() URLParserUtilInterface {
	return &URLParserUtil{}
}

func (u *URLParserUtil) Parse(rawURL string) (*url.URL, error) {
	return url.Parse(rawURL)
}
