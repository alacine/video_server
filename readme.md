## Go 视频网站项目

### 运行启动说明

`initdb.sql` 是原来的开发时用到的脚本，创建数据库还是使用下面的导出的比较好

运行 sql 脚本 `exportsql/video_server.sql`, 这会创建相关的数据库和代码中用到的连接用户

分别在`api`, `scheduler`, `streamserver`目录下执行`go build`

先启动`scheduler`, `streamserver`, 最后启动`api`

另外自动化部署`deployserver`暂时不可用

### API 设计

main -> middleware -> defs(message, err) -> handlers -> dbops -> response

业务量较大的部分是 dbops 部分

用户

| operation        | URL                    | Method | Status Code        |
|------------------|------------------------|--------|--------------------|
| 创建(注册)用户   | /api/users             | POST   | 201, 400, 500      |
| 获取用户基本信息 | /api/users/:uid        | GET    | 200, 400, 500      |
| 获取用户视频     | /api/users/:uid/videos | Get    | 200, 400, 500      |
| 用户登录         | /api/sessions          | POST   | 200, 400, 401, 500 |
| 用户登出         | /api/sessions          | DELETE | 205, 401, 500      |

资源(视频)

| operation    | URL              | Method | Status Code         |
|--------------|------------------|--------|---------------------|
| 获取视频列表 | /api/videos      | Get    | 200, 400, 500       |
| 获取视频信息 | /api/videos/:vid | Get    | 200, 400, 500       |
| 添加视频     | /api/videos/     | POST   | 201, 400, 401, 500  |
| 删除视频     | /api/videos/:vid | DELETE | 204, 400, 401 , 500 |

评论

| operation    | URL                       | Method | Status Code   |
|--------------|---------------------------|--------|---------------|
| 获取视频评论 | /api/videos/:vid/comments | Get    | 200, 400, 500 |
| 发布评论     | /api/videos/:vid/comments | POST   | 201, 400, 500 |

handler -> validation{1. request, 2. user} -> business logic -> response
1. data model
2. error handling


### Streaming

* 静态视频, 非 RTMP(Real-Time Messaging Protocol)
* 独立的服务, 可独立部署
* 统一的 api 格式

| operation    | URL                 | Method | Status Code   |
|--------------|---------------------|--------|---------------|
| stream video | /stream/videos/:vid | Get    | 200, 404, 500 |
| upload video | /stream/videos/:vid | POST   | 201, 400, 500 |

bucket token

bucket 中放置指定数量的 token, 当接受到 request 请求时, 为其分配一个 token,
当发送 response 后, 释放这个 token


### Scheduler

1. RESTful 的 http server
2. Timer(计时器)
3. 生产者/消费者模型下的 task runner

* api -> video_id -> mysql
* dispatcher -> mysql: video_id -> datachannel
* executor -> datachannel video_id -> delete videos

| operation    | URL                    | Method | Status Code |
|--------------|------------------------|--------|-------------|
| delete video | /scheduler/videos/:vid | DELETE | 200, 400    |
