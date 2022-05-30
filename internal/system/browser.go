package system

import (
	"github.com/pkg/browser"
)

type Browser struct{}

func NewBrowser() Browser {
	return Browser{}
}

func (Browser) OpenURL(url string) error {
	return browser.OpenURL(url)
}
