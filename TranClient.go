package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type TranClient struct {
	con     net.Conn
	SendLen int //定义一次发送数据大小

}

func (t *TranClient) SendFile(path string) {

	File, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("找不到文件")
	}
	pathName := strings.Split(path, "/")

	if err != nil {
		panic(err)
	}
	//发送握手包
	_, err = t.con.Write(Shake(0x00, []byte(pathName[len(pathName)-1]), len(File)))
	if err != nil {
		fmt.Println("发送握手包失败")
	}

	ack := make([]byte, 1)
	//t.con.Read(ack)
	if ack[0] == 0x00 {
		fmt.Println("双发握手成功 开始发送文件")
		Robj := NewReialbSend(t.con, File, t.SendLen)
		Robj.Handle()
	}
}

func (t *TranClient) ConnServer(target string, len int) {
	t.SendLen = 512
	con, err := net.Dial("udp", target)

	if err != nil {
		fmt.Println("错误")
	}
	t.SendLen = len

	t.con = con
}
