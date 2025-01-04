package uploads

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/fanky5g/ponzu/content/entities"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/internal/config"
	"github.com/fanky5g/ponzu/internal/constants"
	"github.com/fanky5g/ponzu/internal/content"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/http/response"
)

func NewEditUploadFormHandler(
	uploadsService *Service,
	contentService *content.Service,
	cfg config.ConfigCache,
	publicPath string,
) http.HandlerFunc {
	tmpl, err := getEditUploadTemplate()
	if err != nil {
		log.Fatalf("Failed to build page template: %v", err)
	}

	return func(res http.ResponseWriter, req *http.Request) {
		q := req.URL.Query()
		i := q.Get("id")

		var fileUpload *entities.Upload
		var err error
		if i != "" {
			fileUpload, err = uploadsService.GetUpload(i)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to get file upload")
				return
			}

			if fileUpload == nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			_, ok := interface{}(fileUpload).(item.Identifiable)
			if !ok {
				log.Println("Content type", constants.UploadEntityName, "doesn't implement item.Identifiable")
				return
			}

			fileUpload = &entities.Upload{}
		}

		editUploadForm, err := NewEditUploadFormViewModel(fileUpload, cfg, publicPath, contentService.ContentTypes())
		if err != nil {
			// TODO: handle error
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Respond(
			res,
			req,
			response.NewHTMLResponse(http.StatusOK, tmpl, editUploadForm),
		)

	}
}

func NewSaveUploadHandler(uploadService *Service, publicPath string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		err := req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
		if err != nil {
			log.WithField("Error", err).Warning("Failed to parse form")
			// TODO: handle error
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		t := req.FormValue("type")
		post, err := request.GetFileUploadFromFormData(req.Form)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to get form file")
			// TODO: handle error
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		hook, ok := post.(item.Hookable)
		if !ok {
			log.Println("Type", t, "does not implement item.Hookable or embed item.Item.")
			// TODO: handle error
			res.WriteHeader(http.StatusBadRequest)
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
			// TODO: handle error
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		urlPaths, err := uploadService.UploadFiles(files)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to save files")
			// TODO: handle error
			res.WriteHeader(http.StatusInternalServerError)
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

		response.Respond(res, req, response.NewRedirectResponse(publicPath, "/uploads"))
	}
}
