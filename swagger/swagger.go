package swagger

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Swagger ...
type Swagger struct {
	Version     string                `json:"swagger"`
	Info        Info                  `json:"info"`
	Host        string                `json:"host"`
	Paths       Paths                 `json:"paths"`
	Definitions map[string]Definition `json:"definitions"`
}

// ToCurl convert method to curl string
func (s Swagger) ToCurl() []string {
	curls := []string{}
	for _, path := range s.Paths {
		curls = append(curls, path.ToCurl(s.Host)...)
	}

	return curls
}

// FindPath convert method to curl string
func (s Swagger) FindPath(path string) *Path {
	for p, pInfo := range s.Paths {
		if p == path {
			return &pInfo
		}
	}
	return nil
}

// Paths ...
type Paths map[string]Path

// UnmarshalJSON ...
func (p *Paths) UnmarshalJSON(data []byte) error {
	var b map[string]Path
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	*p = make(map[string]Path)
	for k, v := range b {
		m := Path{
			Path:    k,
			Methods: v.Methods,
		}
		for _, mm := range *m.Methods {
			mm.Path = k
		}
		(*p)[k] = m
	}
	return nil
}

// Path ...
type Path struct {
	Path    string
	Methods *map[string]*Method
}

// UnmarshalJSON ...
func (p *Path) UnmarshalJSON(data []byte) error {
	var b map[string]Method
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}

	p.Methods = &map[string]*Method{}
	for k, v := range b {
		m := &Method{
			Tags:        v.Tags,
			Summary:     v.Summary,
			Description: v.Description,
			OperationID: v.OperationID,
			Parameters:  v.Parameters,
			Responses:   v.Responses,
			Schemes:     v.Schemes,
			Operation:   strings.ToUpper(k),
		}
		(*p.Methods)[k] = m
	}
	return nil
}

// ToCurl convert method to curl string
func (p Path) ToCurl(host string) []string {
	curls := []string{}
	for _, method := range *p.Methods {
		curls = append(curls, method.ToCurl(host))
	}

	return curls
}

// FindMethod convert method to curl string
func (p Path) FindMethod(method string) *Method {
	for m, mInfo := range *p.Methods {
		if m == method {
			return mInfo
		}
	}
	return nil
}

// Info ...
type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

// Method ...
type Method struct {
	Tags        []string            `json:"tags,omitempty"`
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	OperationID string              `json:"operationId"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses"`
	Schemes     []string            `json:"schemes"`
	Host        string
	Path        string
	Operation   string
}

// ToCurl convert method to curl string
func (m Method) ToCurl(host string) string {
	schema := "http"
	// default to first schema
	if len(m.Schemes) > 0 {
		schema = m.Schemes[0]
	}
	url := fmt.Sprintf("%s://%s%s", schema, host, m.Path)

	// check if there are params
	data := ""
	header := "\\"

	reqParms := map[string]string{}
	for _, pars := range m.Parameters {
		for _, par := range pars.Schema.Required {
			reqParms[par] = fmt.Sprintf("{{.%s}}", par)
		}
	}
	if len(reqParms) > 0 {
		// TODO: add option to use form data
		b, err := json.Marshal(reqParms)
		if err != nil {
			fmt.Println("error:", err)
		}
		data = fmt.Sprintf("--data '%s' \\\n", b)
		header = "--header \"Content-Type: application/json\" \\"
	}

	return fmt.Sprintf(`
curl -L %s
--request %s \
%s%s
`,
		header, m.Operation, data, url)
}

// GetRequiredParams convert method to curl string
func (m Method) GetRequiredParams() (params []string) {
	params = []string{}
	for _, pars := range m.Parameters {
		for _, par := range pars.Schema.Required {
			params = append(params, par)
		}
	}
	return
}

// Parameter ...
type Parameter struct {
	Name     string `json:"name"`
	In       string `json:"in"`
	Required bool   `json:"required"`
	Type     string `json:"type,omitempty"`
	Schema   Schema `json:"schema,omitempty"`
}

// Schema ...
type Schema struct {
	Ref      string   `json:"$ref,omitempty"`
	Required []string `json:"required,omitempty"`
}

// Response ...
type Response struct {
	Description string `json:"description,omitempty"`
	Schema      struct {
		Type string `json:"type"`
	} `json:"schema,omitempty"`
}

// Definition ...
type Definition struct {
	Title      string                 `json:"title"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Example    map[string]interface{} `json:"example"`
	Required   []string               `json:"required"`
}

// UnmarshalSwagger manually unmarshal swagger json.. not complete
// If I want to be able to get the path and host information in a method.
// I'll have unmarshal the json manually.
// func (s *Swagger) UnmarshalJSON(data []byte) error {
func UnmarshalSwagger(s *Swagger, data []byte) error {
	var b map[string]interface{}
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}

	// fmt.Printf("%+v\n", b)
	s = &Swagger{
		Version: b["swagger"].(string),
		Info: Info{
			Title:       b["info"].(map[string]interface{})["title"].(string),
			Description: b["info"].(map[string]interface{})["description"].(string),
			Version:     b["info"].(map[string]interface{})["version"].(string),
		},
		Host:  b["host"].(string),
		Paths: Paths{},
		// Definitions: b["definitions"].(map[string]Definition),
	}

	for k, p := range b["paths"].(map[string]interface{}) {
		fmt.Println(p)
		m := map[string]*Method{}
		s.Paths[k] = Path{
			Path: k,
			// Host:    s.Host,
			Methods: &m,
		}
	}

	// *p = make(map[string]Path)
	// for k, v := range b {
	// 	m := Path{
	// 		Path:    k,
	// 		Methods: v.Methods,
	// 	}
	// 	for _, mm := range *m.Methods {
	// 		mm.Path = k
	// 	}
	// 	(*p)[k] = m
	// }
	return nil
}
