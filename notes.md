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
