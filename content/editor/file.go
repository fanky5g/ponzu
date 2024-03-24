package editor

import "fmt"

// File returns the []byte of a <input type="file"> HTML element with a label.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func File(fieldName string, p interface{}, attrs map[string]string) []byte {
	name := TagNameFromStructField(fieldName, p, nil)
	value := ValueFromStructField(fieldName, p, nil).(string)
	tmpl :=
		`<div class="file-input ` + name + ` input-field col s12">
			<label class="active">` + attrs["label"] + `</label>
			<div class="file-field input-field">
				<div class="btn">
					<span>Upload</span>
					<input class="storage" type="file">
				</div>
				<div class="file-path-wrapper">
					<input class="file-path validate" placeholder="` + attrs["label"] + `" type="text">
				</div>
			</div>
			<div class="preview"><div class="img-clip"></div></div>			
			<input class="store ` + name + `" type="hidden" name="` + name + `" value="` + value + `" />
		</div>`

	script :=
		fmt.Sprintf(`<script>
			$(function() {
				var $file = $('.file-input.` + name + `'),
					storage = $file.find('input.storage'),
					store = $file.find('input.store'),
					preview = $file.find('.preview'),
					clip = preview.find('.img-clip'),
					reset = document.createElement('div'),
					img = document.createElement('img'),
					video = document.createElement('video'),
					unknown = document.createElement('div'),
					viewLink = document.createElement('a'),
					viewLinkText = document.createTextNode('Download / Template '),
					iconLaunch = document.createElement('i'),
					iconLaunchText = document.createTextNode('launch'),
					uploadSrc = store.val();
					video.setAttribute
					preview.hide();
					viewLink.setAttribute('href', '` + value + `');
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

				// when ` + name + ` input changes (file is selected), remove
				// the 'name' and 'value' attrs from the hidden store input.
				// add the 'name' attr to ` + name + ` input
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
								.css('width', '100%');
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

					$(reset).addClass('reset ` + name + ` btn waves-effect waves-light grey');
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
					storage.attr('name', '` + name + `');
					clip.empty();
				}
			});	
		</script>`)

	return []byte(tmpl + script)
}
