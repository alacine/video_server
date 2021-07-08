FROM alpine:latest

ENV APP_DIR=/app
RUN apk add --no-cache make go iproute2
RUN addgroup -S appgroup && adduser -DS appuser -G appgroup -h ${APP_DIR}

USER appuser
ENV GOROOT /usr/lib/go
ENV GOPATH /go
RUN mkdir -pv /app/config

ENV PS1='\n[\u@\h \w]\n$ '

WORKDIR $APP_DIR
COPY --from=video_server_build ${GOPATH}/bin/api .
COPY --from=video_server_build ${GOPATH}/bin/conf.json config/
CMD ["./api"]
