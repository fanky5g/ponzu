package content

import (
	"errors"
	"github.com/fanky5g/ponzu/content/workflow"
	"github.com/fanky5g/ponzu/exceptions"
	"github.com/fanky5g/ponzu/internal/http/request"
	"github.com/fanky5g/ponzu/internal/http/response"
	"github.com/fanky5g/ponzu/internal/layouts"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func NewEditContentFormHandler(contentService *Service, publicPath string, layout layouts.Template) http.HandlerFunc {
	tmpl, templateErr := layout.Child("views/edit_content_view.gohtml")
	if templateErr != nil {
		log.Fatalf("Failed to build page template: %v", templateErr)
	}

	return func(res http.ResponseWriter, req *http.Request) {
		contentQuery, err := MapContentQueryFromRequest(req)
		if err != nil {
			log.WithField("Error", err).Info("Error mapping request body")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		entity, err := contentService.Type(contentQuery.Type)
		if err != nil {
			log.WithField("Error", err).Warning("Error fetching content type")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		if contentQuery.ID != "" {
			entity, err = contentService.GetContent(contentQuery.Type, contentQuery.ID)
			if err != nil {
				log.WithField("Error", err).Warning("Error getting content by ID")
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			if entity == nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		editContentForm, err := NewEditContentFormViewModel(
			entity,
			publicPath,
			nil,
		)
		if err != nil {
			log.WithField("Error", err).Warning("Error building ContentFormViewModel")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.Respond(
			res,
			req,
			response.NewHTMLResponse(http.StatusOK, tmpl, editContentForm),
		)
	}
}

func NewSaveContentHandler(contentService *Service, publicPath string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		contentQuery, err := MapContentQueryFromRequest(req)
		if err != nil {
			// TODO: handle error
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		entity, err := contentService.Type(contentQuery.Type)
		if err != nil {
			// TODO: handle error
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		entity, err = request.GetEntityFromFormData(entity, req.PostForm)
		if err != nil {
			log.WithField("Error", err).Warning("Failed to map request entity")
			res.WriteHeader(http.StatusBadRequest)
			return
		}

		if contentQuery.ID == "" {
			contentQuery.ID, err = contentService.CreateContent(contentQuery.Type, entity)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to create content")
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

		} else {
			_, err = contentService.UpdateContent(contentQuery.Type, contentQuery.ID, entity)
			if err != nil {
				log.WithField("Error", err).Warning("Failed to update content")
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		response.Respond(
			res,
			req,
			response.NewRedirectResponse(
				publicPath,
				req.URL.Path+"?type="+contentQuery.Type+"&id="+contentQuery.ID,
			),
		)
	}
}

func NewContentWorkflowTransitionHandler(contentService *Service, publicPath string, layout layouts.Template) http.HandlerFunc {
	tmpl, templateErr := layout.Child("views/edit_content_view.gohtml")
	if templateErr != nil {
		log.Fatalf("Failed to build page template: %v", templateErr)
	}

	return func(res http.ResponseWriter, req *http.Request) {
		contentTransitionInput, err := MapContentTransitionInputFromRequest(req)
		if err != nil {
			exceptions.Log(err)
			res.WriteHeader(getResponseCode(err))
			return
		}

		entity, workflowTransitionErr := contentService.TransitionWorkflowState(
			contentTransitionInput.Type,
			contentTransitionInput.ID,
			workflow.State(contentTransitionInput.TargetState),
		)
		if workflowTransitionErr != nil {
			exceptions.Log(err)
		}

		editContentForm, err := NewEditContentFormViewModel(
			entity,
			publicPath,
			workflowTransitionErr,
		)
		if err != nil {
			exceptions.Log(err)
			res.WriteHeader(getResponseCode(err))
			return
		}

		response.Respond(
			res,
			req,
			response.NewHTMLResponse(getResponseCode(workflowTransitionErr), tmpl, editContentForm),
		)
	}
}

func getResponseCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	var clientException *exceptions.ClientException
	if errors.As(err, &clientException) {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}
