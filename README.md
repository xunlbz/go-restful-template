# Edge Admin
- go-resultful + vue + gorm + gobindata打包（页面+后台）+ swaggerui
- 访问地址: 
```
http://localhost:8889
```

## API
- 用户管理


## 鉴权
- Bearer Token 
- JWT Token Via HTTP Header (Authorization:${token})

## 开发

- 目录说明
  - bindata  静态文件打包生成代码目录
  - cmd      命令行执行预留
  - docs     接口文档预留
  - pkg      后台代码
  - webui    前端代码
  - systemd  systemd脚本

- 开启swagger ui, 此模式为API开发模式将不加载VUE页面
  - ./edge_admin -swagger

## 编译打包

```shell
make build
```
## 静态页面打包成go文件
```shell
go generate
```

## WebSocket
- 路径
```
http://localhost:8889/ws
```
- 消息格式
```json
{"type":"collect","value":"os","interval":1}
{"type":"service_log","value":"edgeAdmin"}
{"type":"container_log","value":"容器id"}
```
 
- Type 分类
  - collect  #参数interval 获取频率，再收集使用率等指标情况下可生效
  - service_log
  - container_log

- 使用参考
[websocket.js](https://github.com/xunlbz/go-restful-template/tree/master/pkg/websocket/websocket.js)