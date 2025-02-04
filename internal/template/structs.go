package template

import "github.com/a-h/templ"

type Tab struct {
	Name    string
	Content templ.Component
	URL     string
}
