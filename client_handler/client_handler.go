package client_handler

import (
	"flag"
	"fmt"
	pb "github.com/chiupc/sentiment_analytic/sentiment_analytic"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	serverAddr         = flag.String("server_addr", "localhost:8444", "The server address in the format of host:port")
)

func GetTextSentiments(client pb.SentimentAnalyticClient, fileName string, columnName string) error{
	fmt.Println("Running GetTextSentiments")
	fmt.Println(filepath.Join(os.Getenv("DATA_PATH"),fileName + ".csv"))
	text, err := ioutil.ReadFile(filepath.Join(os.Getenv("DATA_PATH"),fileName + ".csv"))
	if err != nil{
		return err
	}
	fmt.Println(string(text))
	splits := strings.Split(fileName,".")
	fmt.Println(splits[0])
	in := pb.InputFile{
		ColumnName: columnName,
		FileName:   splits[0],
		Text:       text,
	}
	_, err = client.AnalyzeSentiment(context.Background(),&in)
	if err != nil {
		fmt.Errorf(err.Error())
		return err
	}
	return nil
}

func PrintSomething(){
	fmt.Println("test")
}

func NewSentimentAnalyticGrpcClient() pb.SentimentAnalyticClient {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	//defer conn.Close()
	client := pb.NewSentimentAnalyticClient(conn)
	return client
}