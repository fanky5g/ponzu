package views

import (
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
	"runtime"
)

var (
	t    *template.Template
	Path string
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	Path = filepath.Join(filepath.Dir(b), "../../internal/views")
	t = template.New("views").Funcs(GlobFuncs)

	traverseErr := filepath.WalkDir(Path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			pattern := filepath.Join(path, "*.gohtml")
			matches, err := filepath.Glob(pattern)
			if err != nil {
				return err
			}

			if len(matches) > 0 {
				t, err = t.ParseGlob(pattern)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	if traverseErr != nil {
		panic(traverseErr)
	}
}

func ExecuteTemplate(w io.Writer, name string, data interface{}) error {
	return t.ExecuteTemplate(w, name, data)
}
