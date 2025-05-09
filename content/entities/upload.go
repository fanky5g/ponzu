package entities

import (
	"fmt"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/internal/constants"
	"github.com/fanky5g/ponzu/util"
)

const UploadRepositoryToken = "uploads"

// Upload represents the file uploaded to the system
type Upload struct {
	item.Item

	Name          string `json:"name"`
	Path          string `json:"path"`
	ContentLength int64  `json:"content_length"`
	ContentType   string `json:"content_type"`
}

func (*Upload) EntityName() string {
	return constants.UploadEntityName
}

func (f *Upload) GetTitle() string { return f.Name }

func (*Upload) GetRepositoryToken() string {
	return UploadRepositoryToken
}

// MarshalEditor writes a buffer of templates to edit a Post and partially implements editor.Editable
func (f *Upload) MarshalEditor(publicPath string) ([]byte, error) {
	isEmptyFile := f.Path == ""
	formLabel := "Edit Upload"
	if isEmptyFile {
		formLabel = "Upload New File"
	}

	return editor.Form(f,
		editor.Field{
			View: editor.File(publicPath, "Path", f, map[string]string{
				"label":       formLabel,
				"placeholder": "Upload the file here",
				"PublicPath":  publicPath,
			}),
		},
		editor.Field{
			View: func() []byte {
				if isEmptyFile {
					return nil
				}

				return []byte(`
            <div class="control-block file-attributes">
				<label>` + f.Name + `</label>
				<ul class="mdc-list mdc-list--two-line">
				  <li class="mdc-list-item" tabindex="0">
					<span class="mdc-list-item__ripple"></span>
					<span class="mdc-list-item__text">
					  <span class="mdc-list-item__primary-text">Content-Length</span>
					  <span class="mdc-list-item__secondary-text">` + fmt.Sprintf("%s", util.FmtBytes(float64(f.ContentLength))) + `</span>
					</span>
				  </li>
				  <li class="mdc-list-item">
					<span class="mdc-list-item__ripple"></span>
					<span class="mdc-list-item__text">
					  <span class="mdc-list-item__primary-text">Content-Type</span>
					  <span class="mdc-list-item__secondary-text">` + f.ContentType + `</span>
					</span>
				  </li>
				  <li class="mdc-list-item">
					<span class="mdc-list-item__ripple"></span>
					<span class="mdc-list-item__text">
					  <span class="mdc-list-item__primary-text">Uploaded</span>
					  <span class="mdc-list-item__secondary-text">` + util.FmtTime(f.Timestamp) + `</span>
					</span>
				  </li>
				</ul>
            </div>
            `)
			}(),
		},
	)
}

func (f *Upload) Push() []string {
	return []string{
		"path",
	}
}

// IndexContent determines if FileUpload should be indexed for searching
func (f *Upload) IndexContent() bool {
	return true
}
