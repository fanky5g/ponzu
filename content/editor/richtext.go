package editor

import (
	"bytes"
	"html"
)

// Richtext returns the []byte of a rich text editor (provided by http://summernote.org/) with a label.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func Richtext(fieldName string, p interface{}, attrs map[string]string, args *FieldArgs) []byte {
	// create wrapper for richtext editor, which isolates the editor's css
	iso := []byte(`<div class="iso-texteditor input-field col s12"><label>` + attrs["label"] + `</label>`)
	isoClose := []byte(`</div>`)

	if _, ok := attrs["class"]; ok {
		attrs["class"] += "richtext " + fieldName
	} else {
		attrs["class"] = "richtext " + fieldName
	}

	if _, ok := attrs["id"]; ok {
		attrs["id"] += "richtext-" + fieldName
	} else {
		attrs["id"] = "richtext-" + fieldName
	}

	// create the target element for the editor to attach itself
	div := &Element{
		TagName: "div",
		Attrs:   attrs,
		Name:    "",
		Label:   "",
		Data:    "",
		ViewBuf: &bytes.Buffer{},
	}

	// create a hidden input to store the value from the struct
	val := ValueFromStructField(fieldName, p, args).(string)
	name := TagNameFromStructField(fieldName, p, args)
	input := `<input type="hidden" name="` + name + `" class="richtext-value ` + fieldName + `" value="` + html.EscapeString(val) + `"/>`

	// build the dom tree for the entire richtext component
	iso = append(iso, DOMElement(div)...)
	iso = append(iso, []byte(input)...)
	iso = append(iso, isoClose...)

	script := `
	<script>
		$(function() { 
			var _editor = $('.richtext.` + fieldName + `');
			var hidden = $('.richtext-value.` + fieldName + `');

			_editor.materialnote({
				height: 250,
				placeholder: '` + attrs["placeholder"] + `',
				toolbar: [
					['style', ['style']],
					['font', ['bold', 'italic', 'underline', 'clear', 'strikethrough', 'superscript', 'subscript']],
					['fontsize', ['fontsize']],
					['color', ['color']],
					['insert', ['link', 'picture', 'video', 'hr']],					
					['para', ['ul', 'ol', 'paragraph']],
					['table', ['table']],
					['height', ['height']],
					['misc', ['codeview']]
				],
				// intercept file insertion, storage and insert img with new src
				onImageUpload: function(uploads) {
					var data = new FormData();
					data.append("file", uploads[0]);
					$.ajax({
						data: data,
						type: 'PUT',	
						url: '/edit/upload',
						cache: false,
						contentType: false,
						processData: false,
						success: function(resp) {
							var img = document.createElement('img');
							img.setAttribute('src', resp.data[0].url);
							_editor.materialnote('insertNode', img);
						},
						error: function(xhr, status, err) {
							console.log(status, err);
						}
					})

				}
			});

			// inject entities into editor
			if (hidden.val() !== "") {
				_editor.code(hidden.val());
			}

			// update hidden input with encoded value on different events
			_editor.on('materialnote.change', function(e, entities, $editable) {
				hidden.val(replaceBadChars(entities));			
			});

			_editor.on('materialnote.paste', function(e) {
				hidden.val(replaceBadChars(_editor.code()));			
			});

			// bit of a hack to stop the editor buttons from causing a refresh when clicked 
			$('.note-toolbar').find('button, i, a').on('click', function(e) { e.preventDefault(); });
		});
	</script>`

	return append(iso, []byte(script)...)
}
