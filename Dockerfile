FROM golang:1.12-alpine as builder
RUN apk add git
COPY . /go/src/shuTodo
ENV GO111MODULE on
WORKDIR /go/src/shuTodo
RUN go get && go build

FROM alpine
MAINTAINER longfangsong@icloud.com
COPY --from=builder /go/src/shuTodo/shuTodo /
WORKDIR /
CMD ./shuTodo
ENV PORT 8000
EXPOSE 8000