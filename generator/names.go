package generator

import "strings"

// get the initial field name passed and check it for all possible cases
// MyTitle:string myTitle:string my_title:string -> MyTitle
// error-message:string -> ErrorMessage
func parseName(rawName string) (string, string) {
	parts := strings.Split(rawName, ":")
	name := cleanName(parts[0])
	label := name
	if len(parts) > 1 {
		label = cleanName(parts[1])
	}

	return name, label
}

func cleanName(name string) string {
	// remove _ or - if first character
	if name[0] == '-' || name[0] == '_' {
		name = name[1:]
	}

	// remove _ or - if last character
	if name[len(name)-1] == '-' || name[len(name)-1] == '_' {
		name = name[:len(name)-1]
	}

	// upcase the first character
	name = strings.ToUpper(string(name[0])) + name[1:]

	// remove _ or - character, and upcase the character immediately following
	for i := 0; i < len(name); i++ {
		r := rune(name[i])
		if isUnderscore(r) || isHyphen(r) {
			up := strings.ToUpper(string(name[i+1]))
			name = name[:i] + up + name[i+2:]
		}
	}

	return name
}

// get the field name passed and convert to json-like name
// MyTitle:string myTitle:string my_title:string -> my_title
// error-message:string -> error-message
func getJSONName(name string) string {
	// remove _ or - if first character
	if name[0] == '-' || name[0] == '_' {
		name = name[1:]
	}

	// downcase the first character
	name = strings.ToLower(string(name[0])) + name[1:]

	// check for uppercase character, downcase and insert _ before it if i-1
	// isn'types already _ or -
	for i := 0; i < len(name); i++ {
		r := rune(name[i])
		if isUpper(r) {
			low := strings.ToLower(string(r))
			if name[i-1] == '_' || name[i-1] == '-' {
				name = name[:i] + low + name[i+1:]
			} else {
				name = name[:i] + "_" + low + name[i+1:]
			}
		}
	}

	return name
}

func isUnderscore(char rune) bool {
	return char == '_'
}

func isHyphen(char rune) bool {
	return char == '-'
}

func isUpper(char rune) bool {
	if char >= 'A' && char <= 'Z' {
		return true
	}

	return false
}
