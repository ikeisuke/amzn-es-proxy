package main

import (
  "log"
  "net/http"
  
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

func main() {
  domain := "test"
  sess, err := session.NewSession()
  if err != nil {
    log.Println("failed to create session,", err)
    return
  }
  es := elasticsearchservice.New(sess)
  params := &elasticsearchservice.DescribeElasticsearchDomainInput{
		DomainName: aws.String(domain), // Required
	}
	resp, err := es.DescribeElasticsearchDomain(params)

	if err != nil {
		log.Println(err.Error())
		return
	}
  endpoint := *resp.DomainStatus.Endpoint
  http.HandleFunc("/", handler(endpoint))
  http.ListenAndServe("localhost:8080", nil)
}
