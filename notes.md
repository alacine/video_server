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

| operation    | URL                 | Method | Status Code   |
|--------------|---------------------|--------|---------------|
| stream video | /stream/videos/:vid | Get    | 200, 404, 500 |
| upload video | /stream/videos/:vid | POST   | 201, 400, 500 |

bucket token

bucket 中放置指定数量的 token, 当接受到 request 请求时, 为其分配一个 token,
当发送 response 后, 释放这个 token

### SCHEDULER

定时任务，目前实际只有一个删除视频的任务

1. RESTful 的 http server
2. Timer(计时器)
3. 生产者/消费者模型下的 task runner

* api -> video_id -> mysql
* dispatcher -> mysql: video_id -> datachannel
* executor -> datachannel video_id -> delete videos

| operation    | URL                    | Method | Status Code |
|--------------|------------------------|--------|-------------|
| delete video | /scheduler/videos/:vid | DELETE | 200, 400    |

### 前后端解耦

优势
* 解放生产力，提高合作效率
* 松耦合的架构更灵活，部署更方便，更符合微服务的设计特征
* 性能的提升，可靠性的提升

缺点
* 工作量大，把简单的工作拆的更复杂
* 前后端分离带来的团队成本以及学习成本
* 系统复杂度加大


### RESTful API

* REST(Representational Status Transfer) API
* REST 是一种设计风格，不是任何架构标准
* 当今 RESTful API 通常使用 HTTP 作为通信协议，JSON 作为数据格式

特点
* 统一接口(Uniform Inrterface)
* 无状态(Stateless)
* 可缓存(Cacheable)
* 分层(Layered System)
* CS 模式(Client-server Architecture)

设计原则
* 以 URL(统一资源定位符号)风格设计 API
* 通过不同的 METHOD(GET, POST, PUT, DELETE)来区分资源的 CRUD
* 返回码(Status Code)符合 HTTP 资源描述的规定

### MYSQL

mysqldump --single-transaction --add-drop-database -h 127.0.0.1 -u video --all-databases -p > export.sql
