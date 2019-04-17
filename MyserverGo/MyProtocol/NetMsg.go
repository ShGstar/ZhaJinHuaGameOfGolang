package MyProtocol

///网络消息结构
///每次发送消息都发送这个结构体
type NetMsg struct {

	///操作码
	//账号模块 -0
	//模块匹配 -1
	//聊天模块 -2
	//战斗模块 -3
	opCode uint32

	///子操作码
	subCode uint32

	///传递的参数
	value interface{}
}

//操作码验证
type SubCodeType struct {

	//账号模块 -0
	Account uint32

	//模块匹配 -1
	Math uint32

	//聊天战斗-2
	Chat uint32

	//战斗模块-3
	Fight uint32
}

//子操作码-账号模块 10个
type AccountCodeType struct {
	Register_CREQ    uint32
	Register_SRES    uint32
	Login_CREQ       uint32
	Login_SRES       uint32
	GetUserInfo_CREQ uint32
	GetUserInfo_SRES uint32
	GetRankList_CREQ uint32
	GetRankList_SRES uint32

	UpdateCoinCount_CREQ uint32
	UpdateCoinCount_SRES uint32
}

//子操作码-聊天战斗
type ChatCodeType struct {
	CREQ uint32
	BRO  uint32
}

//子操作码-战斗模块 -14ge
type FightCodeType struct {

	/// 开始战斗的广播
	StartFight_BRO uint32
	Leave_CREQ     uint32
	Leave_BRO      uint32

	/// 开始下注的服务器广播
	StartStakes_BRO uint32

	/// 看牌的客户端请求
	LookCard_CREQ uint32
	LookCard_BRO  uint32
	Follow_CREQ   uint32
	PutStakes_BRO uint32

	/// 加注的客户端请求
	AddStakes_CREQ uint32

	/// 弃牌的客户端请求
	GiveUpCard_CREQ  uint32
	GiveUpCard_BRO   uint32
	CompareCard_CREQ uint32
	CompareCard_BRO  uint32
	GameOver_BRO     uint32
}

//子操作码-匹配模块10ge
type MatchCodeType struct {
	Enter_CREQ uint32
	Enter_SRES uint32
	Enter_BRO  uint32

	Leave_CREQ uint32
	Leave_BRO  uint32

	Ready_CREQ   uint32
	Ready_BRO    uint32
	UnReady_CREQ uint32
	UnReady_BRO  uint32

	StartGame_BRO uint32
}

type PlayerDataOfDto struct {
	State    int //状态
	Name     string
	IconName string
	Coin     int //jinbi
}
