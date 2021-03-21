package main

import (
	"context"
	"fmt"
	pb "github.com/chiupc/sentiment_analytic/sentiment_analytic"
	"testing"
	"github.com/chiupc/sentiment_analytic/client_handler"
)

func Test_sentimentAnalyticServer_AnalyzeSentiment(t *testing.T) {
	ch := client_handler.NewPySentimentAnalyticGrpcClient()
	analyzedFilename, err := ch.AnalyzeSentiment(context.Background(), &pb.InputFile{FileName: "NIO_1614600519_1614686919_83c4f5aabdbb", ColumnName: "userText"})
	if err != nil {
		t.Errorf(err.Error())
	}else{
		fmt.Println(analyzedFilename)
	}
}
