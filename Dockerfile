FROM golang:1.10.2-alpine3.7

RUN apk update ; apk add git

WORKDIR /go/src/urls_checker

RUN go get -u github.com/mailru/easyjson/...
RUN go get github.com/patrickmn/go-cache

COPY . .

RUN go install ./internal/app/urls_checker
RUN go build -o ./myapp ./cmd/urls_checker

CMD "./myapp"


