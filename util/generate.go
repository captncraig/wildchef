package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"strings"
	"text/template"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	dat, err := ioutil.ReadFile("items.json")
	must(err)
	nameLookup := map[string]string{}
	must(json.Unmarshal(dat, &nameLookup))
	constants := map[string]string{}
	for k, v := range nameLookup {
		constants[clean(v)] = v
		nameLookup[k] = clean(v)
	}
	ctx := map[string]interface{}{
		"Consts":     constants,
		"NameLookup": nameLookup,
	}
	buf := &bytes.Buffer{}
	must(gotpl.Execute(buf, ctx))
	fmt.Println(buf.String())
	dat, err = format.Source(buf.Bytes())
	must(err)
	ioutil.WriteFile("constants/constants.go", dat, 0777)
}

func clean(s string) string {
	for _, bad := range []string{" ", "-", "'"} {
		s = strings.Replace(s, bad, "", -1)
	}
	return s
}

var gotpl = template.Must(template.New("").Parse(`package constants

const(
	{{range $k,$v := .Consts}} {{$k}} = "{{$v}}"
	{{end -}}
)
var ItemIds = map[string]string{
	{{range $k,$v := .NameLookup}}"{{$k}}": {{$v}},
	{{end -}}
}
`))
