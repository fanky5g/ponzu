package types

import "github.com/fanky5g/ponzu/content/item"

type Story struct {
	item.Item

	Title  string `json:"title"`
	Body   string `json:"body"`
	Author Author `json:"author"`
}
