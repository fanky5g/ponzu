package editor

// Textarea returns the []byte of a <textarea> HTML element with a label.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
func Textarea(fieldName string, p interface{}, attrs map[string]string) []byte {
    addClassName(attrs, "material-textarea")
	e := NewElement("textarea", attrs["label"], fieldName, p, attrs, nil)

	return DOMElement(e)
}
