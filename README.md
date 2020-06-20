# rwppa

RVPN Web Portal Proxy Adapter

尝试将 RVPN 网页版接口模拟成一个 HTTP 代理。

[有没有将 CGI Proxy 转化为普通 HTTP Proxy 的工具呢？ - V2EX](https://www.v2ex.com/t/670356)

## 使用

1. ```bash
   openssl genrsa -out rootCA.key 4096
   ```

2. ```bash
   openssl req -x509 -new -nodes -key rootCA.key -sha256 -days 1024 -out rootCA.crt
   ```

3. 把 rootCA.crt 安装到受信任证书存储区；

4. 将 rootCA.key 和 rootCA.crt 的内容填入 cert.go；

5. ```
   go build
   ```

6. 在浏览器登录 rvpn.zju.edu.cn，从 Cookies 中找到 TWFID 字段的值；

7. ```
   .\rwppa --id TWFID字段的值
   ```

1~4 是可选的。你也可以信任并使用 cert.go 自带的证书，但出于安全性考虑不推荐这样做。

## 待修复

离实用还有不少距离，欢迎各位大佬 PR 改进！

- ~~JavaScript 文件会被改动导致网页不正常，如典型例子 jquery.min.js~~ 已解决，但不确定是否所有 js 文件都正常
- ~~HTTP 301/302 跳转的 Location 字段还没有替换~~ 已解决
- 更友好的用户交互方式（自动生成证书，自动登录获取 Cookies……）
- "tls: client sent unexpected early data"（TLS 1.3 Early Data 只能等待 Go 语言 [crypto/tls](https://golang.org/src/crypto/tls/handshake_server_tls13.go) 标准库支持）
- 疑似 CORS 丢失，见 www.cc98.org 页面对 api.cc98.org 的请求
- ……

## LICENSE

GPLv3