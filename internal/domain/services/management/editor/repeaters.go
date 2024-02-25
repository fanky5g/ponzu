package editor

import (
	"bytes"
	"fmt"
	"github.com/fanky5g/ponzu/internal/util"
	"log"
	"reflect"
	"strings"
)

func InputRepeater(fieldName string, p interface{}, attrs map[string]string) []byte {
	// find the field values in p to determine pre-filled inputs
	fieldVals := ValueFromStructField(fieldName, p, nil).(string)
	vals := strings.Split(fieldVals, "__ponzu")

	scope := TagNameFromStructField(fieldName, p, nil)
	html := bytes.Buffer{}

	_, err := html.WriteString(`<span class="__ponzu-repeat ` + scope + `">`)
	if err != nil {
		log.Println("Error writing HTML string to InputRepeater buffer")
		return nil
	}

	for i, val := range vals {
		el := &Element{
			TagName: "input",
			Attrs:   attrs,
			Name:    TagNameFromStructFieldMulti(fieldName, i, p),
			Data:    val,
			ViewBuf: &bytes.Buffer{},
		}

		// only add the label to the first input in repeated list
		if i == 0 {
			el.Label = attrs["label"]
		}

		_, err := html.Write(DOMElementSelfClose(el))
		if err != nil {
			log.Println("Error writing DOMElementSelfClose to InputRepeater buffer")
			return nil
		}
	}
	_, err = html.WriteString(`</span>`)
	if err != nil {
		log.Println("Error writing HTML string to InputRepeater buffer")
		return nil
	}

	return append(html.Bytes(), RepeatController(fieldName, p, "input", ".input-field")...)
}

// SelectRepeater returns the []byte of a <select> HTML element plus internal <options> with a label.
// It also includes repeat controllers (+ / -) so the element can be
// dynamically multiplied or reduced.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func SelectRepeater(fieldName string, p interface{}, attrs, options map[string]string) []byte {
	// options are the value attr and the display value, i.e.
	// <option value="{map key}">{map value}</option>
	scope := TagNameFromStructField(fieldName, p, nil)
	html := bytes.Buffer{}
	_, err := html.WriteString(`<span class="__ponzu-repeat ` + scope + `">`)
	if err != nil {
		log.Println("Error writing HTML string to SelectRepeater buffer")
		return nil
	}

	// find the field values in p to determine if an option is pre-selected
	fieldVals := ValueFromStructField(fieldName, p, nil).(string)
	vals := strings.Split(fieldVals, "__ponzu")

	if _, ok := attrs["class"]; ok {
		attrs["class"] += " browser-default"
	} else {
		attrs["class"] = "browser-default"
	}

	// loop through vals and create selects and options for each, adding to views
	if len(vals) > 0 {
		for i, val := range vals {
			sel := &Element{
				TagName: "select",
				Attrs:   attrs,
				Name:    TagNameFromStructFieldMulti(fieldName, i, p),
				ViewBuf: &bytes.Buffer{},
			}

			// only add the label to the first select in repeated list
			if i == 0 {
				sel.Label = attrs["label"]
			}

			// create options for select element
			var opts []*Element

			// provide a call to action for the select element
			cta := &Element{
				TagName: "option",
				Attrs:   map[string]string{"disabled": "true", "selected": "true"},
				Data:    "Select an option...",
				ViewBuf: &bytes.Buffer{},
			}

			// provide a selection reset (will store empty string in db)
			reset := &Element{
				TagName: "option",
				Attrs:   map[string]string{"value": ""},
				Data:    "None",
				ViewBuf: &bytes.Buffer{},
			}

			opts = append(opts, cta, reset)

			for k, v := range options {
				optAttrs := map[string]string{"value": k}
				if k == val {
					optAttrs["selected"] = "true"
				}
				opt := &Element{
					TagName: "option",
					Attrs:   optAttrs,
					Data:    v,
					ViewBuf: &bytes.Buffer{},
				}

				opts = append(opts, opt)
			}

			_, err := html.Write(DOMElementWithChildrenSelect(sel, opts))
			if err != nil {
				log.Println("Error writing DOMElementWithChildrenSelect to SelectRepeater buffer")
				return nil
			}
		}
	}

	_, err = html.WriteString(`</span>`)
	if err != nil {
		log.Println("Error writing HTML string to SelectRepeater buffer")
		return nil
	}

	return append(html.Bytes(), RepeatController(fieldName, p, "select", ".input-field")...)
}

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

// RepeatController generates the javascript to control any repeatable form
// element in an editor based on its type, field name and HTML tag name
func RepeatController(fieldName string, p interface{}, inputSelector, cloneSelector string) []byte {
	scope := TagNameFromStructField(fieldName, p, nil)
	script := `
    <script>
        $(function() {
            // define the scope of the repeater
            var scope = $('.__ponzu-repeat.` + scope + `');

            var getChildren = function() {
                return scope.find('` + cloneSelector + `')
            }

            var resetFieldNames = function() {
                // loop through children, set its name to the fieldName.i where
                // i is the current index number of children array
                var children = getChildren();

                for (var i = 0; i < children.length; i++) {
					var preset = false;					
                    var $el = children.eq(i);
					var name = '` + scope + `.'+String(i);

                    $el.find('` + inputSelector + `').attr('name', name);

					// ensure no other input-like elements besides ` + inputSelector + `
					// get the new name by setting it to an empty string
					$el.find('input, select, textarea').each(function(i, elem) {
						var $elem = $(elem);
						
						// if the elem is not ` + inputSelector + ` and has no value 
						// set the name to an empty string
						if (!$elem.is('` + inputSelector + `')) {
							if ($elem.val() === '' || $elem.is('.file-path')) {
								$elem.attr('name', '');
							} else {
								$elem.attr('name', name);
								preset = true;
							}						
						}
					});      

					// if there is a preset value, remove the name attr from the
					// ` + inputSelector + ` element so it doesn't overwrite db
					if (preset) {
						$el.find('` + inputSelector + `').attr('name', '');														
					}          

                    // reset controllers
                    $el.find('.controls').remove();
                }

                applyRepeatControllers();
            }

            var addRepeater = function(e) {
                e.preventDefault();
                
                var add = e.target;

                // find and clone the repeatable input-like element
                var source = $(add).parent().closest('` + cloneSelector + `');
                var clone = source.clone();

                // if clone has label, remove it
                clone.find('label').remove();
                
                // remove the pre-filled value from clone
                clone.find('` + inputSelector + `').val('');
				clone.find('input').val('');

                // remove controls from clone if already present
                clone.find('.controls').remove();

				// remove input preview on clone if copied from source
				clone.find('.preview').remove();

                // add clone to scope and reset field name attributes
                scope.append(clone);

                resetFieldNames();
            }

            var delRepeater = function(e) {
                e.preventDefault();

                // do nothing if the input is the only child
                var children = getChildren();
                if (children.length === 1) {
                    return;
                }

                var del = e.target;
                
                // pass label onto next input-like element if del 0 index
                var wrapper = $(del).parent().closest('` + cloneSelector + `');
                if (wrapper.find('` + inputSelector + `').attr('name') === '` + scope + `.0') {
                    wrapper.next().append(wrapper.find('label'))
                }
                
                wrapper.remove();

                resetFieldNames();
            }

            var createControls = function() {
                // create + / - controls for each input-like child element of scope
                var add = $('<button>+</button>');
                add.addClass('repeater-add');
                add.addClass('btn-flat waves-effect waves-green');

                var del = $('<button>-</button>');
                del.addClass('repeater-del');
                del.addClass('btn-flat waves-effect waves-red');                

                var controls = $('<span></span>');
                controls.addClass('controls');
                controls.addClass('right');

                // bind listeners to child's controls
                add.on('click', addRepeater);
                del.on('click', delRepeater);

                controls.append(add);
                controls.append(del);

                return controls;
            }

            var applyRepeatControllers = function() {
                // add controls to each child
                var children = getChildren()
                for (var i = 0; i < children.length; i++) {
                    var el = children[i];
                    
                    $(el).find('` + inputSelector + `').parent().find('.controls').remove();
                    
                    var controls = createControls();                                        
                    $(el).append(controls);
                }
            }

			resetFieldNames();
        });

    </script>
    `

	return []byte(script)
}

func NestedRepeater(fieldName string, p interface{}, m func(v interface{}, f *FieldArgs) (string, []Field)) []byte {
	value := ValueByName(fieldName, p, nil)
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		panic(fmt.Sprintf("Ponzu: Type '%s' for field '%s' not supported.", value.Type(), fieldName))
	}

	scope := TagNameFromStructField(fieldName, p, nil)

	tmpl := `
		<div class="col s12 __ponzu-repeat ` + scope + `" name="` + fieldName + `">
			<label class="active">` + fieldName + `</label>
	`

	positionalPlaceHolder := "%pos%"
	fieldArgs := &FieldArgs{
		Parent: fmt.Sprintf("%s.%s", fieldName, positionalPlaceHolder),
	}

	arrayTypeName, fields := m(p, fieldArgs)
	fieldArgs.TypeName = arrayTypeName
	emptyEntryTemplate := Nested("", p, fieldArgs, fields...)

	script := &bytes.Buffer{}
	scriptTmpl := util.MakeScript("nested_repeater")
	if err := scriptTmpl.Execute(script, struct {
		Template              string
		NumItems              int
		Scope                 string
		InputSelector         string
		CloneSelector         string
		PositionalPlaceholder string
	}{
		Template:              string(emptyEntryTemplate),
		NumItems:              value.Len(),
		Scope:                 scope,
		CloneSelector:         fmt.Sprintf(".%s", arrayTypeName),
		InputSelector:         "input",
		PositionalPlaceholder: positionalPlaceHolder,
	}); err != nil {
		panic(err)
	}

	for i := 0; i < value.Len(); i++ {
		_, fields = m(p, &FieldArgs{
			Parent: fmt.Sprintf("%s.%d", fieldName, i),
		})

		fieldTemplate := Nested("", p, fieldArgs, fields...)
		tmpl += string(fieldTemplate)
	}

	tmpl += `</div>`
	return append([]byte(tmpl), script.Bytes()...)
}
