package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type Handler struct {
	Endpoint string
	Client   *http.Client
	s        SignerInterface
}

func NewHandler(endpint string, signer SignerInterface) *Handler {
	return &Handler{
		Endpoint: endpint,
		Client:   http.DefaultClient,
		s:        signer,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	url.Host = h.Endpoint
	url.Scheme = "https"
	req, _ := http.NewRequest(r.Method, url.String(), r.Body)
	req.ContentLength = r.ContentLength
	origin := "http://" + r.Host
	replace := "https://" + h.Endpoint
	for k, vs := range r.Header {
		if k == "Connection" {
			continue
		}
		for _, v := range vs {
			req.Header.Add(k, strings.Replace(v, origin, replace, -1))
		}
	}
	err := h.s.Sign(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := h.Client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := ioutil.ReadAll(res.Body)
	if err == nil {
		for k, v := range res.Header {
			s := strings.Join(v, ", ")
			w.Header().Set(k, strings.Replace(s, origin, replace, -1))
		}
		w.WriteHeader(res.StatusCode)
		w.Write(b)
	}
}
