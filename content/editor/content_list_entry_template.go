package editor

import (
	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/content/item"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func BuildContentListEntryTemplate(e Editable, typeName string) string {
	s, ok := e.(item.Sortable)
	if !ok {
		log.Warnf("Content type %s doesn't implement item.Sortable", typeName)
		return `<li class="col s12">Error retrieving data. Your data type doesn't implement necessary interfaces. (item.Sortable)</li>`
	}

	i, ok := e.(item.Identifiable)
	if !ok {
		log.Warnf("Content type %s doesn't implement item.Identifiable", typeName)
		return `<li class="col s12">Error retrieving data. Your data type doesn't implement necessary interfaces. (item.Identifiable)</li>`
	}

	// use sort to get other info to display in controllers UI post list
	tsTime := time.Unix(s.Time()/1000, 0)
	upTime := time.Unix(s.Touch()/1000, 0)
	updatedTime := upTime.Format("01/02/06 03:04 PM")
	publishTime := tsTime.Format("01/02/06")

	cid := i.ItemID()

	action := "{{ .PublicPath }}/edit/delete"
	link := `<a href="{{ .PublicPath }}/edit?type=` + typeName + `&id=` + cid + `">` + i.ItemID() + `</a>`
	if strings.HasPrefix(typeName, constants.UploadsEntityName) {
		link = `<a href="{{ .PublicPath }}/edit/upload?id=` + cid + `">` + i.ItemID() + `</a>`
		action = "{{ .PublicPath }}/edit/upload/delete"
	}

	return `
			<li class="col s12">
				` + link + `
				<span class="post-detail">Updated: ` + updatedTime + `</span>
				<span class="publish-date right">` + publishTime + `</span>

				<form enctype="multipart/form-data" class="quick-delete-post __ponzu right" action="` + action + `" method="post">
					<span>Delete</span>
					<input type="hidden" name="id" value="` + cid + `" />
					<input type="hidden" name="type" value="` + typeName + `" />
				</form>
			</li>`
}
