package request

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

func GetSelectedItems(req *http.Request) ([]string, error) {
	err := req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse form")
	}

	ids := make([]string, 0)

	idParam := strings.TrimSpace(req.FormValue("id"))
	idsParam := strings.TrimSpace(req.FormValue("ids"))
	if idParam != "" {
		ids = append(ids, idParam)
	} else if idsParam != "" {
		idsToDelete := strings.FieldsFunc(idsParam, func(c rune) bool {
			return c == ','
		})

		for _, idToDelete := range idsToDelete {
			ids = append(ids, strings.TrimSpace(idToDelete))
		}
	}

	return ids, nil
}
