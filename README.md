# TranTools
 # 介绍
 - 基于udp 实现的一个局域网传输工具
 - 为了防止网络原因导致的丢包，所以在应用层实现了可靠传输
 - 经过测试 在丢包环境下 速率为8M/s
 # 如何使用
 - 下载代码
 ```bash
  git clone https://github.com/adminoryuan/TranTools
 ```
 -开启服务端
 - go run . tar="server
 - 客户端
 - go run . tar="cli" host="127.0.0.0:6001" -path="/file/a.dmg" -sl=512 
 
 # 参数说明
- -path 发送文件路径
- -tar 区分客户端和服务端
- -sl 每次发包大小
- -host 服务端地址
 # 布到其它平台
 - 发布到windows
 - env GOOS=windows go build .
 - 发布到mac 平台
 - env GOOS=darwin go build .
