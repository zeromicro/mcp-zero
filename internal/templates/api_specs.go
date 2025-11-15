package templates

const APISpecTemplate = `syntax = "v1"

info (
	title: "{{.ServiceName}} API"
	desc: "API specification for {{.ServiceName}}"
	version: "1.0"
)

{{range .Types}}
type {{.Name}} {
{{range .Fields}}	{{.Name}} {{.Type}} ` + "`json:\"{{.JsonTag}}\"`" + `{{if .Comment}} // {{.Comment}}{{end}}
{{end}}}

{{end}}
service {{.ServiceName}}-api {
{{range .Endpoints}}	@handler {{.Handler}}
	{{.Method}} {{.Path}}{{if .Request}} ({{.Request}}){{end}}{{if .Response}} returns ({{.Response}}){{end}}
{{end}}}
`

type APISpec struct {
	ServiceName string
	Types       []TypeDef
	Endpoints   []EndpointDef
}

type TypeDef struct {
	Name   string
	Fields []FieldDef
}

type FieldDef struct {
	Name    string
	Type    string
	JsonTag string
	Comment string
}

type EndpointDef struct {
	Handler  string
	Method   string
	Path     string
	Request  string
	Response string
}
