package main

import (
	"bytes"
	"encoding/json"
	"go/format"
	"io/ioutil"
	"log"
	"strings"
	"text/template"

	"github.com/captncraig/wildchef/constants"
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
	consts := map[string]string{}
	for k, v := range nameLookup {
		consts[clean(v)] = v
		nameLookup[k] = clean(v)
	}
	ctx := map[string]interface{}{
		"Consts":     consts,
		"NameLookup": nameLookup,
		"Sprites":    constants.SpriteLocations,
	}
	buf := &bytes.Buffer{}
	must(gotpl.Execute(buf, ctx))
	dat, err = format.Source(buf.Bytes())
	must(err)
	ioutil.WriteFile("constants/constants.go", dat, 0777)

	buf.Reset()
	must(tsTpl.Execute(buf, ctx))
	ioutil.WriteFile("gen.ts", buf.Bytes(), 0777)
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

var tsTpl = template.Must(template.New("").Parse(`
export var Names = {
	{{range $k,$v := .Consts}} {{$k}}: "{{$v}}",
	{{end -}}
}

export interface SpriteLoc{
	X: number;
	Y: number;
}

export var Sprites: {[id: string]: SpriteLoc} = {
	{{range $k,$v := .Sprites}} "{{$k}}": {X: {{$v.X}}, Y: {{$v.Y}}},
	{{end -}}
}
`))
