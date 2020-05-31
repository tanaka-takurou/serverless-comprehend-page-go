package main

import (
	"fmt"
	"log"
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type APIResponse struct {
	Message  string `json:"message"`
}

type Response events.APIGatewayProxyResponse

const layout       string = "2006-01-02 15:04"
const languageCode string = "en"

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var jsonBytes []byte
	var err error
	d := make(map[string]string)
	json.Unmarshal([]byte(request.Body), &d)
	if v, ok := d["action"]; ok {
		switch v {
		case "detectsentiment" :
			if m, ok := d["message"]; ok {
				r, e := detectSentiment(m)
				if e != nil {
					err = e
				} else {
					jsonBytes, _ = json.Marshal(APIResponse{Message: r})
				}
			}
		}
	}
	log.Print(request.RequestContext.Identity.SourceIP)
	if err != nil {
		log.Print(err)
		jsonBytes, _ = json.Marshal(APIResponse{Message: fmt.Sprint(err)})
		return Response{
			StatusCode: 500,
			Body: string(jsonBytes),
		}, nil
	}
	return Response {
		StatusCode: 200,
		Body: string(jsonBytes),
	}, nil
}

func detectSentiment(message string)(string, error) {
	svc := comprehend.New(session.New(), &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})

	input := &comprehend.DetectSentimentInput{
		LanguageCode: aws.String(languageCode),
		Text:  aws.String(message),
	}
	res, err := svc.DetectSentiment(input)
	if err != nil {
		return "", err
	}
	return aws.StringValue(res.Sentiment), nil
}

func main() {
	lambda.Start(HandleRequest)
}
