package main

import (
	"io"
	"os"
	"log"
	"bytes"
	"embed"
	"context"
	"html/template"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type PageData struct {
	Title   string
	ApiPath string
}

type Response events.APIGatewayProxyResponse

//go:embed templates
var templateFS embed.FS

const title string = "Sample Comprehend Page"

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	tmp := template.New("tmp")
	var dat PageData
	p := request.PathParameters
	funcMap := template.FuncMap{
		"safehtml": func(text string) template.HTML { return template.HTML(text) },
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int { return a / b },
	}
	buf := new(bytes.Buffer)
	fw := io.Writer(buf)
	dat.ApiPath = os.Getenv("API_PATH")
	if p["proxy"] == "detect-dominantlanguage" {
		dat.Title = title + " | Detect DominantLanguage"
		tmp = template.Must(template.New("").Funcs(funcMap).ParseFS(templateFS, "templates/index_dominantlanguage.html", "templates/view.html", "templates/header.html"))
	} else if p["proxy"] == "detect-entities" {
		dat.Title = title + " | Detect Entities"
		tmp = template.Must(template.New("").Funcs(funcMap).ParseFS(templateFS, "templates/index_entities.html", "templates/view.html", "templates/header.html"))
	} else if p["proxy"] == "detect-keyphrases" {
		dat.Title = title + " | Detect KeyPhrases"
		tmp = template.Must(template.New("").Funcs(funcMap).ParseFS(templateFS, "templates/index_keyphrases.html", "templates/view.html", "templates/header.html"))
	} else if p["proxy"] == "detect-syntax" {
		dat.Title = title + " | Detect Syntax"
		tmp = template.Must(template.New("").Funcs(funcMap).ParseFS(templateFS, "templates/index_syntax.html", "templates/view.html", "templates/header.html"))
	} else {
		dat.Title = title + " | Detect Sentiment"
		tmp = template.Must(template.New("").Funcs(funcMap).ParseFS(templateFS, "templates/index.html", "templates/view.html", "templates/header.html"))
	}
	if e := tmp.ExecuteTemplate(fw, "base", dat); e != nil {
		log.Fatal(e)
	}
	res := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(buf.Bytes()),
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}
	return res, nil
}

func main() {
	lambda.Start(HandleRequest)
}
