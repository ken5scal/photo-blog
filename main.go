package main

import (
	"html/template"
	"net/http"
	"github.com/satori/go.uuid"
	"strings"
	"fmt"
	"crypto/sha1"
	"io"
	"os"
	"path/filepath"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", index)
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}

func index(res http.ResponseWriter, req *http.Request) {
	cookie := getCookie(res, req)

	if req.Method == http.MethodPost {
		mf, fh, err := req.FormFile("nf")
		if err != nil {
			fmt.Println(err)
		}
		defer mf.Close()

		// create sha
		h := sha1.New()
		io.Copy(h, mf)
		ext := strings.Split(fh.Filename, ".")[1]
		fname := fmt.Sprintf("%x", h.Sum(nil)) + "." + ext

		// create new file
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		path := filepath.Join(currentDir, "public", "pics", fname)
		newFile, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
		}
		defer newFile.Close()

		// copy
		mf.Seek(0, 0)
		io.Copy(newFile, mf)
		appendCookieValue(cookie, res, fname)
	}

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

func appendCookieValue(cookie *http.Cookie, res http.ResponseWriter, fname string)  {

	value := cookie.Value
	if !strings.Contains(value, fname) {
		value += "|" + fname
	}

	cookie.Value = value
	http.SetCookie(res, cookie)
}