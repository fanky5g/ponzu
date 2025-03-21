package content

import (
	"github.com/fanky5g/ponzu/content/entities"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/internal/constants"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/http/response"
	"github.com/fanky5g/ponzu/internal/layouts"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func NewEditUploadFormHandler(publicPath string, uploadsService *UploadService, layout layouts.Template) http.HandlerFunc {
	tmpl, tmplErr := layout.Child("views/edit_upload_view.gohtml")
	if tmplErr != nil {
		panic(tmplErr)
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

		var viewModel *EditUploadFormViewModel
		viewModel, err = NewEditUploadFormViewModel(publicPath, fileUpload)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to build view model")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Respond(
			res,
			req,
			response.NewHTMLResponse(http.StatusOK, tmpl, viewModel),
		)
	}
}

func NewSaveUploadHandler(uploadService *UploadService, publicPath string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		err := req.ParseMultipartForm(1024 * 1024 * 4) // maxMemory 4MB
		if err != nil {
			log.WithField("Error", err).Warning("Failed to parse form")
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
