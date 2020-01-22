/*
Package pages provides general tools for rendering the front-end web-pages.
*/
package pages

import (
	"EPIC-Scouting/lib/lumberjack"
	"bytes"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

const (
	// VerbGET represents the HTTP method defined in RFC 7231 § 4.3.1. See https://tools.ietf.org/html/rfc7231 for more information.
	VerbGET = "GET"
	// VerbPOST represents the HTTP method defined in RFC 7231 § 4.3.3.
	VerbPOST = "POST"
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

var log *lumberjack.Lumberjack
var pages []*Page

func init() {
	log = lumberjack.New("Pages")
	// TODO: Check if static directories exist.
}

/*
RegisterPage adds a page's route, HTTP verb, and handlers to the pages array.
*/
func RegisterPage(route string, verb Verb, handlers ...gin.HandlerFunc) {
	if pages == nil {
		pages = make([]*Page, 0, 64) // TODO: Allow for page array of size N+1.
	}
	page := &Page{route, verb, handlers}
	pages = append(pages, page)
}

/*
GetPages returns all pages loaded into the pages variable.
*/
func GetPages() []*Page {
	return pages
}

/*
RenderTemplate renders the requested template by filename.
*/
func RenderTemplate(name string) (*template.Template, error) {
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
}

/*
MakePage concatenates the listed templates and executes passed variables with them.
*/
func MakePage(data interface{}, templateNames ...string) ([]byte, error) {
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
SendPage sends the rendered page to the requesting client via gin-gonic.
*/
func SendPage(c *gin.Context, data interface{}, templateNames ...string) {
	html, err := MakePage(data, templateNames...)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, "text/html", html)
}
