package Handler

import (
	"MyserverGo/MyProtocol"
	"reflect"
)

/// <summary>
/// 客户端进入房间的请求
/// </summary>
/// value为房间列表id
//！！此处应该单线程处理逻辑，需要枷锁
func (p *PlayerData) EnterRoom(value int) {

	//判断一下当前的客户端链接对象是不是在匹配房间里面，如果在，则忽略
	if IsMatching(p.State) {
		return
	}

	room := p.Enter()

	//构造UserDto用户数据传输模型
	//广播给房间内的所有玩家，除了自身，有新的玩家进来了，参数：新进用户的UserDto
	// room.Broadcast(OpCode.Match, MatchCode.Enter_BRO, userDto, client);
	room.Broadcast(MyProtocol.SubCode.Math, MyProtocol.MatchCode.Enter_BRO,
		[]interface{}{p.State, p.Name, p.IconName, p.Coin}, p)

	//给客户端一个相应 参数：房间传输模型，包含房间内的正在等待的玩家以及准备的玩家id集合
	//建立传输模型
	/// 进入房间顺序的用户ID列表
	enterOrderUserIdList := room.readyUIdList

	/// 用户ID与该用户UserDto之间的映射字典
	userIdUserDtoDic := CreateUserMap(room)

	p.Sendmsg(MyProtocol.SubCode.Math, MyProtocol.MatchCode.Enter_SRES,
		[]interface{}{userIdUserDtoDic, enterOrderUserIdList, room.readyUIdList})

}

//客户端离开房间请求
//！！此处应该单线程处理逻辑，需要枷锁
func (p *PlayerData) LeaveRoom(value int) {

	//不在匹配房间 忽略
	if IsMatching(p.State) == false {
		return
	}

	//p.Leave()
	room := p.Leave()

	//广播消息
	// room.Broadcast(OpCode.Match, MatchCode.Leave_BRO, client.Id);
	room.Broadcast(MyProtocol.SubCode.Math, MyProtocol.MatchCode.Leave_BRO, p.State, nil)
}

//客户端发来得准备请求
//！！此处应该单线程处理逻辑，需要枷锁
func (p *PlayerData) Ready(value int) {

	//不在匹配房间 忽略
	if IsMatching(p.State) == false {
		return
	}

	room := GetRoom(p.State)
	room.Ready(p.State)

	//广播
	//room.Broadcast(OpCode.Match, MatchCode.Ready_BRO, client.Id);
	room.Broadcast(MyProtocol.SubCode.Math, MyProtocol.MatchCode.Ready_BRO, p.State, nil)

	//全部都准备了，可以开始游戏了
	if room.IsAllReady() {
		//战斗
		//startFight(room.clientList, roomType)

		//通知房间内的所有玩家，开始游戏了
		room.Broadcast(MyProtocol.SubCode.Math, MyProtocol.MatchCode.StartGame_BRO, nil, nil)
		//销毁房间
		room.DestoryRoom()
	}

}

//取消准备
//！！此处应该单线程处理逻辑，需要枷锁
func (p *PlayerData) UnReady(value int) {

	//不在匹配房间 忽略
	if IsMatching(p.State) == false {
		return
	}

	room := GetRoom(p.State)
	room.UnReady(p.State)
	// room.Broadcast(OpCode.Match, MatchCode.UnReady_BRO, client.Id);
	room.Broadcast(MyProtocol.SubCode.Math, MyProtocol.MatchCode.UnReady_BRO, p.State, nil)
}

///构造用户ID与该用户UserDto之间的映射字典 模型
func CreateUserMap(p MatchRoom) map[int]MyProtocol.PlayerDataOfDto {

	userIdUserDtoDic := make(map[int]MyProtocol.PlayerDataOfDto)

	for i := p.clientList.Front(); i != nil; i = i.Next() {
		t := i.Value
		x := reflect.ValueOf(t).Interface()

		player := x.(*PlayerData)

		playerDto := MyProtocol.PlayerDataOfDto{}

		playerDto.Name = player.Name
		playerDto.IconName = player.IconName
		playerDto.Coin = player.Coin
		playerDto.State = player.State

		userIdUserDtoDic[playerDto.State] = playerDto
	}

	return userIdUserDtoDic
}
