package controllers

import (
	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/entities"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/handler/controllers/router"
	"github.com/fanky5g/ponzu/internal/services/storage"
	"github.com/fanky5g/ponzu/tokens"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func NewEditUploadHandler(r router.Router) http.HandlerFunc {
	storageService := r.Context().Service(tokens.StorageServiceToken).(storage.Service)

	return func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			q := req.URL.Query()
			i := q.Get("id")

			var fileUpload *entities.FileUpload
			var err error
			if i != "" {
				fileUpload, err = storageService.GetFileUpload(i)
				if err != nil {
					log.WithField("Error", err).Warning("Failed to get file upload")
					return
				}

				if fileUpload == nil {
					r.Renderer().BadRequest(res)
					return
				}
			} else {
				_, ok := interface{}(fileUpload).(item.Identifiable)
				if !ok {
					log.Println("Content type", constants.UploadsEntityName, "doesn't implement item.Identifiable")
					return
				}

				fileUpload = &entities.FileUpload{}
			}

			r.Renderer().ManageEditable(res, interface{}(fileUpload).(editor.Editable), constants.UploadsEntityName)
		case http.MethodPost:
			err := req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
			if err != nil {
				log.WithField("Error", err).Warning("Failed to parse form")
				r.Renderer().InternalServerError(res)
				return
			}

			t := req.FormValue("type")
			post, err := request.GetFileUploadFromFormData(req.Form)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to get form file")
				return
			}

			hook, ok := post.(item.Hookable)
			if !ok {
				log.Println("Type", t, "does not implement item.Hookable or embed item.Item.")
				r.Renderer().BadRequest(res)
				return
			}

			err = hook.BeforeSave(res, req)
			if err != nil {
				log.WithField("Error", err).Warningf("Error running BeforeSave method in editHandler for: %s", t)
				return
			}

			// StoreFiles has the SetUpload call (which is equivalent of CreateContent in other controllers)
			files, err := request.GetRequestFiles(req)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to get request files")
				r.Renderer().InternalServerError(res)
				return
			}

			urlPaths, err := storageService.StoreFiles(files)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to save files")
				r.Renderer().InternalServerError(res)
				return
			}

			for name, urlPath := range urlPaths {
				req.PostForm.Set(name, urlPath)
			}

			err = hook.AfterSave(res, req)
			if err != nil {
				log.WithField("Error", err).
					Warningf("Error running AfterSave method in editHandler for: %s", t)
				return
			}

			r.Redirect(req, res, "/uploads")

		default:
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}
