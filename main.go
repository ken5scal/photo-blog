package main

import (
	"html/template"
	"net/http"
	"github.com/satori/go.uuid"
	"strings"
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
	cookie = appendCookieValue(cookie, res)
	xs := strings.Split(cookie.Value, "|")
	tmpl.ExecuteTemplate(res, "index.html", xs)
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

func appendCookieValue(cookie *http.Cookie, res http.ResponseWriter) *http.Cookie {
	p1 := "dog.jpg"
	p2 := "cat.jpg"
	p3 := "rabbit.jpg"

	value := cookie.Value
	if !strings.Contains(value, p1) {
		value += "|" + p1
	}
	if !strings.Contains(value, p2) {
		value += "|" + p2
	}
	if !strings.Contains(value, p3) {
		value += "|" + p3
	}
	cookie.Value = value
	http.SetCookie(res, cookie)
	return cookie
}