package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

//发送文件到服务端
func SendFile(filePath string, //fileSize int64,
	conn net.Conn) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	for {
		buf := make([]byte, 2048)
		//读取文件内容
		n, err := f.Read(buf)
		if err != nil && io.EOF == err {
			fmt.Println("文件传输完成")
			//告诉服务端结束文件接收
			//conn.Write([]byte("finish"))
			return
		}
		//发送给服务端
		conn.Write(buf[:n])
	}
}

//image_Send "filepath" "ip address" "port"
func main() {
	var info []string //filepath + ip + port
	info = append(info, os.Args...)
	//获取文件信息
	fileInfo, err := os.Stat(info[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	ip_addr := info[2] + ":" + info[3] // ip + port
	//创建客户端连接
	conn, err := net.Dial("tcp", ip_addr) //"192.168.0.101:8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	//文件名称
	fileName := fileInfo.Name()
	//发送文件名称到服务端
	conn.Write([]byte(fileName))
	buf := make([]byte, 2048)
	//读取服务端内容
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	revData := string(buf[:n])
	if revData == "ok" {
		//发送文件数据
		SendFile(info[1], conn)
	}
}
