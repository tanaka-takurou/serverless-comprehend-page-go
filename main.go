package main

import (
	"io"
	"log"
	"bytes"
	"context"
	"io/ioutil"
	"html/template"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type PageData struct {
	Title string
	Api   string
}

type ConstantData struct {
	Title string `json:"title"`
	Api   string `json:"api"`
}

type Response events.APIGatewayProxyResponse

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	tmp := template.New("tmp")
	var dat PageData
	r := request.Resource
	funcMap := template.FuncMap{
		"safehtml": func(text string) template.HTML { return template.HTML(text) },
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int { return a / b },
	}
	buf := new(bytes.Buffer)
	fw := io.Writer(buf)
	jsonString, _ := ioutil.ReadFile("constant/constant.json")
	constant := new(ConstantData)
	json.Unmarshal(jsonString, constant)
	dat.Api = constant.Api
	if r == "/detect/dominantlanguage" {
		dat.Title = constant.Title + " | Detect DominantLanguage"
		tmp = template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/index_dominantlanguage.html", "templates/view.html", "templates/header.html"))
	} else if r == "/detect/entities" {
		dat.Title = constant.Title + " | Detect Entities"
		tmp = template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/index_entities.html", "templates/view.html", "templates/header.html"))
	} else if r == "/detect/keyphrases" {
		dat.Title = constant.Title + " | Detect KeyPhrases"
		tmp = template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/index_keyphrases.html", "templates/view.html", "templates/header.html"))
	} else if r == "/detect/syntax" {
		dat.Title = constant.Title + " | Detect Syntax"
		tmp = template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/index_syntax.html", "templates/view.html", "templates/header.html"))
	} else {
		dat.Title = constant.Title + " | Detect Sentiment"
		tmp = template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/index.html", "templates/view.html", "templates/header.html"))
	}
	if e := tmp.ExecuteTemplate(fw, "base", dat); e != nil {
		log.Fatal(e)
	} else {
		log.Print("Event received.")
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
