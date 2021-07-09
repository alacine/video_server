FROM golang:1.16-alpine

RUN apk add --no-cache git bash make
ENV PS1='\n[\u@\h \w]\n$ '
RUN go env -w GOPROXY=https://goproxy.io

ENV PROJECT=/go/src/video_server
WORKDIR /go/src/
RUN mkdir -pv $PROJECT/api $PROJECT/scheduler $PROJECT/streamserver
COPY ./api $PROJECT/api
COPY ./scheduler $PROJECT/scheduler
COPY ./streamserver $PROJECT/streamserver
COPY ./Makefile $PROJECT
#RUN git clone https://github.com/alacine/video_server.git

WORKDIR $PROJECT
RUN cd api/ && go mod download
RUN cd scheduler/ && go mod download
RUN cd streamserver/ && go mod download
RUN make && make install
