package content

import (
	"github.com/fanky5g/ponzu/internal/http/response"
	"net/http"

	"github.com/fanky5g/ponzu/internal/http/request"
	log "github.com/sirupsen/logrus"
)

func NewDeleteUploadHandler(publicPath string, uploadService *UploadService) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		selectedItems, err := request.GetSelectedItems(req)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to parse request")
			res.WriteHeader(http.StatusInternalServerError)
		}

		if len(selectedItems) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err = uploadService.DeleteUpload(selectedItems...)
		if err != nil {
			log.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Respond(
			res,
			req,
			response.NewRedirectResponse(publicPath, "/uploads"),
		)
	}
}
