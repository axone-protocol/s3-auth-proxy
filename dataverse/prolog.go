package dataverse

import (
	"bytes"
	_ "embed"
	"text/template"
)

//go:embed gov-check-program-tpl.pl
var govCheckTplStr string

var govCheckTpl *template.Template

func init() {
	tpl, err := template.New("govCheckProgram").Parse(govCheckTplStr)
	if err != nil {
		panic(err)
	}

	govCheckTpl = tpl
}

func makeGovCheckProgram(govCode, action, subject, zone string) (string, error) {
	buf := bytes.Buffer{}
	if err := govCheckTpl.Execute(&buf, map[string]string{
		"GovCode": govCode,
		"Action":  action,
		"Subject": subject,
		"Zone":    zone,
	}); err != nil {
		return "", err
	}

	return buf.String(), nil
}
