package entities

import (
	"fmt"
	"github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/constants"
	"github.com/fanky5g/ponzu/content/editor"
	"github.com/fanky5g/ponzu/content/item"
	"github.com/fanky5g/ponzu/tokens"
	"path/filepath"
	"time"
)

// FileUpload represents the file uploaded to the system
type FileUpload struct {
	item.Item

	Name          string `json:"name"`
	Path          string `json:"path"`
	ContentLength int64  `json:"content_length"`
	ContentType   string `json:"content_type"`
}

func (*FileUpload) EntityName() string {
	return constants.UploadsEntityName
}

func (f *FileUpload) GetTitle() string { return f.Name }

func (*FileUpload) GetRepositoryToken() tokens.RepositoryToken {
	return tokens.UploadRepositoryToken
}

// MarshalEditor writes a buffer of templates to edit a Post and partially implements editor.Editable
func (f *FileUpload) MarshalEditor(paths config.Paths) ([]byte, error) {
	isEmptyFile := f.Path == ""
	f.Path = filepath.Join(paths.PublicPath, f.Path)
	formLabel := "Edit Upload"
	if isEmptyFile {
		formLabel = "Upload New File"
	}

	view, err := editor.Form(f,
		paths,
		editor.Field{
			View: editor.File("Path", f, map[string]string{
				"label":       formLabel,
				"placeholder": "Upload the file here",
				"PublicPath":  paths.PublicPath,
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
					  <span class="mdc-list-item__secondary-text">` + fmt.Sprintf("%s", FmtBytes(float64(f.ContentLength))) + `</span>
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
					  <span class="mdc-list-item__secondary-text">` + FmtTime(f.Timestamp) + `</span>
					</span>
				  </li>
				</ul>
            </div>
            `)
			}(),
		},
	)

	if err != nil {
		return nil, err
	}

	script := []byte(`
	<script>
		$(function() {
			// change form action to storage-specific endpoint
			var form = $('form');
			form.attr('action', '` + paths.PublicPath + `/edit/upload');
		});
	</script>
	`)

	view = append(view, script...)

	return view, nil
}

func (f *FileUpload) Push() []string {
	return []string{
		"path",
	}
}

// FmtBytes converts the numeric byte size value to the appropriate magnitude
// size in KB, MB, GB, TB, PB, or EB.
func FmtBytes(size float64) string {
	unit := float64(1024)
	BYTE := unit
	KBYTE := BYTE * unit
	MBYTE := KBYTE * unit
	GBYTE := MBYTE * unit
	TBYTE := GBYTE * unit
	PBYTE := TBYTE * unit

	switch {
	case size < BYTE:
		return fmt.Sprintf("%0.f B", size)
	case size < KBYTE:
		return fmt.Sprintf("%.1f KB", size/BYTE)
	case size < MBYTE:
		return fmt.Sprintf("%.1f MB", size/KBYTE)
	case size < GBYTE:
		return fmt.Sprintf("%.1f GB", size/MBYTE)
	case size < TBYTE:
		return fmt.Sprintf("%.1f TB", size/GBYTE)
	case size < PBYTE:
		return fmt.Sprintf("%.1f PB", size/TBYTE)
	default:
		return fmt.Sprintf("%0.f B", size)
	}

}

// FmtTime shows a human-readable time based on the timestamp
func FmtTime(t int64) string {
	return time.Unix(t, 0).Format("03:04 PM Jan 2, 2006") + " (UTC)"
}

// IndexContent determines if FileUpload should be indexed for searching
func (f *FileUpload) IndexContent() bool {
	return true
}
