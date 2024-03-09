package editor

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
