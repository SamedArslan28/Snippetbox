package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"snippetbox.samedarslan28.net/internal/models"
	"snippetbox.samedarslan28.net/ui"
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
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Use fs.Glob() to get a slice of all file paths in the ui.Files embedded
	// filesystem which match the pattern 'html/pages/*.tmpl'.
	pages, err := fs.Glob(ui.Files, "html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// Create a slice containing the filepath patterns for the templates to parse.
		patterns := []string{
			"html/base.gohtml",
			"html/partials/*.gohtml",
			page,
		}

		// Parse the template files from the embedded filesystem using ParseFS.
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
