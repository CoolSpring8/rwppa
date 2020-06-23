# rwppa

RVPN Web Portal Proxy Adapter (based on MITM)

因为不想在电脑上安装 EasyConnect，所以做了一个访问 ZJU 校内网站的轻量级方案——将 [浙江大学 RVPN 网页版](https://rvpn.zju.edu.cn) 模拟为本地 HTTP 代理。

## 使用

1. 前往 [Releases](https://github.com/CoolSpring8/rwppa/releases/) 下载并放于合适位置；
2. 运行，输入上网服务账号、密码（即用于校内拨号 VPN 的账号）以及监听端口（如“127.0.0.1:8080”代表监听 8080 端口，仅接受本机流量）。账号密码仅会被发送至 rvpn.zju.edu.cn，用于登录认证。
3. 点击登录后别着急，先将程序同目录下生成的 rootCA.crt 安装至“受信任的根证书颁发机构”存储区。（之后运行不需要这一步）
4. 用 SwitchyOmega 等设置好代理，start browsing！

注意不要删除 rootCA.crt 和 rootCA.key，否则需要重新添加信任。

## 特性及局限

### 特性：

- 配置完毕并登录后，可访问校内校外 HTTP / HTTPS 网站（例如 CC98，正版软件服务与管理平台等）。
- 支持 HTTP 的 GET，POST，PUT，DELETE，HEAD 方法，一定程度上支持 OPTIONS 方法（仅为模拟）。
- 正确处理大部分网页及网页资源。如果碰到问题，欢迎提 Issue/PR ！

### 局限：

- **不支持 WebSocket，FTP，SSH等非 HTTP 协议。**你可能会想了解 [Hagb/docker-easyconnect](https://github.com/Hagb/docker-easyconnect) 项目。
- 不支持 HTTP 的 PATCH 等方法。
- OPTIONS 方法由于 RVPN 网页版不支持，所以实际不会向服务器发送请求，而是由 rwppa 直接返回一个针对 CORS 预检而配置的较为宽松的结果，安全性有所降低。
- 不支持 TLS 1.3 Early Data（0-RTT）。若原网站启用过 TLS 1.3 Early Data 且 session ticket 还没过期则会[无法访问](https://golang.org/src/crypto/tls/handshake_server_tls13.go)（如 V2EX）。（但是短期内大概不会有开启这个特性的校内网站吧。）

## 其他

提出这个想法的最初讨论：[有没有将 CGI Proxy 转化为普通 HTTP Proxy 的工具呢？ - V2EX](https://www.v2ex.com/t/670356)

欢迎提出 Issue 或 PR！

## LICENSE

GPLv3