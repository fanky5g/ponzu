package content

import (
	"errors"
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/http/response"
	"log"
	"net/http"
)

func NewAPIDeleteContentHandler(contentTypes map[string]content.Builder, contentService *Service) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		isSlug, identifier := request.GetRequestContentId(req)
		if identifier == "" {
			response.Respond(
				res,
				req,
				response.NewJSONResponse(
					http.StatusBadRequest,
					nil,
					errors.New("entities id is required"),
				),
			)
			return
		}

		if isSlug {
			response.Respond(
				res,
				req,
				response.NewJSONResponse(
					http.StatusBadRequest,
					nil,
					errors.New("slug not supported for delete"),
				),
			)
			return
		}

		t := req.URL.Query().Get("type")
		if t == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		_, found := contentTypes[t]
		if !found {
			log.Println("[Delete] attempt to delete entities of unknown type:", t, "from:", req.RemoteAddr)
			res.WriteHeader(http.StatusNotFound)
			return
		}

		err := contentService.DeleteContent(t, identifier)
		if err != nil {
			log.Printf("[Delete] error: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Respond(
			res,
			req,
			response.NewJSONResponse(
				http.StatusOK,
				map[string]interface{}{
					"id":     identifier,
					"status": "deleted",
					"type":   t,
				},
				nil,
			),
		)
	}
}
