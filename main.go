package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	flag "github.com/ogier/pflag"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
)

// options
//   --endpoint (elasticsearch service endpoint)
//   --listen (ip:port)
//   --domain (elasticsearch service domain)
//   --region
//   --profile
// env
//    AWS_PROFILE
//    AWS_ACCESS_KEY_ID
//    AWS_SECRET_ACCESS_KEY
//    AWS_REGION

var endpoint string
var domain string
var listen string
var region string
var profile string
var help bool

func main() {
	flag.StringVar(&endpoint, "endpoint", "", "The Amazon Elasticsearch Service Endpoint to use. ex. search-[domain]-xxxxxx.[region].es.amazonaws.com")
	flag.StringVar(&domain, "domain", "", "The Amazon Elasticsearch Service Domain to retrive endpoint.")
	flag.StringVar(&listen, "listen", "127.0.0.1:9200", "Listen on host:port. If you want to connect any address. :9200 or 0.0.0.0:9200")
	flag.StringVar(&profile, "profile", "", "Use a specific profile from your aws credential file.")
	flag.StringVar(&region, "region", "", "The region to use. Overrides config/env settings.")
	flag.BoolVar(&help, "help", false, "show this message.")
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if endpoint == "" && domain == "" {
		fmt.Fprintln(os.Stderr, "Either --endpoint or --domain is required.")
		flag.Usage()
		os.Exit(-1)
	}

	if endpoint != "" && domain != "" {
		fmt.Fprintln(os.Stderr, "Both --endpoint and --domain is not be enabled at once.")
		flag.Usage()
		os.Exit(-1)
	}

	if profile != "" {
		os.Setenv("AWS_PROFILE", profile)
	}

	if region != "" {
		os.Setenv("AWS_REGION", region)
	}

	if endpoint == "" {
		sess, err := session.NewSession()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to create aws session.")
			os.Exit(-1)
		}
		es := elasticsearchservice.New(sess)
		params := &elasticsearchservice.DescribeElasticsearchDomainInput{
			DomainName: aws.String(domain),
		}
		resp, err := es.DescribeElasticsearchDomain(params)
		if err != nil {
			log.Println(err.Error())
			return
		}
		processing := *resp.DomainStatus.Processing
		endpoint = *resp.DomainStatus.Endpoint
		if processing && endpoint == "" {
			fmt.Fprintln(os.Stderr, "This domain is being initialized. Try again later.")
			return
		}
		fmt.Fprintln(os.Stdout, "Using endpoint "+endpoint)
	}

	fmt.Fprintln(os.Stdout, "Listen "+listen)
	http.HandleFunc("/", handler(endpoint))
	http.ListenAndServe(listen, nil)
}
