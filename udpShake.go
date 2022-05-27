package main

import (
	"net"
)

type downFile struct {
	conn    *net.UDPConn
	name    string
	RecvLen int
	updAddr *net.UDPAddr
}

//封装握手包
func Shake(megType byte, path []byte, fileSize int) []byte {
	haskData := make([]byte, 0)

	haskData = append(haskData, megType)

	bodys := IntToBytes(len(path))
	//文件名称长度
	haskData = append(haskData, bodys...)

	//文件名称
	haskData = append(haskData, path...)

	FileLen := IntToBytes(fileSize)

	//文件长度
	haskData = append(haskData, FileLen...)

	return haskData
}
func UnShake(shakeBytes []byte) downFile {
	readlen := 0
	var down downFile = downFile{}
	//fmt.Println(shakeBytes)
	if shakeBytes[readlen] == 0x00 {
		readlen++
		NameLen := BytesToInt(shakeBytes[readlen : readlen+8])

		readlen += 8

		FileName := string(shakeBytes[readlen : readlen+NameLen])

		down.name = FileName

		readlen += NameLen

		FileSize := BytesToInt(shakeBytes[readlen:len(shakeBytes)])

		down.RecvLen = FileSize

	}
	return down
}
