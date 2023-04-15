package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/IDOMATH/portfolio/config"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var pathToTemplates = "./templates"
var app *config.AppConfig

func NewRenderer(a *config.AppConfig) {
	app = a
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string) error {
	var tc map[string]*template.Template
	if app.UseCache {
		// Get the template cache from app config
		tc = app.TemplateCache
	} else {
		// This is just for testing, so we rebuild the cache ever request
		tc, _ = CreateTemplateCache()
	}

	template, ok := tc[tmpl]
	if !ok {
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	err := template.Execute(buf, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser")
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return cache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return cache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return cache, err
			}
		}
		cache[name] = ts
	}
	return cache, nil
}
