package controllers

import (
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/content"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func NewDeleteHandler(r router.Router) http.HandlerFunc {
	contentService := r.Context().Service(tokens.ContentServiceToken).(content.Service)
	contentTypes := r.Context().Types().Content

	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		err := req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
		if err != nil {
			log.WithField("Error", err).Warning("Failed to parse form")
			r.Renderer().InternalServerError(res)
			return
		}

		id := req.FormValue("id")
		t := req.FormValue("type")
		ct := t

		if id == "" || t == "" {
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		p, ok := contentTypes[ct]
		if !ok {
			log.Println("Type", t, "does not implement item.Hookable or embed item.Item.")
			r.Renderer().BadRequest(res)
			return
		}

		post := p()
		hook, ok := post.(item.Hookable)
		if !ok {
			log.Println("Type", t, "does not implement item.Hookable or embed item.Item.")
			r.Renderer().BadRequest(res)
			return
		}

		post, err = contentService.GetContent(t, id)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to get content")
			r.Renderer().InternalServerError(res)
			return
		}

		reject := req.URL.Query().Get("reject")
		if reject == "true" {
			err = hook.BeforeReject(res, req)
			if err != nil {
				log.Println("Error running BeforeReject method in deleteHandler for:", t, err)
				return
			}
		}

		err = hook.BeforeAdminDelete(res, req)
		if err != nil {
			log.Println("Error running BeforeAdminDelete method in deleteHandler for:", t, err)
			return
		}

		err = hook.BeforeDelete(res, req)
		if err != nil {
			log.Println("Error running BeforeDelete method in deleteHandler for:", t, err)
			return
		}

		err = contentService.DeleteContent(t, id)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to delete content")
			r.Renderer().InternalServerError(res)
			return
		}

		err = hook.AfterDelete(res, req)
		if err != nil {
			log.Println("Error running AfterDelete method in deleteHandler for:", t, err)
			return
		}

		err = hook.AfterAdminDelete(res, req)
		if err != nil {
			log.Println("Error running AfterDelete method in deleteHandler for:", t, err)
			return
		}

		if reject == "true" {
			err = hook.AfterReject(res, req)
			if err != nil {
				log.Println("Error running AfterReject method in deleteHandler for:", t, err)
				return
			}
		}

		r.Redirect(req, res, "/contents?type="+ct)
	}
}
