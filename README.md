# pafrwp

Proxy Adapter for RVPN Web Portal

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
   .\pafrwp --id TWFID字段的值
   ```

1~4 是可选的。你也可以信任并使用 cert.go 自带的证书，但出于安全性考虑不推荐这样做。

## 待修复

现在只达到证明概念可行的程度，离实用还有不少距离，欢迎各位大佬 PR 改进！

- JavaScript 文件会被改动导致网页不正常，如典型例子 jquery.min.js
- ~~HTTP 301/302 跳转的 Location 字段还没有替换~~
- 更友好的用户交互方式（自动生成证书，自动登录获取 Cookies……）
- "tls: client sent unexpected early data" (upstream issue?)
- ……