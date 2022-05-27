package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

//实现可靠传输发送文件
type reilabSend struct {
	conn       net.Conn
	close      chan struct{}
	wait       sync.Locker
	file       []byte
	maxFileLen int
	since      int
}

func NewReialbSend(conn net.Conn, f []byte, Since int) reilabSend {

	return reilabSend{
		conn:  conn,
		file:  f,
		close: make(chan struct{}),
		since: Since,
		wait:  &sync.Mutex{},
	}
}

func (r *reilabSend) Handle() {
	r.maxFileLen = len(r.file)
	go r.send(0)
	go r.Read()
}

//检测是否丢包
func (r *reilabSend) Read() {
	for {

		byt := make([]byte, 32)
		r.conn.SetReadDeadline(time.Now().Add(time.Second * 2))
		_, err := r.conn.Read(byt)
		if err != nil {
			continue
		}

		serid := BytesToInt(byt[:8]) //获取需要重传的序列号
		sendLen := BytesToInt(byt[8:16])

		if sendLen == r.maxFileLen {
			//	fmt.Println("需要全部重发")
			go r.send(serid)
			continue
			//从这里重新发送
		}
		r.wait.Lock() //重传时 短暂

		newBody, _ := r.NewBodys(serid, sendLen)

		r.conn.Write(newBody)

		time.Sleep(10 * time.Microsecond)
		r.wait.Unlock()

	}

}
func (r *reilabSend) NewBodys(alredyLen int, sendLen int) ([]byte, int) {
	newBodys := make([]byte, 0)

	//头8个字节表示当前已经的发送的长度

	if alredyLen+r.since < r.maxFileLen {

		newBodys = append(newBodys, IntToBytes(alredyLen+sendLen)...) //发送的最后长度
		newBodys = append(newBodys, IntToBytes(sendLen)...)
		newBodys = append(newBodys, r.file[alredyLen:(alredyLen+sendLen)]...)

		alredyLen += r.since

	} else {
		newBodys = append(newBodys, IntToBytes(r.maxFileLen)...) //发送的最后长度
		newBodys = append(newBodys, IntToBytes(r.maxFileLen-alredyLen)...)
		newBodys = append(newBodys, r.file[alredyLen:r.maxFileLen]...)

		alredyLen = r.maxFileLen
	}
	return newBodys, alredyLen
}

func (r *reilabSend) send(alredy int) {
	maxlen := len(r.file) //需要发送的最大长度

	alredyLen := alredy // 保存已经发送的长度

	for alredyLen < maxlen {

		bodys, alerLen := r.NewBodys(alredyLen, r.since)

		r.conn.Write(bodys)

		alredyLen = alerLen

		r.wait.Lock()

		r.wait.Unlock()

		time.Sleep(1 * time.Microsecond)
		//一直发不用管丢包

	}
	fmt.Println("发送完毕")

}
