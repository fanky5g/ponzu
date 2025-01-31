package types

import "github.com/fanky5g/ponzu/content/item"

type Product struct {
	item.Item

	Name        string   `json:"name"`
	Description string   `json:"description"`
	Images      []string `json:"images" reference:"Photo"`
}
