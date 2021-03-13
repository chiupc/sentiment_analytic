package main

import (
	language "cloud.google.com/go/language/apiv1"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

//Attempts to authenticate to twitter developer api, if fails then don't accept request
type GoogleNLPApiHandler struct {
	Transport *http.Transport
	NetClient *http.Client
	GoogleApiClient *language.Client
	TimeOut int
}

func NewGoogleNLPApiHandler() *GoogleNLPApiHandler {
	gah := &GoogleNLPApiHandler{
		Transport: &http.Transport{
			Proxy:                  http.ProxyFromEnvironment,
			TLSHandshakeTimeout:    60 * time.Second,
			MaxIdleConns:           0,
			MaxIdleConnsPerHost:    100,
			MaxConnsPerHost:        100,
			IdleConnTimeout:        time.Duration(getEnvInt("GOOGLE_API_TIMEOUT")) *  time.Second,
		},
		TimeOut: func() int {
			timeout, _ := strconv.Atoi(os.Getenv("GOOGLE_API_TIMEOUT"))
			return timeout
		}(),
	}

	gah.NetClient = &http.Client{
		Transport: gah.Transport,
		Timeout:   time.Duration(getEnvInt("GOOGLE_API_TIMEOUT")) * time.Second,
	}
	var err error
	//gah.GoogleApiClient, err = language.NewClient(context.Background(), option.WithCredentialsJSON(GoogleCredentialsJSON), option.WithHTTPClient(gah.NetClient))

	gah.GoogleApiClient, err = language.NewClient(context.Background(), option.WithCredentialsJSON(GoogleCredentialsJSON()))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return gah
}

func (g *GoogleNLPApiHandler) AnalyzeSentiment(text string) (string, float32){
	ctx := context.Background()

	// Detects the sentiment of the text.
	sentiment, err := g.GoogleApiClient.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})

	if err != nil {
		log.Fatalf("Failed to analyze text: %v", err)
	}

	fmt.Printf("Text: %v\n", text)
	var label string
	if sentiment.DocumentSentiment.Score == 0{
		label = "neutral"
	}else if sentiment.DocumentSentiment.Score > 0 {
		label = "positive"
	} else {
		label = "negative"
	}
	return label, sentiment.DocumentSentiment.Score
}

func GoogleCredentialsJSON() []byte{
	//logrus.Info(os.Getenv("GOOGLE_API_TIMEOUT"))
	cred, err := ioutil.ReadFile(os.Getenv("GOOGLE_CREDENTIALS_JSON"))
	if err != nil {
		panic(err)
	}
	return cred
}
