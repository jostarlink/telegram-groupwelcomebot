package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"strings"
)

type settingsPage struct {
	Title  string
	Status bool
	Config map[string]template.HTML
}

type rulesPage struct {
	Title string
	Rules template.HTML
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := template.New("index")
	t, _ = template.ParseFiles("html/index.html")
	t.Execute(w, nil)
}

func rules(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	s := storage.GetByString(param.ByName("channel"))

	e := template.HTMLEscapeString(s["rules"])
	b := strings.Replace(e, "\n", "<br>\n", -1)

	c := &rulesPage{
		Title: s["name"],
		Rules: template.HTML(b),
	}

	t := template.New("rules")
	t, _ = template.ParseFiles("html/rules.html")
	t.Execute(w, c)
}

func settings(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	s := storage.GetChannelConfig(param.ByName("channel"))

	e := make(map[string]template.HTML)

	for k, v := range s.Config {
		e[k] = template.HTML(strings.Replace(template.HTMLEscapeString(v), "\n", "<br>\n", -1))
	}

	c := &settingsPage{
		Title:  s.Config["name"],
		Status: s.Enabled,
		Config: e,
	}

	t := template.New("channel")
	t, _ = template.ParseFiles("html/settings.html")
	t.Execute(w, c)
}
