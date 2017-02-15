package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "os"
  "os/user"

  "github.com/aws/aws-sdk-go/aws/signer/v4"
  "github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/sha1sum/aws_signing_client"
)

func handler(endpoint string) func(w http.ResponseWriter, r *http.Request) {
  u, _ := user.Current()
  creds := credentials.NewSharedCredentials(u.HomeDir + "/.aws/credentials", os.Getenv("AWS_PROFILE"))
  signer := v4.NewSigner(creds)
  client, _ := aws_signing_client.New(signer, nil, "es", "ap-northeast-1")
  return func(w http.ResponseWriter, r *http.Request) {
    esr, _ := http.NewRequest(r.Method, "https://" + endpoint + "/", nil)
    res, _ := client.Do(esr)
    b, err := ioutil.ReadAll(res.Body)
    if err == nil {
      fmt.Fprintf(w, string(b))
    }
  }
}
