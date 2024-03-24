package editor

import "strings"

// Tags returns the []byte of a tag input (in the style of Materialze 'Chips') with a label.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func Tags(fieldName string, p interface{}, attrs map[string]string) []byte {
	name := TagNameFromStructField(fieldName, p, nil)

	// get the saved tags if this is already an existing post
	values := ValueFromStructField(fieldName, p, nil).(string)
	var tags []string
	if strings.Contains(values, "__ponzu") {
		tags = strings.Split(values, "__ponzu")
	}

	// case where there is only one tag stored, thus has no separator
	if len(values) > 0 && !strings.Contains(values, "__ponzu") {
		tags = append(tags, values)
	}

	tmpl := `
	<div class="col s12 __ponzu-tags ` + name + `">
		<label class="active">` + attrs["label"] + ` (Type and press "Enter")</label>
		<div class="chips ` + name + `"></div>
	`

	var initial []string
	i := 0
	for _, tag := range tags {
		tagName := TagNameFromStructFieldMulti(fieldName, i, p)
		tmpl += `<input type="hidden" class="__ponzu-tag ` + tag + `" name=` + tagName + ` value="` + tag + `"/>`
		initial = append(initial, `{tag: '`+tag+`'}`)
		i++
	}

	script := `
	<script>
		$(function() {
			var tags = $('.__ponzu-tags.` + name + `');
			$('.chips.` + name + `').material_chip({
				data: [` + strings.Join(initial, ",") + `],
				secondaryPlaceholder: '+` + name + `'
			});		

			// handle events specific to tags
			var chips = tags.find('.chips');
			
			chips.on('chip.add', function(e, chip) {
				chips.parent().find('.empty-tag').remove();
				
				var input = $('<input>');
				input.attr({
					class: '__ponzu-tag '+chip.tag.split(' ').join('__'),
					name: '` + name + `.'+String(tags.find('input[type=hidden]').length),
					value: chip.tag,
					type: 'hidden'
				});
				
				tags.append(input);
			});

			chips.on('chip.delete', function(e, chip) {
				// convert tag string to class-like selector "some tag" -> ".some.tag"
				var sel = '.__ponzu-tag.' + chip.tag.split(' ').join('__');
				chips.parent().find(sel).remove();

				// iterate through all hidden tag inputs to re-name them with the correct ` + name + `.index
				var hidden = chips.parent().find('input[type=hidden]');
				
				// if there are no tags, set a blank
				if (hidden.length === 0) {
					var input = $('<input>');
					input.attr({
						class: 'empty-tag',
						name: '` + name + `',
						type: 'hidden'
					});
					
					tags.append(input);
				}
				
				// re-name hidden storage elements in necessary format 
				for (var i = 0; i < hidden.length; i++) {
					$(hidden[i]).attr('name', '` + name + `.'+String(i));
				}
			});
		});
	</script>
	`

	tmpl += `</div>`

	return []byte(tmpl + script)
}
