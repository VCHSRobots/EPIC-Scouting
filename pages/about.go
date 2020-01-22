package pages

func init() {
	RegisterPage("/about", VerbGET, nil)
}
