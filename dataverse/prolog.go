package dataverse

import (
	"bytes"
	"html/template"
)

var govCheckTplStr = `:- consult('{{ .GovCode }}').

action('{{ .Action }}').
subject('{{ .Subject }}').
zone('{{ .Zone }}').

tell(Result, Evidence) :-
bagof(P:Modality, paragraph(P, Modality), Evidence),
(   member(_: 'prohibited', Evidence) -> Result = 'prohibited'
;   member(_: 'permitted', Evidence) -> Result = 'permitted'
).`

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
