package main

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gambit/awsgo"
	"github.com/gambit/database"
	"github.com/gambit/handlers"
)

func main() {
	lambda.Start(LambdaExecute)
}

func LambdaExecute(context context.Context, request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
	awsgo.InitilizeAWS()

	if !ParametersValidate() {
		panic("Error na consulta de parâmetros")
	}

	var res *events.APIGatewayProxyResponse
	path := strings.Replace(request.RawPath, os.Getenv("UrlPrefix"), "", -1)
	method := request.RequestContext.HTTP.Method
	body := request.Body
	header := request.Headers

	database.ReadSecret()

	status, message := handlers.Handlers(path, method, body, header, request)

	headersResp := map[string]string{
		"Content-Type": "application/json",
	}

	res = &events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(message),
		Headers:    headersResp,
	}

	return res, nil

}

func ParametersValidate() bool {
	_, getParameter := os.LookupEnv("SecretName")
	if !getParameter {
		return getParameter
	}

	_, getParameter = os.LookupEnv("UrlPrefix")
	if !getParameter {
		return getParameter
	}

	return getParameter
}
