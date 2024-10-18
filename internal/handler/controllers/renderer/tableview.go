package renderer

import (
	"bytes"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/fanky5g/ponzu/internal/handler/controllers/resources/viewparams/table"
)

func (r *renderer) TableView(
	w http.ResponseWriter,
	templateName string,
	data *table.Table) {
	data.PublicPath = r.ctx.Paths().PublicPath
	buf := &bytes.Buffer{}
	tableViewTmpl := r.TemplateFromDir(templateName)

	if err := tableViewTmpl.Execute(buf, data); err != nil {
		log.WithField("Error", err).Warning("Failed to render table")
		r.InternalServerError(w)
		return
	}

	r.InjectInAdminView(w, buf)
}
