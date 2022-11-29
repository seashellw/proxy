# PROXY

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
