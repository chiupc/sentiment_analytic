package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"github.com/sharkdetector/sentiment_analytic/client_handler"
)

//func clean(s []byte) string {
//	j := 0
//	for _, b := range s {
//		if ('a' <= b && b <= 'z') ||
//			('A' <= b && b <= 'Z') ||
//			('0' <= b && b <= '9') ||
//			b == ' ' {
//			s[j] = b
//			j++
//		} else {
//			fmt.Print(b)
//		}
//	}
//	return string(s[:j])
//}

func TestPeterSO(b *testing.T) {
	//for N := 0; N < b.N; N++ {
	//b.StopTimer()
	//b.StartTimer()
	cleaned := clean(strShakespeare)
	ioutil.WriteFile(`NIO_cleaned.csv`, cleaned, os.FileMode(666))
	//}
}

var strShakespeare = func() []byte {
	// The Complete Works of William Shakespeare by William Shakespeare
	// http://www.gutenberg.org/files/100/100-0.txt
	data, err := ioutil.ReadFile(`NIO_raw.csv`)
	if err != nil {
		panic(err)
	}
	return data
}()

var strShakespeareCleaned = func() []byte {
	// The Complete Works of William Shakespeare by William Shakespeare
	// http://www.gutenberg.org/files/100/100-0.txt
	data, err := ioutil.ReadFile(`NIO_cleaned.csv`)
	if err != nil {
		panic(err)
	}
	return data
}()

func TestAnalyzeSentiment(t *testing.T) {
	g := NewGoogleNLPApiHandler()
	text := clean(strShakespeare)
	//ioutil.WriteFile(`NIO_cleaned.csv`, text, os.FileMode(666))
	br := bytes.NewReader(text)
	//f, _ := os.Open(`NIO_cleaned.csv`)
	c := csv.NewReader(br)
	lines, err := c.ReadAll()
	fmt.Println(c.ReadAll())
	var fieldId int
	for i, header := range lines[0]{
		if header == "userText"{
			fieldId = i
		}
	}
	if err != nil{
		t.Error(err.Error())
	}
	f, err := os.OpenFile(filepath.Join("NIO_processed.csv"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//write header
	lines[0] = append(lines[0],"score")
	lines[0] = append(lines[0],"sentiment")
	f.WriteString(strings.Join(lines[0],",") + "\n")
	for i, line := range lines[1:]{
		//fmt.Println(line[fieldId])
		sentiment, score := g.AnalyzeSentiment(line[fieldId])
		lines[i+1][fieldId] = "\"" + lines[i+1][fieldId] + "\""
		lines[i+1] = append(lines[i+1],fmt.Sprintf("%f", score))
		lines[i+1] = append(lines[i+1],sentiment)
		f.WriteString(strings.Join(lines[i+1],",") + "\n")
	}
}

func TestGetTextSentiments(t *testing.T){
	client := client_handler.NewSentimentAnalyticGrpcClient()
	client_handler.GetTextSentiments(client, "NIO_raw.csv","userText")
}

func TestFileCheckSum(t *testing.T){
	text := strShakespeare
	md5 := md5.Sum(text)
	sha1 := sha1.Sum(text)
	sha256 := sha256.Sum256(text)

	fmt.Printf("%x\n", md5)
	fmt.Printf("%x\n", sha1)
	fmt.Printf("%x\n", sha256)

}
