package Handler

import (
	"fmt"
	"net"

	"MyserverGo/Evemont"

	"MyserverGo/MyProtocol"
)

//玩家
type PlayerData struct {
	Conn     *net.TCPConn
	Name     string
	Password string
	State    int //状态
	IconName string
	Coin     int //jinbi
}

//玩家个数
var NumberPlayer int = 0

//选择opcode的操作
func (p *PlayerData) ChooseOfValueOfOpcode(opcode uint32, subcode uint32, src []byte) {
	///操作码
	//账号模块 -0
	//模块匹配 -1
	//聊天模块 -2
	//战斗模块 -3

	switch opcode {
	case MyProtocol.SubCode.Account: //账号模块
		//return ChooseSubCodeOfAccount(subcode, src, conn)
		p.ChooseSubCodeOfAccount(subcode, src)
	case MyProtocol.SubCode.Fight: //战斗模块
		//ChooseSubCodeOfFight(subcode, src, conn)
		p.ChooseSubCodeOfFight(subcode, src)
	case MyProtocol.SubCode.Chat: //聊天模块
		//ChooseSubCodeOfChat(subcode, src, conn)
		p.ChooseSubCodeOfChat(subcode, src)
	case MyProtocol.SubCode.Math: //匹配模块
		p.ChooseSubCodeOfMatch(subcode, src)
		//ChooseSubCodeOfMatch(subcode, src, conn)
	}

	//return Evemont.PlayerData{}, false
}

//子操作码各个模块的选择判断

//账户模块选择
func (p *PlayerData) ChooseSubCodeOfAccount(subcode uint32, src []byte) {

	var name string
	var password string
	switch subcode {
	case MyProtocol.AccountCode.Login_CREQ:
		//登陆，先解包，然后进行逻辑操作
		name, password = MyProtocol.UpackSubCodeAccountStr(src)
		p.Name = name
		p.Password = password
		p.IconName = name
		p.State = NumberPlayer
		NumberPlayer++

		//UpackSubCodeAccountStr(src)
		p.Login()

		//添加到玩家队列中
		Evemont.InitAndAddPlayer(name, password, p.Conn, 0)

		fmt.Println("AccountCode.Login_CREQ")
		//return player, true
	case MyProtocol.AccountCode.Register_CREQ: //注册

	case MyProtocol.AccountCode.GetUserInfo_CREQ: //客户端获取用户信息请求
		p.GetUserInfo(name, password)
		fmt.Println("AccountCode.GetUserInfo_CREQ")
	case MyProtocol.AccountCode.GetRankList_CREQ: //客户端获取排行榜的请求处理
		p.GetRankList()

	case MyProtocol.AccountCode.UpdateCoinCount_CREQ: //客户端发来更新金币数量的请求
		value := MyProtocol.UpackIntger(src)
		p.UpdateCoinCount(value)

		//num :=UpackIntger()
	}

	//return Evemont.PlayerData{}, false
}

//战斗模块选择
func (p *PlayerData) ChooseSubCodeOfFight(subcode uint32, src []byte) {

	switch subcode {
	case MyProtocol.FightCode.Leave_CREQ: //客户端离开请求处理

	case MyProtocol.FightCode.LookCard_CREQ: //客户端看牌请求处理

	case MyProtocol.FightCode.Follow_CREQ: //客户端跟注请求处理

	case MyProtocol.FightCode.AddStakes_CREQ: //客户端加注请求处理

	case MyProtocol.FightCode.GiveUpCard_CREQ: //客户端弃牌请求处理

	case MyProtocol.FightCode.CompareCard_CREQ: //比牌请求处理
	}

}

//聊天模块
func (p *PlayerData) ChooseSubCodeOfChat(subcode uint32, src []byte) {

	switch subcode {
	case MyProtocol.ChatCode.CREQ: //聊天请求处理
		// target := UpackString(src)

	}

}

//匹配模块
func (p *PlayerData) ChooseSubCodeOfMatch(subcode uint32, src []byte) {

	switch subcode {
	case MyProtocol.MatchCode.Enter_CREQ: //进入房间请求
		value := MyProtocol.UpackIntger(src)
		p.EnterRoom(value)
		//intjiebao
	case MyProtocol.MatchCode.Leave_CREQ: //离开的请求
		value := MyProtocol.UpackIntger(src)
		p.LeaveRoom(value)
		//int解包

	case MyProtocol.MatchCode.Ready_CREQ: //发来的准备请求
		value := MyProtocol.UpackIntger(src)
		p.Ready(value)
		//int解包

	case MyProtocol.MatchCode.UnReady_CREQ: //取消准备的请求
		value := MyProtocol.UpackIntger(src)
		p.UnReady(value)
		//int解包

	}

}
