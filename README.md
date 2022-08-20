# HTTP 代理

- 支持 HTTP，HTTPS，HTTP2，WS，WSS
- 配置文件：`config.json`

## 配置示例

```json
{
  "Service": [
    {
      "Target": "http://localhost:8080/api2",
      "Path": "/api1"
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

- `Service` 代理服务列表
  - `Target` 目标 URL，包括协议，主机，端口，以及路径的前缀
  - `Path` 源路径前缀
  - 路径不能以 `/` 结尾
  - 在本例中，代理服务器会将 `https://0.0.0.0/api1/ping` 转发到 `http://localhost:8080/api2/ping`
- `DynamicService` 动态代理服务
  - `Path` 路径
  - `Query` 请求参数中，目标 url 的参数名
  - 在本例中，代理服务器会将 `https://0.0.0.0/proxy?url=http%3A%2F%2Flocalhost%3A8080%2Fapi` 转发到 `http://localhost:8080/api`
- `HTTPS` 证书配置项
  - `CertFile` 证书
  - `KeyFile` 私钥
