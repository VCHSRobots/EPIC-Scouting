package web

/*
HeaderData defines the variables which may be passed to a page's header.
*/
type HeaderData struct {
	Title       string   // The title of a page.
	StyleSheets []string // Any additional CSS documents to include when rendering the page.
}
