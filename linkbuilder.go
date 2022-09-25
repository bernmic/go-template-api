package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Links map[string]string

type LinkBuilder struct {
	Links   map[string]string
	request *http.Request
}

func New(r *http.Request, init ...string) *LinkBuilder {
	if len(init)%2 != 0 {
		log.Printf("odd number of init values for new links collection")
	}

	l := LinkBuilder{
		Links:   make(map[string]string),
		request: r,
	}
	for i := 0; i < len(init)/2; i += 2 {
		l.Links[init[i]] = init[i+1]
	}
	return &l
}

func (l *LinkBuilder) Add(name string, uri string) {
	// use what's in the request
	p := l.baseUrl()
	l.Links[name] = p + uri
}

func (l *LinkBuilder) baseUrl() string {
	// check for "forwarded" header
	h := l.request.Header.Get("forwarded")
	if h != "" {
		p, err := l.forwarded(h)
		if err == nil {
			return p
		}
		log.Printf("warning: parsing forwarded header '%s': %v", h, err)
	}

	// check for "x-forwarded-*" header
	h = l.request.Header.Get("x-forwarded-for")
	if h != "" {
		p, err := l.xforwarded(l.request.Header)
		if err == nil {
			return p
		}
		log.Printf("warning: parsing x-forwarded header '%s': %v", h, err)
	}

	p := l.request.Host
	if l.request.TLS == nil {
		p = "http://" + p
	} else {
		p = "https://" + p
	}
	return p
}

func (l *LinkBuilder) forwarded(h string) (string, error) {
	proxies := strings.Split(h, ",")
	if len(proxies) == 0 {
		return "", errors.New("forwarded header has no entries")
	}
	parts := strings.Split(proxies[0], ";")
	if len(parts) == 0 {
		return "", errors.New("no data in forwarded header")
	}
	f := ""
	p := "http"
	for _, e := range parts {
		kv := strings.Split(strings.TrimSpace(e), "=")
		if len(kv) != 2 {
			return "", fmt.Errorf("illegal value '%s' in forwarded header", e)
		}
		switch strings.ToLower(kv[0]) {
		case "for":
			f = strings.TrimSpace(kv[1])
		case "proto":
			p = strings.TrimSpace(kv[1])
		}
	}
	if f == "" {
		return "", errors.New("missing for value in forwarded header")
	}
	return fmt.Sprintf("%s://%s", p, f), nil
}

func (l *LinkBuilder) xforwarded(header http.Header) (string, error) {
	//f := header.Get("x-forwarded-for")
	//p := header.Get("x-forwarded-proto")
	return "", errors.New("usage of x-forwarded-for not implemented yet")
}
