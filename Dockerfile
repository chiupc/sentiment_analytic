FROM golang:1.14.7

ADD . /SentimentAnalyticGo
WORKDIR /SentimentAnalyticGo
RUN go get -d -v ./...
RUN go install -v ./...
CMD ["go", "run", "server.go"]