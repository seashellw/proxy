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
  "CDNList": [
    "https://cdn-1259243245.cos.ap-shanghai.myqcloud.com",
    "https://esm.sh",
    "https://cdn.jsdelivr.net"
  ],
  "HTTPS": {
    "CertFile": "/key/cert.csr",
    "KeyFile": "/key/key.key"
  }
}
```
