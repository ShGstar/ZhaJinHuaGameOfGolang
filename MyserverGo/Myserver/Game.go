package Myserver

import (
	"fmt"
	"net"

	"MyserverGo/Handler"
	"MyserverGo/MyProtocol"
)

func Gameplay(bgchan chan []byte, conn *net.TCPConn) {
	//处理客户端发来的信息

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	/*
		首次连接为登陆，获取player信息，连接后只运行一次
	*/
	//没数据会阻塞
	fmt.Println("阻塞")
	bag := <-bgchan

	//解包
	fmt.Println("bag := <-bgchan")
	Opcode, Subcode, src := MyProtocol.Upack(bag) //第一次解包

	player := Handler.PlayerData{}
	player.Conn = conn
	player.ChooseOfValueOfOpcode(Opcode, Subcode, src) //选择和解包value

	for {
		//没数据会阻塞
		fmt.Println("阻塞")
		bag := <-bgchan

		//解包
		fmt.Println("bag := <-bgchan")
		Opcode, Subcode, src := MyProtocol.Upack(bag) //第一次解包

		player.ChooseOfValueOfOpcode(Opcode, Subcode, src) //选择和解包value

		fmt.Println("over Gameplay")
	}

}
