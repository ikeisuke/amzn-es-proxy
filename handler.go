package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func handler(endpoint string) func(w http.ResponseWriter, r *http.Request) {
  profile := os.Getenv("AWS_PROFILE")
  region := os.Getenv("AWS_REGION")
  if region == "" {
    region = os.Getenv("AWS_DEFAULT_REGION")
  }
	signer := NewSigner(profile, region, "es")
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL
		url.Host = endpoint
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
		err := signer.Sign(req)
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
}
