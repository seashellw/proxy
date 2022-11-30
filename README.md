# PROXY

## 配置示例

```json
{
  "Password": "",
  "Service": [
    {
      "Path": "/server",
      "Target": "http://server/server"
    }
  ],
  "Static": [
    {
      "Path": "/home",
      "Dir": "/root/home"
    }
  ],
  "Redirect": [
    {
      "Path": "/",
      "Target": "http://server/home"
    }
  ],
  "HTTPS": {
    "CertFile": "/key/cert.csr",
    "KeyFile": "/key/key.key"
  }
}
```
