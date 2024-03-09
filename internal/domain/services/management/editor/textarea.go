package editor

// Textarea returns the []byte of a <textarea> HTML element with a label.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func Textarea(fieldName string, p interface{}, attrs map[string]string) []byte {
	// add materialize css class to make UI correct
	className := "materialize-textarea"
	if _, ok := attrs["class"]; ok {
		class := attrs["class"]
		attrs["class"] = class + " " + className
	} else {
		attrs["class"] = className
	}

	e := NewElement("textarea", attrs["label"], fieldName, p, attrs, nil)

	return DOMElement(e)
}
