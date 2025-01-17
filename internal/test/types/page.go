package types

import "github.com/fanky5g/ponzu/content/item"

type Page struct {
	item.Item

	Title         string             `json:"title"`
	URL           string             `json:"url"`
	ContentBlocks *PageContentBlocks `json:"content_blocks"`
}
