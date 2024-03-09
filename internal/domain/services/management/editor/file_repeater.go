package editor

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

// FileRepeater returns the []byte of a <input type="file"> HTML element with a label.
// It also includes repeat controllers (+ / -) so the element can be
// dynamically multiplied or reduced.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func FileRepeater(fieldName string, p interface{}, attrs map[string]string) []byte {
	// find the field values in p to determine if an option is pre-selected
	fieldVals := ValueFromStructField(fieldName, p, nil).(string)
	vals := strings.Split(fieldVals, "__ponzu")

	addLabelFirst := func(i int, label string) string {
		if i == 0 {
			return `<label class="active">` + label + `</label>`
		}

		return ""
	}

	tmpl :=
		`<div class="file-input %[5]s %[4]s input-field col s12">
			%[2]s
			<div class="file-field input-field">
				<div class="btn">
					<span>Upload</span>
					<input class="storage %[4]s" type="file" />
				</div>
				<div class="file-path-wrapper">
					<input class="file-path validate" placeholder="Add %[5]s" type="text" />
				</div>
			</div>
			<div class="preview"><div class="img-clip"></div></div>			
			<input class="store %[4]s" type="hidden" name="%[1]s" value="%[3]s" />
		</div>`
	// 1=nameidx, 2=addLabelFirst, 3=val, 4=className, 5=fieldName
	script :=
		`<script>
			$(function() {
				var $file = $('.file-input.%[2]s'),
					storage = $file.find('input.storage'),
					store = $file.find('input.store'),
					preview = $file.find('.preview'),
					clip = preview.find('.img-clip'),
					reset = document.createElement('div'),
					img = document.createElement('img'),
					video = document.createElement('video'),
					unknown = document.createElement('div'),
					viewLink = document.createElement('a'),
					viewLinkText = document.createTextNode('Download / View '),
					iconLaunch = document.createElement('i'),
					iconLaunchText = document.createTextNode('launch'),
					uploadSrc = store.val();
					video.setAttribute
					preview.hide();
					viewLink.setAttribute('href', '%[3]s');
					viewLink.setAttribute('target', '_blank');
					viewLink.appendChild(viewLinkText);
					viewLink.style.display = 'block';
					viewLink.style.marginRight = '10px';					
					viewLink.style.textAlign = 'right';
					iconLaunch.className = 'material-icons tiny';
					iconLaunch.style.position = 'relative';
					iconLaunch.style.top = '3px';
					iconLaunch.appendChild(iconLaunchText);
					viewLink.appendChild(iconLaunch);
					preview.append(viewLink);
				
				// when %[2]s input changes (file is selected), remove
				// the 'name' and 'value' attrs from the hidden store input.
				// add the 'name' attr to %[2]s input
				storage.on('change', function(e) {
					resetImage();
				});

				if (uploadSrc.length > 0) {
					var ext = uploadSrc.substring(uploadSrc.lastIndexOf('.'));
					ext = ext.toLowerCase();
					switch (ext) {
						case '.jpg':
						case '.jpeg':
						case '.webp':
						case '.gif':
						case '.png':
							$(img).attr('src', store.val());
							clip.append(img);
							break;
						case '.mp4':
						case '.webm':
							$(video)
								.attr('src', store.val())
								.attr('type', 'video/'+ext.substring(1))
								.attr('controls', true)
								.css('width', '100%%');
							clip.append(video);
							break;
						default:
							$(img).attr('src', '/static/dashboard/img/ponzu-file.png');
							$(unknown)
								.css({
									position: 'absolute', 
									top: '10px', 
									left: '10px',
									border: 'solid 1px #ddd',
									padding: '7px 7px 5px 12px',
									fontWeight: 'bold',
									background: '#888',
									color: '#fff',
									textTransform: 'uppercase',
									letterSpacing: '2px' 
								})
								.text(ext);
							clip.append(img);
							clip.append(unknown);
							clip.css('maxWidth', '200px');
					}
					preview.show();

					$(reset).addClass('reset %[2]s btn waves-effect waves-light grey');
					$(reset).html('<i class="material-icons tiny">clear<i>');
					$(reset).on('click', function(e) {
						e.preventDefault();
						preview.animate({"opacity": 0.1}, 200, function() {
							preview.slideUp(250, function() {
								resetImage();
							});
						})
						
					});
					clip.append(reset);
				}

				function resetImage() {
					store.val('');
					store.attr('name', '');
					storage.attr('name', '%[1]s');
					clip.empty();
				}
			});	
		</script>`
	// 1=nameidx, 2=className

	name := TagNameFromStructField(fieldName, p, nil)

	html := bytes.Buffer{}
	_, err := html.WriteString(`<span class="__ponzu-repeat ` + name + `">`)
	if err != nil {
		log.Println("Error writing HTML string to FileRepeater buffer")
		return nil
	}

	for i, val := range vals {
		className := fmt.Sprintf("%s-%d", name, i)
		nameidx := TagNameFromStructFieldMulti(fieldName, i, p)

		_, err := html.WriteString(fmt.Sprintf(tmpl, nameidx, addLabelFirst(i, attrs["label"]), val, className, fieldName))
		if err != nil {
			log.Println("Error writing HTML string to FileRepeater buffer")
			return nil
		}

		_, err = html.WriteString(fmt.Sprintf(script, nameidx, className, val))
		if err != nil {
			log.Println("Error writing HTML string to FileRepeater buffer")
			return nil
		}
	}
	_, err = html.WriteString(`</span>`)
	if err != nil {
		log.Println("Error writing HTML string to FileRepeater buffer")
		return nil
	}

	return append(html.Bytes(), RepeatController(fieldName, p, "input.storage", "div.file-input."+fieldName)...)
}
