## Go 视频网站项目

### 前后端解耦

优势
* 解放生产力，提高合作效率
* 松耦合的架构更灵活，部署更方便，更符合微服务的设计特征
* 性能的提升，可靠性的提升

缺点
* 工作量大，把简单的工作拆的更复杂
* 前后端分离带来的团队成本一级学习成本
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

### API 设计

main -> middleware -> defs(message, err) -> handlers -> dbops -> response

业务量较大的部分是 dbops 部分

用户

| operation        | URL              | Method | Status Code             |
|------------------|------------------|--------|-------------------------|
| 创建(注册)用户   | /user            | POST   | 201, 400, 500           |
| 用户登录         | /user/:user_name | POST   | 200, 400, 500           |
| 获取用户基本信息 | /user/:user_name | GET    | 200, 400, 401, 403, 500 |
| 用户注销         | /user/:user_name | DELETE | 204, 400, 401, 403, 500 |

资源(视频)

| operation        | URL                             | Method | Status Code             |
|------------------|---------------------------------|--------|-------------------------|
| List all videos  | /user/:user_name/videos         | Get    | 200, 400, 500           |
| Get one video    | /user/:user_name/videos/:vid-id | Get    | 200, 400, 500           |
| Delete one video | /user/:user_name/videos/:vid-id | DELETE | 204, 400, 401, 403, 500 |

评论

| operation        | URL                                  | Method | Status Code             |
|------------------|--------------------------------------|--------|-------------------------|
| show comments    | /videos/:vid-id/comments             | Get    | 200, 400, 500           |
| post a comment   | /videos/:vid-id/comments             | POST   | 201, 400, 500           |
| delete a comment | /videos/:vid-id/comments/:comment-id | DELETE | 204, 400, 401, 403, 500 |

handler -> validation{1. request, 2. user} -> business logic -> response
1. data model
2. error handling
