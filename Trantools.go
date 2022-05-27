package main

import (
	"flag"
	"strings"
)

type Trantools struct {
}

func (t Trantools) Start() {

	TarIP := flag.String("host", "0.0.0.0:6001", "输入目标地址")
	sendFile := flag.String("path", "file/msqls.dmg", "请输入发送文件")
	MaxLen := flag.Int("sl", 512, "请输入每次发送长度")
	tar := flag.String("tar", "server", "你是客户端还是服务端")
	flag.Parse()

	if strings.EqualFold(*tar, "server") {
		n := Netupd{}

		go n.startService()
	}

	cli := TranClient{}
	cli.ConnServer(*TarIP, *MaxLen)

	cli.SendFile(*sendFile)

}
