package client_handler

import (
	pb "github.com/chiupc/sentiment_analytic/sentiment_analytic"
	"flag"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"strings"
)

var (
	serverAddr         = flag.String("server_addr", "localhost:8444", "The server address in the format of host:port")
)

func GetTextSentiments(client pb.SentimentAnalyticClient, fileName string, columnName string){
	var text = func() []byte {
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			panic(err)
		}
		return data
	}()
	splits := strings.Split(fileName,"_")
	in := pb.InputFile{
		ColumnName: columnName,
		FileName:   splits[0],
		Text:       text,
	}
	_,err := client.AnalyzeSentiment(context.Background(),&in)
	if err != nil {
		logrus.Errorf(err.Error())
	}
}

func NewSentimentAnalyticGrpcClient() pb.SentimentAnalyticClient {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		logrus.Fatalf("fail to dial: %v", err)
	}
	//defer conn.Close()
	client := pb.NewSentimentAnalyticClient(conn)
	return client
}