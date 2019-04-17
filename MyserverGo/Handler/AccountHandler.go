package Handler

import (
	"bytes"
	"fmt"

	"MyserverGo/Evemont"

	"MyserverGo/MyProtocol"
	"sync"
)

/*
	登陆请求
*/

var mu sync.Mutex

func (p *PlayerData) Login() {

	if Evemont.Player.Length() == 0 { //没有玩家

		Evemont.Player.Push(p) //添加到队列
		//打包和发送数据

		p.Sendmsg(MyProtocol.SubCode.Account, MyProtocol.AccountCode.Login_SRES, 0)
	} else {
		//有玩家，查看用户名是否重复，如果重复则让玩家重新，目前没实现，，，

		Evemont.Player.Push(p) //添加到队列
	}
	//发送数据
}

//客户端发来的更新金币数量的请求
//需要单线程执行
func (p *PlayerData) UpdateCoinCount(value int) {

	mu.Lock()
	p.Coin = p.Coin + value
	mu.Unlock()
}

/// 客户端获取排行榜的请求的处理
func (p *PlayerData) GetRankList() {

	// client.SendMsg(OpCode.Account, AccountCode.GetRankList_SRES, dto);
	//p.Sendmsg(MyProtocol.SubCode.Account, MyProtocol.AccountCode.GetRankList_SRES,
	//[]interface{}{p.State, p.Name, p.IconName, p.Coin})

}

/// 客户端获取用户信息的请求
func (p *PlayerData) GetUserInfo(name string, password string) {

	//姓名 密码 姓名 金币（为0）
	p.Sendmsg(MyProtocol.SubCode.Account, MyProtocol.AccountCode.GetUserInfo_SRES, []interface{}{p.State, p.Name, p.Name, 0})

}

func (p *PlayerData) Sendmsg(opcode uint32, subcode uint32, value interface{}) {

	//打包code
	bufcode := MyProtocol.PackCode(opcode, subcode)

	fmt.Println(".............")

	//打包value
	bufvalue := MyProtocol.PackValue(opcode, subcode, value)

	var buffer bytes.Buffer
	buffer.Write(bufcode)
	buffer.Write(bufvalue)

	//构造消息体
	packet := MyProtocol.PackMsg(buffer.Bytes())

	i, err := p.Conn.Write(packet)

	if err != nil {
		fmt.Println("SendMsg Error!")
		fmt.Println(err)
	}

	fmt.Println("发送的信息", i)
}
