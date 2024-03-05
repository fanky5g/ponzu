package editor

// Input returns the []byte of an <input> HTML element with a label.
// IMPORTANT:
// The `fieldName` argument will cause a panic if it is not exactly the string
// form of the struct field that this editor input is representing
//
//	type Person struct {
//		item.Item
//		editor editor.Editor
//
//		Name string `json:"name"`
//		//...
//	}
//
//	func (p *Person) MarshalEditor() ([]byte, error) {
//		view, err := editor.Form(p,
//			editor.Field{
//				View: editor.Input("Name", p, map[string]string{
//					"label":       "Name",
//					"type":        "text",
//					"placeholder": "Enter the Name here",
//				}),
//			}
//		)
//	}
func Input(fieldName string, p interface{}, attrs map[string]string, args *FieldArgs) []byte {
	e := NewElement("input", attrs["label"], fieldName, p, attrs, args)

	return DOMElementSelfClose(e)
}
