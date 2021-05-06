package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"time"
	"sync"
)

const ShellToUse = "bash"
var lock sync.Mutex

func Shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func Handler(conn net.Conn) {
	buf := make([]byte, 2048)
	//读取客户端发送的内容
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fileName := string(buf[:n])
	//获取客户端ip+port
	addr := conn.RemoteAddr().String()
	fmt.Println(addr + ": 客户端传输的文件名为--" + fileName)
	recvStart_Time := time.Now()
	fmt.Print("开始接收数据： ")
	fmt.Println(recvStart_Time)
	//告诉客户端已经接收到文件名
	conn.Write([]byte("ok"))
	//创建文件
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	defer f.Close()
	//循环接收客户端传递的文件内容
	for {
		buf := make([]byte, 2048)
		n, _ := conn.Read(buf)
		//结束协程
		if n == 0 {
			recvFinish_Time := time.Now()
			fmt.Print("接收数据完成： ")
			fmt.Println(recvFinish_Time)
			lock.Lock()  //加锁
			ExecErr, _, _ := Shellout("./darknet detector test cfg/coco.data cfg/yolov4.cfg yolov4.weights -dont_show -ext_output " + fileName)
			if ExecErr != nil {
				fmt.Printf("error: %v\n", ExecErr)
			}
			lock.Unlock()//解锁
			computeFInish_Time := time.Now()
			//fmt.Println(addr + ": 计算结束")
			fmt.Print("计算结束时间： ")
			fmt.Println(computeFInish_Time)
			ExecErr, _, _ = Shellout("mv result.txt " + "result/" + fileName + ".txt")
			if ExecErr != nil {
				fmt.Printf("error: %v\n", ExecErr)
			}
			//fmt.Println(addr + ": 协程结束")
			runtime.Goexit()
		}
		f.Write(buf[:n])
	}
}

func main() {
   //runtime.GOMAXPROCS(1)//只有一个物理核心，并行也串行
	//创建tcp监听
	listen, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listen.Close()
	fmt.Println("应用启动成功，请开始发送数据！")
	for {
		//阻塞等待客户端
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		//创建协程
		go Handler(conn)
	}
}
