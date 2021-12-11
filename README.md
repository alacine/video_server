## Go 视频网站项目

[![build](https://github.com/alacine/video_server/actions/workflows/build.yml/badge.svg)](https://github.com/alacine/video_server/actions/workflows/build.yml)
[![golangci-lint](https://github.com/alacine/video_server/actions/workflows/cilint.yml/badge.svg)](https://github.com/alacine/video_server/actions/workflows/cilint.yml)

在线视频网站项目的后端，前端在 [video_server_vue](https://github.com/alacine/video_server_vue)

分了 api、stream、scheduler 三个服务。项目主要用来熟悉 Go 以及学习 Makefile、
Docker、docker-compose、github actions、Jenkins 的基本使用

### 运行启动说明

**需要有 docker 和 docker-compsoe**

无论那种运行方式都是把数据库放在了 docker 里面

```bash
# 查看帮助信息
make help
```

* 本地环境

```bash
# 构建
make

# 启动数据库
make startdb
# 启动所有，用 nohup 挂在后台
make run-deamon

# 状态
make status

# 停止服务
make stop-deamon

# 停止服务且停止数据库
make stopall
```

* 本地 docker 环境

```bash
# 构建 docker 基础镜像，以及创建容器中会挂载出来的目录 local-cache
make build-in-docker

# 启动、停止
docker-compose up|down
```

* 清理操作

```bash
# 删除除了 loacl-cache 目录外所有二进制文件、nohup 产生的日志、基础镜像
make clean

# clean，并且删除挂载卷、local-cache/
make restore
```

另外自动化部署服务`deployserver`暂不需要

### API

对外业务接口服务

### STREAM

静态视频服务，在线观看、上传

### SCHEDULER

定时任务，目前实际只有一个删除视频的任务
