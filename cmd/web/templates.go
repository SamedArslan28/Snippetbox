package main

import (
	"html/template"
	"path/filepath"
	"snippetbox.samedarslan28.net/internal/models"
	"time"
)

type templateData struct {
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	CurrentYear     int
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)

	// Get all the pages (templates) from the "pages" directory
	pages, err := filepath.Glob("./ui/html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}

	// Iterate over each page and parse the templates
	for _, page := range pages {
		// Get the name of the page (file name) to use as the template name in the cache
		name := filepath.Base(page)

		// Parse the base template
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.gohtml")
		if err != nil {
			return nil, err
		}

		// Parse the partials templates and append them to the base template
		ts, err = ts.ParseGlob("./ui/html/partials/*.gohtml")
		if err != nil {
			return nil, err
		}

		// Parse the specific page template and append it to the others
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template to the cache
		cache[name] = ts
	}

	// Return the cache with all templates
	return cache, nil
}
