# rwppa

RVPN Web Portal Proxy Adapter (based on MITM)

尝试将[浙江大学 RVPN 网页版](https://rvpn.zju.edu.cn)模拟成一个 HTTP 代理。

提出这个想法的讨论：[有没有将 CGI Proxy 转化为普通 HTTP Proxy 的工具呢？ - V2EX](https://www.v2ex.com/t/670356)

## 特性及局限

### 特性：

- 配置完毕并登录后，可访问校内校外 HTTP / HTTPS 网站（例如 CC98，正版软件服务与管理平台等），原理基于模拟 RVPN 网页版。
- 支持 HTTP 的 GET，POST，DELETE，HEAD 方法，一定程度上支持 OPTIONS 方法（模拟）。

### 局限：

- **不支持 WebSocket，FTP，SSH等非 HTTP 协议。**
- 不支持 HTTP 的 PATCH 等方法。
- OPTIONS 方法由于 RVPN 网页版不支持，所以实际不会向服务器发送请求，而是由 rwppa 直接返回一个针对 CORS 预检而配置的较为宽松的结果，安全性有所降低。
- 不支持 TLS 1.3 Early Data（0-RTT）。若原网站启用过 TLS 1.3 Early Data且 session ticket 还没过期则会[无法访问](https://golang.org/src/crypto/tls/handshake_server_tls13.go)（如 V2EX）。（但是短期内大概不会有开启这个特性的校内网站吧。）

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
   ./rwppa
   ```

1~4 是可选的。你也可以信任并使用 cert.go 自带的证书，但出于安全性考虑不推荐这样做。

由于使用了 Fyne GUI，现在编译还需要 GCC 等额外依赖。

## 其他

欢迎 PR 改进！

## LICENSE

GPLv3