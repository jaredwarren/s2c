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

// UnmarshalJSON ...
// func (s *Swagger) UnmarshalJSON(data []byte) error {
// 	var b map[string]interface{}
// 	if err := json.Unmarshal(data, &b); err != nil {
// 		return err
// 	}

// 	// fmt.Printf("%+v\n", b)
// 	s = &Swagger{
// 		Version: b["swagger"].(string),
// 		Info: Info{
// 			Title:       b["info"].(map[string]interface{})["title"].(string),
// 			Description: b["info"].(map[string]interface{})["description"].(string),
// 			Version:     b["info"].(map[string]interface{})["version"].(string),
// 		},
// 		Host:  b["host"].(string),
// 		Paths: Paths{},
// 		// Definitions: b["definitions"].(map[string]Definition),
// 	}

// 	for k, p := range b["paths"].(map[string]interface{}) {
// 		m := map[string]*Method{}
// 		s.Paths[k] = Path{
// 			Path:    k,
// 			Host:    s.Host,
// 			Methods: &m,
// 		}
// 	}

// 	// *p = make(map[string]Path)
// 	// for k, v := range b {
// 	// 	m := Path{
// 	// 		Path:    k,
// 	// 		Methods: v.Methods,
// 	// 	}
// 	// 	for _, mm := range *m.Methods {
// 	// 		mm.Path = k
// 	// 	}
// 	// 	(*p)[k] = m
// 	// }
// 	return nil
// }

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

	// ms := map[string]Method{}
	p.Methods = &map[string]*Method{}
	for k, v := range b {
		// fmt.Println("  - ", k)
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
		// ms[k] = m
	}
	// p.Methods = &ms
	// fmt.Printf("%+v\n", p.Methods)
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
// TODO: add data and/or args
func (m Method) ToCurl(host string) string {
	schema := "http"
	// default to first schema
	if len(m.Schemes) > 0 {
		schema = m.Schemes[0]
	}
	url := fmt.Sprintf("%s://%s%s", schema, host, m.Path)

	// check if there are params
	data := ""

	return fmt.Sprintf(`
curl -L --header "Content-Type: application/json" \
-- request %s \
%s%s
`,
		m.Operation, data, url)
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

// SwagToCurl ...
func SwagToCurl(sw Swagger, path, method string) (err error) {
	// 	for p, pInfo := range sw.Paths {
	// 		if path == "" || p == path {
	// 			for m, mInfo := range pInfo {
	// 				if method == "" || m == method {
	// 					fmt.Printf(" -> %s -> %s -> %+v\n", p, m, mInfo)
	// 					// output curl

	// 					// url
	// 					// TODO: figure out how to set this. (schema for http?)
	// 					schema := "http"
	// 					domain := "localhost"
	// 					port := ":8080"
	// 					url := fmt.Sprintf("%s://%s%s%s", schema, domain, port, p)

	// 					// check if there are params
	// 					data := ""

	// 					fmt.Printf(`
	// curl -L --header "Content-Type: application/json" \
	// -- request %s \
	// %s%s
	// 					`, m, data, url)
	// 				}
	// 			}

	// 		}
	// 	}

	// TODO: run through Paths and generate curl
	return
}
