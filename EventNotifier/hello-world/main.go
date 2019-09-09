package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/davidsbond/lux"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")

	//AWS
	sess, _ = session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: credentials.NewEnvCredentials(),
	})

	svc = s3.New(sess)
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayProxyResponse{}, ErrNoIP
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello, %v", string(ip)),
		StatusCode: 200,
	}, nil
}

func getUserEvents(w lux.ResponseWriter, r *lux.Request) {
	// userid := r.QueryStringParameters["id"]
	out, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("concertnotifier"),
		Key:    aws.String("users.json"),
	})

	if err != nil {
		fmt.Errorf("failed to download file, %v", err)
	}
	w.Write([]byte("Ana are mere"))
	fmt.Println(out)
}

func main() {
	router := lux.NewRouter()

	router.Handler("GET", getUserEvents).Queries("id", "*")

	lambda.Start(router)
}
