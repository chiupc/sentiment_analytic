package client_handler

import (
	"fmt"
	pb "github.com/chiupc/sentiment_analytic/sentiment_analytic"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var logger = logrus.New()

func init(){
	logger.SetFormatter(&logrus.JSONFormatter{})
	switch logLevel := strings.ToLower(os.Getenv("LOGLEVEL")); logLevel {
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "trace":
		logger.SetLevel(logrus.TraceLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
}

func GetTextSentiments(ctx context.Context, fileName string, columnName string, analyzerEngine string) error{
	log := logger.WithFields(logrus.Fields{
		"ctx-id": ctx.Value("ctx-id"),
		"func": "GetTextSentiments",
	})
	splits := strings.Split(fileName, ".")
	if analyzerEngine == "GCP-NLP" {
		client := NewSentimentAnalyticGrpcClient()
		text, err := ioutil.ReadFile(filepath.Join(os.Getenv("DATA_PATH"), fileName+".csv"))
		if err != nil {
			return err
		}

		in := pb.InputFile{
			ColumnName: columnName,
			FileName:   splits[0],
			Text:       text,
		}
		_, err = client.AnalyzeSentiment(context.Background(), &in)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		return nil
	}else{
		client := NewPySentimentAnalyticGrpcClient()
		_, err := client.AnalyzeSentiment(context.Background(), &pb.InputFile{FileName: splits[0], ColumnName: "userText"})
		if err != nil {
			log.Error(err.Error())
		}
		return err
	}
}

func NewSentimentAnalyticGrpcClient() pb.SentimentAnalyticClient {
	//flag.Parse()
	var serverAddr string
	if os.Getenv("ENVIRONMENT") == "DEV"{
		serverAddr = fmt.Sprintf("%s:%d","localhost",os.Getenv("GRPC_SENTIMENTANALYTIC_PORT"))
	}else if os.Getenv("ENVIRONMENT") == "TEST"{
		serverAddr = fmt.Sprintf("%s:%d",os.Getenv("GRPC_SENTIMENTANALYTIC_TEST_VIP"),os.Getenv("GRPC_SENTIMENTANALYTIC_PORT"))
	}else if os.Getenv("ENVIRONMENT") == "PROD"{
		serverAddr = fmt.Sprintf("%s:%d",os.Getenv("GRPC_SENTIMENTANALYTIC_PROD_VIP"),os.Getenv("GRPC_SENTIMENTANALYTIC_PORT"))
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		logger.Errorf("fail to dial: %v", err)
	}
	//defer conn.Close()
	client := pb.NewSentimentAnalyticClient(conn)
	return client
}

func NewPySentimentAnalyticGrpcClient() pb.SentimentAnalyticClient {
	//flag.Parse()
	log := logger.WithFields(logrus.Fields{
		"func": "NewPySentimentAnalyticGrpcClient",
	})
	var serverAddr string
	if os.Getenv("ENVIRONMENT") == "DEV"{
		log.Infof("Environment = %s", os.Getenv("ENVIRONMENT"))
		serverAddr = fmt.Sprintf("%s:%s","localhost",os.Getenv("PY_GRPC_SENTIMENTANALYTIC_PORT"))
	}else if os.Getenv("ENVIRONMENT") == "TEST"{
		serverAddr = fmt.Sprintf("%s:%s",os.Getenv("PY_GRPC_SENTIMENTANALYTIC_TEST_VIP"),os.Getenv("PY_GRPC_SENTIMENTANALYTIC_PORT"))
	}else if os.Getenv("ENVIRONMENT") == "PROD"{
		serverAddr = fmt.Sprintf("%s:%d",os.Getenv("PY_GRPC_SENTIMENTANALYTIC_PROD_VIP"),os.Getenv("PY_GRPC_SENTIMENTANALYTIC_PORT"))
	}
	log.Infof("Grpc server = %s", serverAddr)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		logger.Errorf("fail to dial: %v", err)
	}
	//defer conn.Close()
	client := pb.NewSentimentAnalyticClient(conn)
	return client
}