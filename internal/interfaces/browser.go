package interfaces

type BrowserManager interface {
	OpenURL(url string) error
}
