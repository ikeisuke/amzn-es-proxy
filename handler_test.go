package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type mockSigner struct {
}

func (s *mockSigner) Sign(req *http.Request) error {
	req.Header.Add("Authorization", "mock token")
	return nil
}

func remoteServer() *httptest.Server {
	return httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello HTTP Test\n---\n")
		fmt.Fprintf(w, "%s %s %s\n", r.Method, r.RequestURI, r.Proto)
		fmt.Fprintf(w, "ContentLength: %d\n", r.ContentLength)
		fmt.Fprintf(w, "HEADER:\n")
		keys := []string{
			"Accept",
			"Accept-Encoding",
			"Accept-Language",
			"Content-Length",
			"User-Agent",
			"Origin",
			"Referer",
			"Connection",
		}
		for i := range keys {
			key := keys[i]
			value := r.Header.Get(key)
			fmt.Fprintf(w, "\t%v: %v\n", key, value)
		}
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Fprintf(w, "%s\n", string(body))
	}))
}

func TestGet(t *testing.T) {
	server := remoteServer()
	defer server.Close()

	parsed, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("%v", err)
	}

	endpoint := parsed.Host
	path := "/path/to/resource"

	expectResponse, err := get("https://" + endpoint + path)
	if err != nil {
		t.Fatalf("%v", err)
	}
	expectBody, err := ioutil.ReadAll(expectResponse.Body)

	handler := NewHandler(endpoint, &mockSigner{})
	handler.Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	proxy := httptest.NewServer(handler)
	defer proxy.Close()

	res, err := get(proxy.URL + path)
	if err != nil {
		t.Fatalf("%v", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("%v", err)
	}
	expected := string(expectBody)
	expected = strings.Replace(expected, "Connection: keep-alive", "Connection: ", 1)
	data := string(body)
	if data != expected {
		t.Log(string(expectBody))
		t.Log(string(body))
		t.Fatalf("%v", "not equal response")
	}
}

func TestPost(t *testing.T) {
	server := remoteServer()
	defer server.Close()

	parsed, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("%v", err)
	}

	endpoint := parsed.Host
	path := "/path/to/resource"

	expectResponse, err := post("https://" + endpoint + path)
	if err != nil {
		t.Fatalf("%v", err)
	}
	expectBody, err := ioutil.ReadAll(expectResponse.Body)

	handler := NewHandler(endpoint, &mockSigner{})
	handler.Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	proxy := httptest.NewServer(handler)
	defer proxy.Close()

	res, err := post(proxy.URL + path)
	if err != nil {
		t.Fatalf("%v", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("%v", err)
	}
	expected := string(expectBody)
	expected = strings.Replace(expected, "Connection: keep-alive", "Connection: ", 1)
	data := string(body)
	if data != expected {
		t.Log(string(expectBody))
		t.Log(string(body))
		t.Fatalf("%v", "not equal response")
	}
}

func get(rUrl string) (*http.Response, error) {
	return send("GET", rUrl, "")
}

func post(rUrl string) (*http.Response, error) {
	return send("POST", rUrl, "'post body'")
}

func send(method string, rUrl string, body string) (*http.Response, error) {
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, rUrl, reader)
	if err != nil {
		return nil, err
	}

	parsed, err := url.Parse(rUrl)
	if err != nil {
		return nil, err
	}
	sourceURL := parsed.Scheme + "://" + parsed.Host

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", "ja,en-US;q=0.8,en;q=0.6")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Origin", sourceURL)
	req.Header.Add("Referer", sourceURL+"/path/to/referer")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
