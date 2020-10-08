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
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/comprehend/types"
)

type APIResponse struct {
	Message  string `json:"message"`
}

type Response events.APIGatewayProxyResponse

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
		comprehendClient = getComprehendClient()
	}

	input := &comprehend.DetectSentimentInput{
		LanguageCode: types.LanguageCodeJa,
		Text: aws.String(message),
	}
	res, err := comprehendClient.DetectSentiment(ctx, input)
	if err != nil {
		return "", err
	}
	return string(res.Sentiment), nil
}

func detectDominantLanguage(ctx context.Context, message string)(string, error) {
	if comprehendClient == nil {
		comprehendClient = getComprehendClient()
	}

	input := &comprehend.DetectDominantLanguageInput{
		Text: aws.String(message),
	}
	res, err := comprehendClient.DetectDominantLanguage(ctx, input)
	if err != nil {
		return "", err
	}
	results, err2 := json.Marshal(res.Languages)
	if err2 != nil {
		return "", err2
	}
	return string(results), nil
}

func detectEntities(ctx context.Context, message string)(string, error) {
	if comprehendClient == nil {
		comprehendClient = getComprehendClient()
	}

	input := &comprehend.DetectEntitiesInput{
		LanguageCode: types.LanguageCodeJa,
		Text: aws.String(message),
	}
	res, err := comprehendClient.DetectEntities(ctx, input)
	if err != nil {
		return "", err
	}
	results, err2 := json.Marshal(res.Entities)
	if err2 != nil {
		return "", err2
	}
	return string(results), nil
}

func detectKeyPhrases(ctx context.Context, message string)(string, error) {
	if comprehendClient == nil {
		comprehendClient = getComprehendClient()
	}

	input := &comprehend.DetectKeyPhrasesInput{
		LanguageCode: types.LanguageCodeJa,
		Text: aws.String(message),
	}
	res, err := comprehendClient.DetectKeyPhrases(ctx, input)
	if err != nil {
		return "", err
	}
	results, err2 := json.Marshal(res.KeyPhrases)
	if err2 != nil {
		return "", err2
	}
	return string(results), nil
}

func detectSyntax(ctx context.Context, message string)(string, error) {
	if comprehendClient == nil {
		comprehendClient = getComprehendClient()
	}

	input := &comprehend.DetectSyntaxInput{
		LanguageCode: types.SyntaxLanguageCodeEn,
		Text: aws.String(message),
	}
	res, err := comprehendClient.DetectSyntax(ctx, input)
	if err != nil {
		return "", err
	}
	results, err2 := json.Marshal(res.SyntaxTokens)
	if err2 != nil {
		return "", err2
	}
	return string(results), nil
}

func getComprehendClient() *comprehend.Client {
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		log.Print(err)
	}
	cfg.Region = os.Getenv("REGION")
	return comprehend.NewFromConfig(cfg)
}

func main() {
	lambda.Start(HandleRequest)
}
