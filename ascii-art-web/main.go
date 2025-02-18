package main

import (
	"html/template"
	"log"
	"net/http"
	"test/asciiart"
)

var templates map[string]*template.Template

var errorMessages = map[int]string{
	http.StatusNotFound:            "Page not found",
	http.StatusBadRequest:          "Invalid request",
	http.StatusForbidden:           "Access forbidden",
	http.StatusInternalServerError: "Internal server error",
	http.StatusMethodNotAllowed:    "Method not allowed",
}

type PageData struct {
	Art string
}

type ErrorData struct {
	ErrorCode    int
	ErrorMessage string
}

func main() {
	templates = map[string]*template.Template{
		"home":  template.Must(template.ParseFiles("frontend/index.html")),
		"error": template.Must(template.ParseFiles("frontend/error.html")),
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", artHandler)
	fs := http.FileServer(http.Dir("frontend/static"))
	http.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/static/" {
			data := ErrorData{ErrorCode: http.StatusForbidden, ErrorMessage: getErrorMessage(http.StatusForbidden)}
			renderTemplate(w, "error", data)
			return
		}
		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
	}))

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", recoverHandler(http.DefaultServeMux)))
}

func getErrorMessage(code int) string {
	if msg, exists := errorMessages[code]; exists {
		return msg
	}
	return "An unexpected error occurred"
}

func recoverHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				data := ErrorData{ErrorCode: http.StatusInternalServerError, ErrorMessage: getErrorMessage(http.StatusInternalServerError)}
				renderTemplate(w, "error", data)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		data := ErrorData{ErrorCode: http.StatusNotFound, ErrorMessage: getErrorMessage(http.StatusNotFound)}
		renderTemplate(w, "error", data)
		return
	}
	renderTemplate(w, "home", nil)
}

func artHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		text := r.FormValue("text")
		if text == "" || len(text) > 400 {
			data := ErrorData{ErrorCode: http.StatusBadRequest, ErrorMessage: getErrorMessage(http.StatusBadRequest)}
			renderTemplate(w, "error", data)
			return
		}
		bannerName := r.FormValue("banner")

		_, exists := ascii.Banners[bannerName]
		if !exists {
			data := ErrorData{ErrorCode: http.StatusBadRequest, ErrorMessage: getErrorMessage(http.StatusBadRequest)}
			renderTemplate(w, "error", data)
			return
		}
		asciiArt := ascii.GenAscii(text, bannerName)
		
		data := PageData{Art: asciiArt}
		renderTemplate(w, "home", data)
		return
	}
	renderTemplate(w, "home", nil)
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := templates[name]
	if !ok {
		log.Println("Template not found:", name)
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error rendering template:", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
