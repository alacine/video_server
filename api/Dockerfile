FROM video_server_base

USER appuser
WORKDIR $APP_DIR

RUN mkdir -pv $APP_DIR/config
COPY --from=video_server_build --chown=appuser:appgroup ${GOPATH}/bin/api .
COPY --from=video_server_build --chown=appuser:appgroup ${GOPATH}/bin/conf.json config/

CMD ["./api"]
