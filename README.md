# HTTP 代理

- 支持 HTTP，HTTPS，HTTP2，WS，WSS
- 配置文件：`config.json`
- 支持可视化配置：http(s)://0.0.0.0:9000

## 构建

```
pnpm install
pnpm build
```

## 配置示例

```json
{
  "Password": "",
  "Service": [
    {
      "Target": "http://localhost:8080/api2",
      "Path": "/api1"
    }
  ],
  "Static": [
    {
      "Path": "/home",
      "Target": "/app/html"
    }
  ],
  "DynamicService": {
    "Path": "/proxy",
    "Query": "url"
  },
  "HTTPS": {
    "CertFile": "/key/cert.csr",
    "KeyFile": "/key/key.key"
  }
}
```

## 说明

- `Password` 配置密码
- `Service` 代理服务列表
  - `Target` 目标 URL，包括协议，主机，端口，以及路径的前缀
  - `Path` 源路径前缀
- `DynamicService` 动态代理服务
  - `Path` 路径
  - `Query` 请求参数中，目标 url 的参数名
- `Static` 静态资源服务
  - `Path` 请求路径前缀
  - `Target` 本地静态资源目录
- `HTTPS` 证书配置项
  - `CertFile` 证书
  - `KeyFile` 私钥
