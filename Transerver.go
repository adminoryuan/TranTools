package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

const ()

//接收服务
type Netupd struct {
}

var draProcess DrawProcess

func (n Netupd) startService() {
	listend, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 6001})

	if err != nil {
		panic(err)
	}

	defer listend.Close()
	fmt.Println("服务已经启动", 6001)

	for {

		bf := make([]byte, 1024)

		lens, addr, err := listend.ReadFromUDP(bf[:])

		if err != nil {
			panic(err)

		}
		fmt.Println("收到握手包", addr.IP)

		shakeBytes := bf[:lens]

		df := UnShake(shakeBytes)
		df.updAddr = addr
		df.conn = listend

		draProcess = Init(int64(df.RecvLen))

		//df.conn.Write([]byte{0x01}) //表示客户端可以发送了
		n.DownLoadFile(df)

		fmt.Println("文件接收完毕")
		break

	}
}

//开始下载文件
func (n Netupd) DownLoadFile(down downFile) {
	//go draProcess.Draw()
	BytePool := sync.Pool{
		New: func() interface{} {
			return make([]byte, 1500)
		},
	}

	recvLens := 0
	file, err := os.Create(down.name)

	if err != nil {
		panic(err)
	}

	for recvLens < down.RecvLen {
		bytes := BytePool.Get().([]byte) //最大mtu

		down.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		lens, err := down.conn.Read(bytes)
		if err != nil {

			ackBodys := make([]byte, 0)

			ackBodys = append(ackBodys, IntToBytes(recvLens)...)

			ackBodys = append(ackBodys, IntToBytes(down.RecvLen)...) //丢包长度

			down.conn.WriteToUDP(ackBodys, down.updAddr)

			continue

		} else if err == io.EOF {
			fmt.Println("链接已经断开")
		}

		recvBody := bytes[:lens] //

		//1 10  //
		SerialId := BytesToInt(recvBody[:8]) //包id

		Sendlen := BytesToInt(recvBody[8:16]) //数据长度

		//判断是否丢包
		if SerialId > recvLens+Sendlen {

			RertLen := SerialId - Sendlen
			for i := recvLens; i < RertLen; {
				ackBodys := make([]byte, 0)

				ackBodys = append(ackBodys, IntToBytes(i)...)

				ackBodys = append(ackBodys, IntToBytes(1024)...) //丢包长度

				go down.conn.WriteToUDP(ackBodys, down.updAddr)

				recvBody1 := make([]byte, 1500)

				down.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
				lens1, err := down.conn.Read(recvBody1)

				if err == nil {

					temLen := BytesToInt(recvBody1[8:16]) //数据长度

					recvLens += temLen

					i += temLen

					file.Write(recvBody1[16:lens1])

				}

			}

		}
		//	draProcess.ProSignle <- int64(recvLens)
		BytePool.Put(bytes)

		file.Write(recvBody[16:])

		recvLens += Sendlen

	}

	//down.conn.WriteToUDP(ackBodys, down.updAddr)
	fmt.Println("接收完毕")

}
