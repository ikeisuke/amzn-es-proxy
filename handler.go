package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	Endpoint string
	s        *Signer
}

func NewHandler(endpint string, profile string, region string) *Handler {
	return &Handler{
		Endpoint: endpint,
		s:        NewSigner(profile, region, "es"),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	url.Host = h.Endpoint
	url.Scheme = "https"
	req, _ := http.NewRequest(r.Method, url.String(), r.Body)
	req.ContentLength = r.ContentLength
	for k, vs := range r.Header {
		if k == "Connection" {
			continue
		}
		for _, v := range vs {
			req.Header.Add(k, strings.Replace(v, "http://localhost:8080", "https://"+endpoint, -1))
		}
	}
	err := h.s.Sign(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := ioutil.ReadAll(res.Body)
	if err == nil {
		for k, v := range res.Header {
			s := strings.Join(v, ", ")
			w.Header().Set(k, strings.Replace(s, "http://localhost:8080", "https://"+endpoint, -1))
		}
		w.WriteHeader(res.StatusCode)
		w.Write(b)
	}
}
