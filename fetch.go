package ssmgo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"strings"
)

var parameters []*Parameter

// ParameterFetcher interface for dealing with parameter retrieval
type ParameterFetcher interface {
	Initialize(path string, profile string) error
	GetParameterByName(name string) *Parameter
}

// LiveParameterFetcher implements the Parameter Fetcher as it relates to AWS
type LiveParameterFetcher struct {

}

// initialize loads up the parameters from aws
// and locally caches them for later retrieval
func (a *LiveParameterFetcher) Initialize(path string, profile string) error {
	var input =   &ssm.GetParametersByPathInput{}
	var config = &aws.Config{Region: aws.String("us-west-2")}
	//var session, err = session.NewSession(config)
	var sess *session.Session
	var err error = nil

	if profile == "" {
		sess, err = session.NewSession(config)
	} else {
		sess, err = session.NewSessionWithOptions(session.Options{
			Profile: profile,
			Config:  *config,
		})
	}

	if err != nil {
		return err
	}

	ssmsvc := ssm.New(sess)
	input.SetWithDecryption(false)
	input.SetPath(path)
	req, err :=  ssmsvc.GetParametersByPath(input)

	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, p := range req.Parameters{
		split := strings.Split(*p.Name, "/")
		name := strings.ToUpper(split[len(split)-1])

		parameters = append(parameters, &Parameter{
			Name: name,
			Value: *p.Value,
		})
	}

	return nil
}

// GetParameterByName returns a pointer to the Parameter
// as found by the Name property
func (a *LiveParameterFetcher)GetParameterByName(name string) *Parameter {
	for _, p := range parameters {
		if p.Name == name {
			return p
		}
	}

	return nil
}


