package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"os"
	"regexp"
	"strings"
)

//go:embed html/template.html
var htmlTemplate string

//go:embed html/dttemplate.html
var htmlDTTemplate string

var currentID uint

var (
	escapedHTML = regexp.MustCompile(`\\<([^<>]*)>`)
	dt          = regexp.MustCompile(`<<(.*?)>>`)
	bold        = regexp.MustCompile(`\*(.*?)\*`)
)

// TODO:implement html escaping
func process(s string) template.HTML {
	s = escapedHTML.ReplaceAllString(s, `&lt;$1&gt;`)
	s = bold.ReplaceAllString(s, `<b>$1</b>`)
	s = dt.ReplaceAllStringFunc(s, func(repl string) string {
		repl = strings.ReplaceAll(repl, "<<", "")
		repl = strings.ReplaceAll(repl, ">>", "")
		tpl, err := template.New("dt" + string(currentID)).Parse(htmlDTTemplate)
		if err != nil {
			log.Fatalf("failed to parse dt template: %s", err)
		}
		var buf bytes.Buffer
		if err := tpl.Execute(&buf, struct {
			Id      uint
			TimeFmt template.JSStr
		}{Id: currentID, TimeFmt: template.JSStr(repl)}); err != nil {
			log.Fatalf("failed to execute dt template: %s", err)
		}
		currentID += 1
		return buf.String()
	})
	return template.HTML(s)
}

func (ss slideshow) toHTML() ([]byte, error) {
	funcmap := template.FuncMap{
		"process": process,
		"getScript": func() template.HTML {
			if ss.Script.Has {
				script, err := os.ReadFile(ss.Script.Val)
				if err != nil {
					log.Fatalf("couldn't read script at url %q: %s", ss.Script.Val, err)
				}
				return template.HTML(fmt.Sprintf("<script>%s</script>", script))
			}
			return template.HTML("")
		},
		"getStyles": func() template.HTML {
			if ss.Styles.Has {
				styles, err := os.ReadFile(ss.Styles.Val)
				if err != nil {
					log.Fatalf("couldn't read stylesheet at url %q: %s", ss.Styles.Val, err)
				}
				return template.HTML(fmt.Sprintf("<style>%s</style>", string(styles)))
			}
			return template.HTML("")
		},
	}
	tmpl, err := template.New("page").Funcs(funcmap).Parse(htmlTemplate)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ss); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
