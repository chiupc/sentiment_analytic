package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	pb "github.com/chiupc/sentiment_analytic/sentiment_analytic"
	"github.com/chiupc/sentiment_analytic/client_handler"
)

var (
	port       = flag.Int("port", 8444, "The server port")
)

type sentimentAnalyticServer struct {
	pb.UnimplementedSentimentAnalyticServer
}

func (s *sentimentAnalyticServer) AnalyzeSentiment(ctx context.Context, file *pb.InputFile) (*pb.OutputFile, error) {
	//TODO - Process before
	fmt.Println("AnalyzeSentiment is triggered")
	if os.Getenv("SENTIMENT_ENGINE") == "GCP-NLP" {
		g := NewGoogleNLPApiHandler()
		text := clean(file.Text)
		ioutil.WriteFile(filepath.Join(filepath.Join(os.Getenv("DATA_PATH"), file.FileName+"_processed.csv")), text, os.FileMode(666))
		br := bytes.NewReader(text)
		c := csv.NewReader(br)
		lines, err := c.ReadAll()
		var fieldId int
		for i, header := range lines[0] {
			if header == file.ColumnName {
				fieldId = i
			}
		}
		if err != nil {
			return nil, err
		}
		f, err := os.OpenFile(filepath.Join(os.Getenv("DATA_PATH"), file.FileName+"_processed.csv"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		//write header
		lines[0] = append(lines[0], "score")
		lines[0] = append(lines[0], "sentiment")
		f.WriteString(strings.Join(lines[0], ",") + "\n")
		for i, line := range lines[1:] {
			//fmt.Println(line[fieldId])
			sentiment, score := g.AnalyzeSentiment(line[fieldId])
			lines[i+1][fieldId] = "\"" + lines[i+1][fieldId] + "\""
			lines[i+1] = append(lines[i+1], fmt.Sprintf("%f", score))
			lines[i+1] = append(lines[i+1], sentiment)
			f.WriteString(strings.Join(lines[i+1], ",") + "\n")
		}

		return &pb.OutputFile{Filename: file.FileName+"_processed.csv"}, nil
	}else{ //else if os.Getenv("SENTIMENT_ENGINE") == "VADER"
		pyGrpc := client_handler.NewPySentimentAnalyticGrpcClient()
		out, err := pyGrpc.AnalyzeSentiment(context.Background(),&pb.InputFile{
			ColumnName: file.ColumnName,
			FileName:   file.FileName,
			Text:       nil,
		})
		if err != nil {
			logrus.Errorf(err.Error())

		}
		return out, err
	}
}

func newServer() *sentimentAnalyticServer {
	s := &sentimentAnalyticServer{}
	return s
}

func main() {
	godotenv.Load(".env")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		logrus.Errorf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSentimentAnalyticServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}