## Go 视频网站项目

### API 设计

main -> middleware -> defs(message, err) -> handlers -> dbops -> response

业务量较大的部分是 dbops 部分

用户

| operation        | URL                  | Method | Status Code             |
|------------------|----------------------|--------|-------------------------|
| 创建(注册)用户   | /api/register        | POST   | 201, 400, 500           |
| 用户登录         | /api/login           | POST   | 200, 400, 500           |
| 获取用户基本信息 | /api/user/:user_name | GET    | 200, 400, 401, 403, 500 |
| 用户注销         | /api/user/:user_name | DELETE | 204, 400, 401, 403, 500 |

资源(视频)

| operation        | URL                         | Method | Status Code             |
|------------------|-----------------------------|--------|-------------------------|
| List user videos | /api/user/:user_name/videos | Get    | 200, 400, 500           |
| Get one video    | /api/videos/:vid            | Get    | 200, 400, 500           |
| Delete one video | /api/videos/:vid            | DELETE | 204, 400, 401, 403, 500 |

评论

| operation        | URL                                   | Method | Status Code             |
|------------------|---------------------------------------|--------|-------------------------|
| show comments    | /api/videos/:vid/comments             | Get    | 200, 400, 500           |
| post a comment   | /api/videos/:vid/comments             | POST   | 201, 400, 500           |

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
| upload video | /upload/videos/:vid | POST   | 200, 400, 500 |

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
