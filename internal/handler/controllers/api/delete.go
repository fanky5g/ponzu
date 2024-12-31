package api

import (
	"errors"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/tokens"
	"log"
	"net/http"
)

func NewDeleteContentHandler(r router.Router) http.HandlerFunc {
	contentTypes := r.Context().Types().Content
	contentService := r.Context().Service(tokens.ContentServiceToken).(*content.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		isSlug, identifier := request.GetRequestContentId(req)
		if identifier == "" {
			r.Renderer().Error(res, http.StatusBadRequest, errors.New("entities id is required"))
			return
		}

		if isSlug {
			r.Renderer().Error(res, http.StatusBadRequest, errors.New("slug not supported for delete"))
		}

		t := req.URL.Query().Get("type")
		if t == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		p, found := contentTypes[t]
		if !found {
			log.Println("[Delete] attempt to delete entities of unknown type:", t, "from:", req.RemoteAddr)
			res.WriteHeader(http.StatusNotFound)
			return
		}

		hook, ok := p().(item.Hookable)
		if !ok {
			log.Println("[Delete] error: Type", t, "does not implement item.Hookable or embed item.Item.")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err := hook.BeforeAPIDelete(res, req)
		if err != nil {
			log.Println("[Delete] error calling BeforeAPIDelete:", err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		err = hook.BeforeDelete(res, req)
		if err != nil {
			log.Println("[Delete] error calling BeforeDelete:", err)
			return
		}

		err = contentService.DeleteContent(t, identifier)
		if err != nil {
			log.Printf("[Delete] error: %v\n", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = hook.AfterDelete(res, req)
		if err != nil {
			log.Println("[Delete] error calling AfterDelete:", err)
			return
		}

		err = hook.AfterAPIDelete(res, req)
		if err != nil {
			log.Println("[Delete] error calling AfterAPIDelete:", err)
			return
		}

		r.Renderer().Json(res, http.StatusOK, map[string]interface{}{
			"id":     identifier,
			"status": "deleted",
			"type":   t,
		})
	}
}
