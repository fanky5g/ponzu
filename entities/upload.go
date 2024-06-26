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
	f.Path = filepath.Join(paths.PublicPath, f.Path)

	view, err := editor.Form(f,
		paths,
		editor.Field{
			View: func() []byte {
				if f.Path == "" {
					return nil
				}

				return []byte(`
            <div class="input-field col s12">
				<h5>` + f.Name + `</h5>
				<ul>
					<li><span class="grey-text text-lighten-1">Content-Length:</span> ` + fmt.Sprintf("%s", FmtBytes(float64(f.ContentLength))) + `</li>
					<li><span class="grey-text text-lighten-1">Content-Type:</span> ` + f.ContentType + `</li>
					<li><span class="grey-text text-lighten-1">Uploaded:</span> ` + FmtTime(f.Timestamp) + `</li>
				</ul>
            </div>
            `)
			}(),
		},
		editor.Field{
			View: editor.File("Path", f, map[string]string{
				"label":       "File Upload",
				"placeholder": "Upload the file here",
			}),
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
			
			// hide default fields & labels unnecessary for the config
			var fields = $('.default-fields');
			fields.css('position', 'relative');
			fields.find('input:not([type=submit])').remove();
			fields.find('label').remove();
			fields.find('button').css({
				position: 'absolute',
				top: '-10px',
				right: '0px'
			});

			var contentOnly = $('.entities-only.__ponzu');
			contentOnly.hide();
			contentOnly.find('input, textarea, select').attr('name', '');

			// adjust layout of td so save button is in same location as usual
			fields.find('td').css('float', 'right');

			// stop some fixed config settings from being modified
			fields.find('input[name=client_secret]').attr('name', '');

			// hide save, show delete
			if ($('h5').length > 0) {
				fields.find('.save-post').hide();
				fields.find('.delete-post').show();
			} else {
				fields.find('.save-post').show();
				fields.find('.delete-post').hide();
			}
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
