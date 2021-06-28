FROM golang:1.16-alpine

RUN apk add --no-cache git bash make
ENV PS1='\n[\u@\h \w]\n$ '
RUN go env -w GOPROXY=https://goproxy.io

WORKDIR /go/src/
RUN git clone https://github.com/alacine/video_server.git

WORKDIR /go/src/video_server
RUN cd api/ && go mod download
RUN cd scheduler/ && go mod download
RUN cd streamserver/ && go mod download
RUN make && make install
