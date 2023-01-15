package render

import (
	"bytes"
	// "fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/aman1111068/bookings/pkg/config"
	"github.com/aman1111068/bookings/pkg/models"
)

//RenderTemplate renders templates using html/template
// func RenderTemplateWithoutCache(w http.ResponseWriter, tmpl string) {
// 	parsedTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.html")
// 	err := parsedTemplate.Execute(w, nil)
// 	if err != nil {
// 		fmt.Println("error parsing template:", err)
// 		return
// 	}
// }

// var tc = make(map[string]*template.Template)

//RenderTemplate renders templates using html/template
// func RenderTemplateWithCache(w http.ResponseWriter, t string) {
// 	var tmpl *template.Template
// 	var err error
// 	// check to see if we already have template in our cache
// 	_, inMap := tc[t]
// 	if !inMap {
// 		//need to create the template and add to cache
// 		log.Println("creating template and adding to cache")
// 		err = createTemplateCache(t)
// 		if err != nil {
// 			fmt.Println("error creating template cache:", err)
// 		}

// 	} else {
// 		// we have the template in the cache
// 		log.Println("using cached template")
// 	}

// 	tmpl = tc[t]
// 	err = tmpl.Execute(w, nil)
// 	if err != nil {
// 		fmt.Println("error executing template:", err)
// 		return
// 	}
// }

// func createTemplateCache(t string) error {
// 	// create an slice of template strings
// 	templates := []string{
// 		fmt.Sprintf("./templates/%s", t),
// 		"./templates/base.layout.html",
// 	}

// 	//parse the template
// 	tmpl, err := template.ParseFiles(templates...)
// 	if err != nil {
// 		fmt.Println("error parsing template:", err)
// 		return err
// 	}

// 	//add template to cache (map)
// 	tc[t] = tmpl

// 	return nil
// }

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

//RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateComplexTemplateCache()

	}

	// Calling CreateComplexTemplateCache separately to create template cache
	// tc, err := CreateComplexTemplateCache()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// get requested template from cache
	t, inMap := tc[tmpl]
	if !inMap {
		log.Fatal("Could not get the template from template cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing template to browser", err)
	}
}

func CreateComplexTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files named *.page.html from ./templates
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	// range through all the files(in pages slice) ending with *.page.html
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil

}
