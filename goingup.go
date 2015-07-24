package goingup

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var app *App

// App represents the whole goingup application
type App struct {
	Options   AppOptions
	Pages     []Page
	Templates *template.Template
}

// NewApp creates a new App instance
func NewApp() *App {
	app = &App{
		Options: AppOptions{
			Port:            80,
			TemplateDir:     "templates",
			StaticAssetsDir: "static/",
			StaticAssetsURL: "/static/",
			LoginAction:     "/login",
			RegisterAction:  "/register",
		},
	}
	return app
}

// AddPage _
func (a *App) AddPage(url string, title string, tmpl string) error {
	if url == "" {
		return fmt.Errorf("Cannot create page with no URL")
	}

	if tmpl == "" {
		tmpl = "page"
	}

	a.Pages = append(a.Pages, Page{
		URL:      url,
		Title:    title,
		Template: tmpl,
	})

	return nil
}

// Run finalizes all options and calls the ListenAndServe function to serve
// requests
func (a *App) Run() {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir(a.Options.StaticAssetsDir))
	r.PathPrefix(a.Options.StaticAssetsURL).Handler(http.StripPrefix(a.Options.StaticAssetsURL, fs))

	for _, page := range a.Pages {
		r.HandleFunc(page.URL, makePageHandler(page))
	}

	a.Templates = template.Must(template.ParseGlob(a.Options.TemplateDir + "/*"))
	http.ListenAndServe(":"+strconv.Itoa(a.Options.Port), newLogHandler(r))
}