package content

import (
	"github.com/fanky5g/ponzu/internal/http/response"
	"net/http"

	"github.com/fanky5g/ponzu/internal/http/request"
	log "github.com/sirupsen/logrus"
)

func NewDeleteHandler(publicPath string, contentService *Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		selectedItems, err := request.GetSelectedItems(req)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to parse request")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(selectedItems) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		t := req.FormValue("type")

		if len(selectedItems) == 0 || t == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err = contentService.DeleteContent(t, selectedItems...)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to delete content")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Respond(
			res,
			req,
			response.NewRedirectResponse(publicPath, "/contents?type="+t),
		)
	}
}
