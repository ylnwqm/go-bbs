FROM golang:1.14.2-alpine3.11

ENV GOPROXY=https://goproxy.cn

WORKDIR /work

COPY . .

RUN go build -o go-bbs main.go

CMD ["./start.sh"]