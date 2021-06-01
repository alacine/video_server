#!/bin/bash
# 自动化部署
set -e -o pipefail

source_code_location=~/go/src/video_server/
project_location=~/video_server/
deploy_log_file=${project_location}"deploy_log.txt"
api_log_file=${project_location}"api_log.txt"
streamserver_log_file=${project_location}"streamserver_log.txt"
scheduler_log_file=${project_location}"scheduler_log.txt"

echo "kill api server..." > $deploy_log_file
kill -9 $(pgrep api)
echo "kill streamserver..." >> $deploy_log_file
kill -9 $(pgrep streamserver)
echo "kill scheduler..." >> $deploy_log_file
kill -9 $(pgrep scheduler)

cd $source_code_location
echo "git pull..." >> $deploy_log_file
git pull
echo "copy server to project location..." >> $deploy_log_file
if [ ! -d ${project_location}"config" ]; then
    mkdir ${project_location}"config"
fi
cp ${source_code_location}"conf.json" ${project_location}"/config/conf.json"
cp ${source_code_location}"api/api" ${project_location}"api"
cp ${source_code_location}"streamserver/streamserver" ${project_location}"streamserver"
cp ${source_code_location}"scheduler/scheduler" ${project_location}"scheduler"

cd $project_location
echo "start api..." >> $deploy_log_file
nohup ./api > $api_log_file 2>&1 &
echo "start streamserver..." >> $deploy_log_file 
nohup ./streamserver > $streamserver_log_file 2>&1 &
echo "start scheduler..." >> $deploy_log_file
nohup ./scheduler > $scheduler_log_file 2>&1 &
