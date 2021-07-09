FROM alpine:latest

ENV APP_DIR=/app
RUN apk add --no-cache make go iproute2
RUN addgroup -S appgroup --gid 1000 \
    && adduser -DS appuser -G appgroup -h ${APP_DIR} --uid 1000

USER appuser
ENV GOROOT /usr/lib/go
ENV GOPATH /go

ENV PS1='\n[\u@\h \w]\n$ '
