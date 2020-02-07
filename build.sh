#!/bin/bash
set -e -o pipefail

# Build Web UI
cd ~/go/src/github.com/alacine/video_server/web
go install
if [ ! -d ~/go/bin/video_server_web_ui ]; then
    mkdir ~/go/bin/video_server_web_ui
fi
cp ~/go/bin/web ~/go/bin/video_server_web_ui/web
cp -R ~/go/src/github.com/alacine/video_server/templates_example ~/go/bin/video_server_web_ui/
