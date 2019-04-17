package Myserver

import (
	"encoding/binary"
	"fmt"
	"net"

	"MyserverGo/Evemont"
	"MyserverGo/Handler"
	//Myserver "../Myserver"
)

var MAXBUFFERREAD int = 2048 //一次最多读取的字节数

var MAXBUFFERCHAN int = 20 //通道缓冲的大小

//建立TCP服务器
func StattServer() (bool, error) {

	addr, err := net.ResolveTCPAddr("tcp4", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	lister, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		fmt.Println(err)
		//	panic(err)
		return false, err
	}

	//此处 初始化环境资源，初始化客户队列，和房间
	Evemont.InitEvE()
	Handler.InitRoomQueue()

	fmt.Println("server success!")

	//循环接受客户端的信息
	for {
		conn, err := lister.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("接收到一个客户端连接")
		fmt.Println(conn.RemoteAddr())

		//客户端队列，及客户端对象连接池
		Evemont.ClientNum.Push(conn)

		//给每一个客户端开一个协程
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()

			//单独开一个协程处理包
			bagchan := make(chan []byte, MAXBUFFERCHAN)

			go Gameplay(bagchan, conn)

			ReadTcpData(conn, bagchan) //读取数据包
		}()

	}

	return true, nil
}

//读取数据包 tcp流
func ReadTcpData(conn *net.TCPConn, bagchan chan []byte) {

	var yiqian []byte

	//循环读取数据包
	for {
		myread := make([]byte, MAXBUFFERREAD) //缓冲区
		fmt.Println("read前")
		length, err := conn.Read(myread)
		fmt.Println("read后")

		if err != nil {
			//连接已经断开
			fmt.Println("连接已经断开")
			break
		} //for

		fmt.Println(length)
		fmt.Println(string(myread))

		//这一次读取的数据，以前剩余的数据组合验证是不是一个数据包
		myread = append(yiqian, myread[:length]...) //现在的全部数据

	A:
		if len(myread) > 1 {
			changdu := binary.BigEndian.Uint32(myread[:4])
			//数据包的长度可以获得,至少是2

			fmt.Println("长度:", changdu)
			fmt.Println("len(myread)", len(myread))

			if int(changdu) <= len(myread) {
				//数据至少有一个包

				bag := myread[4:]
				myread = myread[4:]

				bagchan <- bag
				fmt.Println("bagchan <- bag")

				myread = myread[changdu:]
				fmt.Println("获取数据包头的长度：", changdu)
				fmt.Println("len(myread): ", len(myread))
				goto A
			} else {
				//不到一个包
				yiqian = myread
			}

		} else {
			yiqian = myread
		}
	}

}
