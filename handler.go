package main

import (
  //"fmt"
  "net/http"
  "io/ioutil"
  "os"
  "os/user"
  "strings"

  "github.com/aws/aws-sdk-go/aws/signer/v4"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/sha1sum/aws_signing_client"
)

func handler(endpoint string) func(w http.ResponseWriter, r *http.Request) {
  u, _ := user.Current()
  creds := credentials.NewSharedCredentials(u.HomeDir + "/.aws/credentials", os.Getenv("AWS_PROFILE"))
  signer := v4.NewSigner(creds)
  client, _ := aws_signing_client.New(signer, nil, "es", os.Getenv("AWS_REGION"))
  return func(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    url := r.URL
    url.Host = endpoint
    url.Scheme = "https"
    esr, _ := http.NewRequest(r.Method, url.String(), r.Body)
    for k, vs := range r.Header {
      for _, v := range vs {
        esr.Header.Add(k, v)
      }
    }
    for _, c := range r.Cookies() {
      esr.AddCookie(c)
    }
    res, _ := client.Do(esr)
    b, err := ioutil.ReadAll(res.Body)
    if err == nil {
      for k, v := range res.Header {
        w.Header().Set(k, strings.Join(v, ", "))
      }
      w.WriteHeader(res.StatusCode)
      w.Write(b)
    }
  }
}
