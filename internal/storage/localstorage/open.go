package localstorage

import (
	"net/http"
)

func (c *Client) Open(name string) (http.File, error) {
	return c.fs.Open(name)
}
