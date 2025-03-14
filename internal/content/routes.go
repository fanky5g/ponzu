package content

import (
	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/internal/layouts"
	"github.com/fanky5g/ponzu/internal/search"
	"github.com/fanky5g/ponzu/internal/templates"
	"net/http"
)

type RouterInterface interface {
	AuthorizedRoute(pattern string, handler func() http.HandlerFunc)
	APIAuthorizedRoute(pattern string, handler func() http.HandlerFunc)
}

func RegisterRoutes(
	r RouterInterface,
	publicPath string,
	contentTypes map[string]content.Builder,
	contentService *Service,
	uploadService *UploadService,
	searchService *search.Service,
	dashboardLayout layouts.Template) error {
	mapper := NewMapper()

	templateNames, templateNameMatchErr := templates.GlobNames("views/datatable/*.gohtml")
	if templateNameMatchErr != nil {
		return templateNameMatchErr
	}

	dataTable, templateErr := dashboardLayout.Child(templateNames...)
	if templateErr != nil {
		return templateErr
	}

	uploadsTableTemplateNames, templateNameMatchErr := templates.GlobNames("views/uploadsdatatable/*.gohtml")
	if templateNameMatchErr != nil {
		return templateNameMatchErr
	}

	uploadsTable, templateErr := dashboardLayout.Child(uploadsTableTemplateNames...)
	if templateErr != nil {
		return templateErr
	}

	r.AuthorizedRoute("GET /edit", func() http.HandlerFunc {
		return NewEditContentFormHandler(contentService, publicPath, dashboardLayout)
	})

	r.AuthorizedRoute("POST /edit", func() http.HandlerFunc {
		return NewSaveContentHandler(contentService, publicPath)
	})

	r.AuthorizedRoute("POST /edit/workflow", func() http.HandlerFunc {
		return NewContentWorkflowTransitionHandler(contentService, publicPath, dashboardLayout)
	})

	r.AuthorizedRoute("GET /edit/upload", func() http.HandlerFunc {
		return NewEditUploadFormHandler(publicPath, uploadService, dashboardLayout)
	})

	r.AuthorizedRoute("POST /edit/upload", func() http.HandlerFunc {
		return NewSaveUploadHandler(uploadService, publicPath)
	})

	r.APIAuthorizedRoute("GET /api/references", func() http.HandlerFunc {
		return NewListReferencesHandler(contentService, mapper)
	})

	r.APIAuthorizedRoute("GET /api/references/{id}", func() http.HandlerFunc {
		return NewGetReferenceHandler(contentService, mapper)
	})

	r.AuthorizedRoute("GET /uploads", func() http.HandlerFunc {
		return NewUploadContentsHandler(publicPath, uploadService, uploadsTable)
	})

	r.AuthorizedRoute("GET /uploads/search", func() http.HandlerFunc {
		return NewUploadSearchHandler(publicPath, searchService, uploadsTable)
	})

	r.AuthorizedRoute("/contents", func() http.HandlerFunc {
		return NewContentsHandler(publicPath, contentService, contentTypes, dataTable)
	})

	r.AuthorizedRoute("/contents/search", func() http.HandlerFunc {
		return NewSearchContentHandler(contentTypes, searchService)
	})

	r.AuthorizedRoute("GET /contents/export", func() http.HandlerFunc {
		return NewExportHandler(contentService)
	})

	r.AuthorizedRoute("GET /edit/delete", func() http.HandlerFunc {
		return NewDeleteHandler(publicPath, contentService)
	})

	r.AuthorizedRoute("GET /edit/upload/delete", func() http.HandlerFunc {
		return NewDeleteUploadHandler(publicPath, uploadService)
	})

	r.APIAuthorizedRoute("POST /api/content", func() http.HandlerFunc {
		return NewAPIContentHandler(contentService, uploadService, contentTypes)
	})

	r.APIAuthorizedRoute("/api/search", func() http.HandlerFunc {
		return NewSearchContentHandler(contentTypes, searchService)
	})

	return nil
}
