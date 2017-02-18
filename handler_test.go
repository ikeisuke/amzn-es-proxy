package main

import (
  "testing"
  "net/http"
  "net/http/httptest"
  "net/url"
  "io/ioutil"
)

type mockSigner struct {
}

func (s *mockSigner) Sign(req *http.Request) error {
  req.Header.Add("Authorization", "mock token")
  return nil
}

func TestHandler(t *testing.T) {
  endpoint := "example.com"
  path     := "/path/to/resource"

  expectResponse, err := request("https://" + endpoint + path)
  if err != nil {
    t.Fatalf("%v", err);
  }
  expectBody, err := ioutil.ReadAll(expectResponse.Body)

  server := httptest.NewServer(NewHandler(endpoint, &mockSigner{}))
  defer server.Close()

  res, err := request(server.URL + path)
  if err != nil {
    t.Fatalf("%v", err);
  }
  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    t.Fatalf("%v", err);
  }
  if string(body) != string(expectBody) {
    t.Fatalf("%v", "not equal response")
  }
}

func request(targetURL string) (*http.Response, error) {
    req, err := http.NewRequest("GET", targetURL, nil)
    if err != nil {
      return nil, err
    }

    parsed, err := url.Parse(targetURL)
    if err != nil {
      return nil, err
    }
    sourceURL := parsed.Scheme + "://" + parsed.Host

    req.Header.Add("Accept", "*/*")
    req.Header.Add("Accept-Language", "ja,en-US;q=0.8,en;q=0.6")
    req.Header.Add("Connection", "keep-alive")
    req.Header.Add("Origin", sourceURL)
    req.Header.Add("Referer", sourceURL + "/path/to/referer")
    req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")

    res, err := http.DefaultClient.Do(req)
    if err != nil {
      return nil, err
    }
    return res, nil
}
