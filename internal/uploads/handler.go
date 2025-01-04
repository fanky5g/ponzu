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

		// StoreFiles has the SetUpload call (which is equivalent of CreateContent in other controllers)
		files, err := request.GetRequestFiles(req)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to get request files")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		savedFiles, err := uploadService.UploadFiles(files)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to save files")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, savedFile := range savedFiles {
			req.PostForm.Set(savedFile.Name, savedFile.Path)
		}

		if len(savedFiles) == 1 {
			response.Respond(
				res,
				req,
				response.NewRedirectResponse(publicPath, "/edit/upload?type=upload&id="+savedFiles[0].ID),
			)
			return
		}

		response.Respond(res, req, response.NewRedirectResponse(publicPath, "/uploads"))
	}
}
