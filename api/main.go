package main

import (
	"os"
	"fmt"
	"log"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
)

type APIResponse struct {
	Message  string `json:"message"`
}

type Response events.APIGatewayProxyResponse

var cfg aws.Config
var comprehendClient *comprehend.Client

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
				r, e := detectSentiment(ctx, m)
				if e != nil {
					err = e
				} else {
					jsonBytes, _ = json.Marshal(APIResponse{Message: r})
				}
			}
		case "detectdominantlanguage" :
			if m, ok := d["message"]; ok {
				r, e := detectDominantLanguage(ctx, m)
				if e != nil {
					err = e
				} else {
					jsonBytes, _ = json.Marshal(APIResponse{Message: r})
				}
			}
		case "detectentities" :
			if m, ok := d["message"]; ok {
				r, e := detectEntities(ctx, m)
				if e != nil {
					err = e
				} else {
					jsonBytes, _ = json.Marshal(APIResponse{Message: r})
				}
			}
		case "detectkeyphrases" :
			if m, ok := d["message"]; ok {
				r, e := detectKeyPhrases(ctx, m)
				if e != nil {
					err = e
				} else {
					jsonBytes, _ = json.Marshal(APIResponse{Message: r})
				}
			}
		case "detectsyntax" :
			if m, ok := d["message"]; ok {
				r, e := detectSyntax(ctx, m)
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

func detectSentiment(ctx context.Context, message string)(string, error) {
	if comprehendClient == nil {
		comprehendClient = comprehend.New(cfg)
	}

	input := &comprehend.DetectSentimentInput{
		LanguageCode: comprehend.LanguageCodeJa,
		Text: aws.String(message),
	}
	req := comprehendClient.DetectSentimentRequest(input)
	res, err := req.Send(ctx)
	if err != nil {
		return "", err
	}
	return string(res.DetectSentimentOutput.Sentiment), nil
}

func detectDominantLanguage(ctx context.Context, message string)(string, error) {
	if comprehendClient == nil {
		comprehendClient = comprehend.New(cfg)
	}

	input := &comprehend.DetectDominantLanguageInput{
		Text: aws.String(message),
	}
	req := comprehendClient.DetectDominantLanguageRequest(input)
	res, err := req.Send(ctx)
	if err != nil {
		return "", err
	}
	results, err2 := json.Marshal(res.DetectDominantLanguageOutput.Languages)
	if err2 != nil {
		return "", err2
	}
	return string(results), nil
}

func detectEntities(ctx context.Context, message string)(string, error) {
	if comprehendClient == nil {
		comprehendClient = comprehend.New(cfg)
	}

	input := &comprehend.DetectEntitiesInput{
		LanguageCode: comprehend.LanguageCodeJa,
		Text: aws.String(message),
	}
	req := comprehendClient.DetectEntitiesRequest(input)
	res, err := req.Send(ctx)
	if err != nil {
		return "", err
	}
	results, err2 := json.Marshal(res.DetectEntitiesOutput.Entities)
	if err2 != nil {
		return "", err2
	}
	return string(results), nil
}

func detectKeyPhrases(ctx context.Context, message string)(string, error) {
	if comprehendClient == nil {
		comprehendClient = comprehend.New(cfg)
	}

	input := &comprehend.DetectKeyPhrasesInput{
		LanguageCode: comprehend.LanguageCodeJa,
		Text: aws.String(message),
	}
	req := comprehendClient.DetectKeyPhrasesRequest(input)
	res, err := req.Send(ctx)
	if err != nil {
		return "", err
	}
	results, err2 := json.Marshal(res.DetectKeyPhrasesOutput.KeyPhrases)
	if err2 != nil {
		return "", err2
	}
	return string(results), nil
}

func detectSyntax(ctx context.Context, message string)(string, error) {
	if comprehendClient == nil {
		comprehendClient = comprehend.New(cfg)
	}

	input := &comprehend.DetectSyntaxInput{
		LanguageCode: comprehend.SyntaxLanguageCodeEn,
		Text: aws.String(message),
	}
	req := comprehendClient.DetectSyntaxRequest(input)
	res, err := req.Send(ctx)
	if err != nil {
		return "", err
	}
	results, err2 := json.Marshal(res.DetectSyntaxOutput.SyntaxTokens)
	if err2 != nil {
		return "", err2
	}
	return string(results), nil
}

func init() {
	var err error
	cfg, err = external.LoadDefaultAWSConfig()
	cfg.Region = os.Getenv("REGION")
	if err != nil {
		log.Print(err)
	}
}

func main() {
	lambda.Start(HandleRequest)
}
