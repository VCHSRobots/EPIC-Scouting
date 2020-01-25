/*
Package web provides general tools for rendering the front-end web-pages.
THIS PACKAGE IS UNDER QUARANTINE AS OF 2020-01-23. IT IS NOT USED BY ANYTHING.
*/

package web

import (
	"EPIC-Scouting/lib/lumberjack"
	"bytes"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

const (
	// GET represents the HTTP method defined in RFC 7231 § 4.3.1. See https://tools.ietf.org/html/rfc7231 for more information.
	GET = "GET"
	// POST represents the HTTP method defined in RFC 7231 § 4.3.3.
	POST = "POST"
)

/*
A Verb is an HTTP method.
*/
type Verb string

/*
A Page is a struct with a Route, an HTTP verb, and any number of middleware handlers.
*/
type Page struct {
	Route    string
	Verb     Verb
	Handlers []gin.HandlerFunc
}

var pages []*Page

/*
RegisterPage adds a page's route, HTTP verb, and handlers to the pages array.
*/
func RegisterPage(route string, verb Verb, handlers ...gin.HandlerFunc) {
	log := lumberjack.New("RegisterPage")
	log.Debug("-----------------------> Loaded web! <-------------------------")
	print("Registered page: ")
	print(route)
	log.Debugf("Page %q handlers: %q", route, handlers) // TODO: TEMP
	if handlers == nil {
		log.Fatalf("Unable to register page with method %q, route %q: handler(s) is nil", verb, route)
	}
	if pages == nil {
		pages = make([]*Page, 0, 64) // TODO: Allow for page array of size N+1.
	}
	page := &Page{route, verb, handlers}
	pages = append(pages, page)
}

/*
LoadTemplates loads templates by glob pattern.
*/
func LoadTemplates(pattern string) {
	template.Must(template.New("").Delims("{{", "}}").Funcs(template.FuncMap).ParseGlob(pattern))
}

/*
RenderTemplate renders the requested template by filename.
*/
func RenderTemplate(name string) (*template.Template, error) {
	/*
		log := lumberjack.New("RenderTemplate")
		fileName := "./static/templates/" + name + ".tmpl"
		bytes, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Errorf("Unable to load template %q from %q: %s", name, fileName, err)
			return nil, err
		}
		template, err := template.New(name).Parse(string(bytes))
		if err != nil {
			log.Errorf("Unable to parse template %q from %q: %s", name, fileName, err)
		}
		return template, nil
	*/

}

/*
MakePage concatenates the listed templates and executes passed variables with them.
*/
func MakePage(data interface{}, templateNames ...string) ([]byte, error) {
	log := lumberjack.New("MakePage")
	templates := make([]*template.Template, 0, len(templateNames))
	for _, t := range templateNames {
		template, err := RenderTemplate(t)
		if err != nil {
			return []byte{}, err
		}
		templates = append(templates, template)
	}
	html := new(bytes.Buffer)
	for i, t := range templates {
		err := t.Execute(html, data)
		if err != nil {
			log.Errorf("Error executing template %q: %s", templateNames[i], err)
			return html.Bytes(), err
		}
	}
	return html.Bytes(), nil
}

/*
SendPage renders and then sends a page to the requesting client via gin-gonic.
*/
func SendPage(c *gin.Context, data interface{}, templateNames ...string) {
	html, err := MakePage(data, templateNames...)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, "text/html", html)
}
