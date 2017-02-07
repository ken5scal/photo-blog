package main

import (
	"html/template"
	"net/http"
	"github.com/satori/go.uuid"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", index)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func index(res http.ResponseWriter, req *http.Request) {
	cookie := getCookie(res, req)
	tmpl.ExecuteTemplate(res, "index.html", cookie)
}
func getCookie(res http.ResponseWriter, request *http.Request) *http.Cookie {
	cookie, err := request.Cookie("session")
	if err != nil {
		sessionID := uuid.NewV4()
		cookie = &http.Cookie{
			Name: "session",
			Value: sessionID.String(),
		}
		http.SetCookie(res, cookie)
	}

	return cookie
}