package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

//go:embed templates
var templates embed.FS

//go:embed assets
var assets embed.FS

func (conf *Configuration) startHttpListener() {
	http.HandleFunc("/", conf.defaultHandler)
	l("info", fmt.Sprintf("Starting http server on port %d", conf.Http.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%04d", conf.Http.Port), nil))
}

func (conf *Configuration) defaultHandler(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Path
	if u == "" || u == "/" || strings.HasPrefix(u, "/index.") {
		u = "/index"
	}

	if f, err := templates.Open(templateFile(u)); err == nil {
		f.Close()
		conf.handleTemplates(u, w, r)
		return
	} else if f, err := assets.Open(assetFile(u)); err == nil {
		f.Close()
		conf.handleAssets(w, r)
		return
	}
	conf.renderNotFound(w, r)
	accessLog(r, http.StatusNotFound, "not found")
}

type IndexData struct {
	Version string
}

func (conf *Configuration) handleTemplates(uri string, w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFS(templates, templateFile(uri))
	if err != nil {
		accessLog(r, http.StatusInternalServerError, fmt.Sprintf("Error parsing template %s: %v", r.URL.Path, err))
		conf.renderServerError(w, r)
		return
	}

	switch uri {
	case "/index":
		err = t.Execute(w, IndexData{Version: version})
		if err != nil {
			accessLog(r, http.StatusInternalServerError, fmt.Sprintf("error executing template /index.gohtml: %v", err))
			conf.renderServerError(w, r)
			return
		}
	}
	accessLog(r, http.StatusOK, "ok")
}

func (conf *Configuration) handleAssets(w http.ResponseWriter, r *http.Request) {
	data, err := assets.ReadFile(assetFile(r.URL.Path))
	if err != nil {
		accessLog(r, http.StatusNotFound, err.Error())
		conf.renderNotFound(w, r)
		return
	}
	accessLog(r, 200, "")
	lc := strings.ToLower(r.RequestURI)
	switch {
	case strings.HasSuffix(lc, ".css"):
		w.Header().Add("Content-Type", "text/css")
	case strings.HasSuffix(lc, ".jpg"), strings.HasSuffix(lc, ".jpeg"):
		w.Header().Add("Content-Type", "image/jpeg")
	case strings.HasSuffix(lc, ".png"):
		w.Header().Add("Content-Type", "image/png")
	case strings.HasSuffix(lc, ".gif"):
		w.Header().Add("Content-Type", "image/gif")
	case strings.HasSuffix(lc, ".ico"):
		w.Header().Add("Content-Type", "image/x-icon")
	case strings.HasSuffix(lc, ".html"), strings.HasSuffix(lc, "htm"):
		w.Header().Add("Content-Type", "text/html")
	case strings.HasSuffix(lc, ".js"):
		w.Header().Add("Content-Type", "application/javascript")
	case strings.HasSuffix(lc, ".json"):
		w.Header().Add("Content-Type", "application/json")
	case strings.HasSuffix(lc, ".map"):
		w.Header().Add("Content-Type", "application/json")
	case strings.HasSuffix(lc, ".svg"):
		w.Header().Add("Content-Type", "image/svg+xml")
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		accessLog(r, http.StatusInternalServerError, "write data")
	} else {
		accessLog(r, http.StatusOK, "ok")
	}
}

func templateFile(u string) string {
	return "templates" + u + ".gohtml"
}

func assetFile(u string) string {
	return "assets" + u
}

type NotFoundData struct {
	Page string
}

func (conf *Configuration) renderNotFound(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFS(templates, templateFile("/404"))
	if err != nil {
		accessLog(r, http.StatusInternalServerError, fmt.Sprintf("Error parsing template %s: %v", r.URL.Path, err))
		conf.renderServerError(w, r)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	err = t.Execute(w, NotFoundData{Page: r.URL.Path})
	if err != nil {
		accessLog(r, http.StatusInternalServerError, fmt.Sprintf("error executing template /404.gohtml: %v", err))
		fmt.Fprintf(w, "Could not find the page you requested: %s.", r.RequestURI)
		return
	}
}

func (conf *Configuration) renderServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := fmt.Fprintf(w, "Internal Server Error: %s.", r.RequestURI)
	if err != nil {
		accessLog(r, http.StatusInternalServerError, "internal server error")
	}
}

func accessLog(r *http.Request, httpCode int, payload string) {
	switch httpCode {
	case http.StatusInternalServerError, http.StatusBadRequest:
		log.Printf("error %s %s, %d, %s", r.Method, r.RequestURI, httpCode, payload)
	case http.StatusNotFound, http.StatusUnauthorized:
		log.Printf("warning %s %s, %d, %s", r.Method, r.RequestURI, httpCode, payload)
	default:
		log.Printf("info %s %s, %d, %s", r.Method, r.RequestURI, httpCode, payload)
	}
}

func l(severity string, payload string) {
	log.Printf("error %s %s\n", severity, payload)
}
