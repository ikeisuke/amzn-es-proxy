package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

type SignerInterface interface {
	Sign(req *http.Request) error
}

type Signer struct {
	Region  string
	Service string
	Signer  *v4.Signer
}

func NewSigner(profile string, region string, service string) *Signer {
	ec2m := ec2metadata.New(session.New(), &aws.Config{
		HTTPClient: &http.Client{Timeout: time.Second},
	})
	creds := credentials.NewChainCredentials([]credentials.Provider{
		&credentials.SharedCredentialsProvider{
			Profile: profile,
		},
		&ec2rolecreds.EC2RoleProvider{
			Client: ec2m,
		},
	})
	s := new(Signer)
	s.Region = region
	s.Service = service
	s.Signer = v4.NewSigner(creds)
	return s
}

func (s *Signer) Sign(req *http.Request) error {
	var body io.ReadSeeker
	if req.Body != nil {
		buf, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}
		body = bytes.NewReader(buf)
	}
	region := s.Region
	_, err := s.Signer.Sign(req, body, "es", region, time.Now())
	return err
}
