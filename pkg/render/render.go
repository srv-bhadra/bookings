package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"practicewebapplication/pkg/config"
	"practicewebapplication/pkg/handlers/models"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}

// RenderTemplate renders template using HTML template
// COMPLEX TEMPLATE CACHE SYSTEM
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	//tc is a variable to hold template cache
	var tc map[string]*template.Template

	if app.UseCache {
		//get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	/*
		//create a template cache
		tc, err := CreateTemplateCache()
		if err != nil {
			log.Fatal(err)
		}
	*/

	//get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	//render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing template to browser ", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	//myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{} //same as above

	//get all the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	//Glob returns the names of all files matching pattern or nil if there is no matching file
	if err != nil {
		return myCache, err
	}

	//range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		//returns *.page.tmpl

		ts, err := template.New(name).ParseFiles(page)
		//ParseFiles creates a new Template and parses the template definitions from the named files
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			//ParseGlob creates a new Template and parses the template definitions from the files identified by the pattern
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}

/* SIMPLE TEMPLATE CACHE SYSTEM

// RenderTemplate renders template using HTML template
func RenderTemplateTest(w http.ResponseWriter, tmpl string) {
	parsedTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl")
	err := parsedTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("error parsing template : ", err)
		return
	}
}

var tc = make(map[string]*template.Template)

func RenderTemplate(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	//check to see if we already have the template in our cache
	_, inMap := tc[t]
	if !inMap {
		//need to create a template
		log.Println("Creating template and adding to cache")
		err = createTemplateCache(t)
		if err != nil {
			log.Println(err)
		}
	} else {
		//we have the template in the cache
		log.Println("using cached template")
	}

	tmpl = tc[t]

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

func createTemplateCache(t string) error {

	//add more templates manually here
	templates := []string {
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.tmpl",
	}

	//parse the template
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	//add template to cache
	tc[t] = tmpl

	return nil
}

*/
