package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/espekkaya/jiralert-dockerize/jiralert/pkg/config"
)

const (
	docsUrl   = "https://github.com/espekkaya/jiralert-dockerize#readme"
	templates = `
    {{ define "page" -}}
      <html>
      <head>
        <title>JIRAlert</title>
        <style type="text/css">
          body { margin: 0; font-family: "Helvetica Neue", Helvetica, Arial, sans-serif; font-size: 14px; line-height: 1.42857143; color: #333; background-color: #fff; }
          .navbar { display: flex; background-color: #222; margin: 0; border-width: 0 0 1px; border-style: solid; border-color: #080808; }
          .navbar > * { margin: 0; padding: 15px; }
          .navbar * { line-height: 20px; color: #9d9d9d; }
          .navbar a { text-decoration: none; }
          .navbar a:hover, .navbar a:focus { color: #fff; }
          .navbar-header { font-size: 18px; }
          body > * { margin: 15px; padding: 0; }
          pre { padding: 10px; font-size: 13px; background-color: #f5f5f5; border: 1px solid #ccc; }
          h1, h2 { font-weight: 500; }
          a { color: #337ab7; }
          a:hover, a:focus { color: #23527c; }
        </style>
      </head>
      <body>
        <div class="navbar">
          <div class="navbar-header"><a href="/">JIRAlert</a></div>
          <div><a href="/config">Configuration</a></div>
          <div><a href="/metrics">Metrics</a></div>
          <div><a href="/debug/pprof">Profiling</a></div>
          <div><a href="{{ .DocsUrl }}">Help</a></div>
        </div>
        {{template "content" .}}
      </body>
      </html>
    {{- end }}

    {{ define "content.home" -}}
      <p>This is <a href="{{ .DocsUrl }}">JIRAlert</a>, a
        <a href="https://prometheus.io/docs/alerting/configuration/#webhook_config">webhook receiver</a> for
        <a href="https://prometheus.io/docs/alerting/alertmanager/">Prometheus Alertmanager</a>.
    {{- end }}

    {{ define "content.config" -}}
      <h2>Configuration</h2>
      <pre>{{ .Config }}</pre>
    {{- end }}

    {{ define "content.error" -}}
      <h2>Error</h2>
      <pre>{{ .Err }}</pre>
    {{- end }}
    `
)

type tdata struct {
	DocsUrl string

	// `/config` only
	Config string

	// `/error` only
	Err error
}

var (
	allTemplates   = template.Must(template.New("").Parse(templates))
	homeTemplate   = pageTemplate("home")
	configTemplate = pageTemplate("config")
	// errorTemplate  = pageTemplate("error")
)

func pageTemplate(name string) *template.Template {
	pageTemplate := fmt.Sprintf(`{{define "content"}}{{template "content.%s" .}}{{end}}{{template "page" .}}`, name)
	return template.Must(template.Must(allTemplates.Clone()).Parse(pageTemplate))
}

// HomeHandlerFunc is the HTTP handler for the home page (`/`).
func HomeHandlerFunc() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := homeTemplate.Execute(w, &tdata{
			DocsUrl: docsUrl,
		}); err != nil {
			w.WriteHeader(500)
		}
	}
}

// ConfigHandlerFunc is the HTTP handler for the `/config` page. It outputs the configuration marshaled in YAML format.
func ConfigHandlerFunc(config *config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := configTemplate.Execute(w, &tdata{
			DocsUrl: docsUrl,
			Config:  config.String(),
		}); err != nil {
			w.WriteHeader(500)
		}
	}
}
