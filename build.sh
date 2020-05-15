#!/bin/bash
# 重新打包并重启服务，用于开发
set -e -o pipefail

go install github.com/alacine/video_server/api
go install github.com/alacine/video_server/streamserver
go install github.com/alacine/video_server/scheduler
#go install github.com/alacine/video_server/deployserver

if [ ! -d ~/go/bin/video_server ]; then
    mkdir ~/go/bin/video_server
fi

kill -9 $(pgrep api)
kill -9 $(pgrep streamserver)
kill -9 $(pgrep scheduler)

cp conf.json ~/go/bin/video_server/conf.json
cp ~/go/bin/api ~/go/bin/video_server/api
cp ~/go/bin/streamserver ~/go/bin/video_server/streamserver
cp ~/go/bin/scheduler ~/go/bin/video_server/scheduler
