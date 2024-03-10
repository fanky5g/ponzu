package views

import (
	"bytes"
	conf "github.com/fanky5g/ponzu/config"
	"github.com/fanky5g/ponzu/internal/util"
)

// ForgotPassword ...
func ForgotPassword(appName string, paths conf.Paths) ([]byte, error) {
	a := View{
		Logo:       appName,
		PublicPath: paths.PublicPath,
	}

	buf := &bytes.Buffer{}
	tmpl := util.MakeTemplate("start_admin", "forgot_password", "end_admin")
	err := tmpl.Execute(buf, a)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
