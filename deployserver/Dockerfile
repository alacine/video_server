FROM video_server_base

USER appuser
WORKDIR $APP_DIR

COPY --from=video_server_build --chown=appuser:appgroup ${GOPATH}/bin/deployserver .

CMD ["./deployserver"]
