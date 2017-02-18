package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

func handler(endpoint string) func(w http.ResponseWriter, r *http.Request) {
	u, _ := user.Current()
	creds := credentials.NewSharedCredentials(u.HomeDir+"/.aws/credentials", os.Getenv("AWS_PROFILE"))
	signer := v4.NewSigner(creds)
	client := http.DefaultClient
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
		var body io.ReadSeeker
		if req.Body != nil {
			buf, _ := ioutil.ReadAll(req.Body)
			body = bytes.NewReader(buf)
		}
		region := os.Getenv("AWS_REGION")
		if region == "" {
			region = os.Getenv("AWS_DEFAULT_REGION")
		}
		_, err := signer.Sign(req, body, "es", region, time.Now())
		res, err := client.Do(req)
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
