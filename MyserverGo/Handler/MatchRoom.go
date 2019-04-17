/*
	匹配房间
*/
package Handler

import (
	"bytes"
	"container/list"
	List "container/list"
	"fmt"

	"MyserverGo/MyProtocol"
	QueueI "MyserverGo/container"
)

type MatchRoom struct {
	/// 房间ID，唯一标识
	RommId int

	/// 房间内的玩家
	///元素为PlayerData
	clientList List.List

	/// 房间内准备的玩家ID列表 int
	readyUIdList List.List
}

//重用房间队列
var RoomQueue *QueueI.Queue

//正在匹配的用户ID与房间ID的映射字典
var UserIdRoomIdDic = make(map[int]int)

/// 正在匹配的房间ID与之对应的房间数据模型之间的映射字典
var RoomIdModelDic = make(map[int]MatchRoom)

//初始化匹配房间资源
func InitRoomQueue() {
	//fmt.Println("*****************************初始化RoomQueue")
	RoomQueue = QueueI.NewQueue(0)
	Number = 0
}

//房间ID
var Number int

/// 进入匹配房间
func (p *PlayerData) Enter() MatchRoom {

	//先遍历正在匹配的房间数据模型字典中有没有未满的房间，如果有，加进去
	for _, value := range RoomIdModelDic {
		if value.IsEmpty() {
			continue
		}
		value.Enter(p)

		UserIdRoomIdDic[p.State] = value.RommId
	}
	//如果执行到这里，代表正在匹配的房间数据模型字典中没有空位了，自己开一间房
	//room := MatchRoom{}

	if RoomQueue.Length() > 0 {
		t := RoomQueue.Dequeue()
		room := anyclisQueueOfMatchRoom(t)
		room.Enter(p)
		UserIdRoomIdDic[p.State] = room.RommId
		RoomIdModelDic[room.RommId] = room

		return room
	} else {
		room := MatchRoom{}
		room.RommId = Number
		Number++
		room.Enter(p)
		UserIdRoomIdDic[p.State] = room.RommId
		RoomIdModelDic[room.RommId] = room
		return room
	}

}

//离开房间
func (p *PlayerData) Leave() MatchRoom {

	// int roomId = userIdRoomIdDic[userId];
	// MatchRoom room = roomIdModelDic[roomId];
	// room.Leave(DatabaseManager.GetClientPeerByUserId(userId));
	// userIdRoomIdDic.Remove(userId);

	roomId := UserIdRoomIdDic[p.State]
	room := RoomIdModelDic[roomId]
	room.Leave(p)
	delete(UserIdRoomIdDic, p.State)

	//如果房间为空，将房间加入到房间重用队列，从正在匹配的房间字典中移除掉
	if room.IsEmpty() {
		delete(RoomIdModelDic, roomId)
		RoomQueue.Push(room)
	}

	return room

}

//*************************************MatchRoom***********//
/// 获取房间是否为空
func (p *MatchRoom) IsEmpty() bool {
	return p.clientList.Len() == 0
}

/// 进入房间
func (p *MatchRoom) Enter(player *PlayerData) {
	p.clientList.PushBack(player)
}

//离开房间
func (p *MatchRoom) Leave(player *PlayerData) {

	for e := p.clientList.Front(); e != nil; e = e.Next() {
		flag := e.Value.(PlayerData)
		if flag.Name == player.Name {
			p.clientList.Remove(e)
			break
		}
	}
}

/// 销毁房间 游戏开始时调用
func (p *MatchRoom) DestoryRoom() {

	delete(RoomIdModelDic, p.RommId)

	for e := p.clientList.Front(); e != nil; e = e.Next() {
		player := e.Value.(PlayerData)
		delete(UserIdRoomIdDic, player.State)
	}

	//清空list
	var next *list.Element
	for e := p.clientList.Front(); e != nil; e = next {
		next = e.Next()
		p.clientList.Remove(e)
	}

	var next2 *list.Element
	for e := p.readyUIdList.Front(); e != nil; e = next2 {
		next2 = e.Next()
		p.readyUIdList.Remove(e)
	}

	RoomQueue.Push(p)
}

//准备
func (p *MatchRoom) Ready(userid int) {
	p.readyUIdList.PushBack(userid)
}

//取消准备
func (p *MatchRoom) UnReady(userid int) {
	for e := p.readyUIdList.Front(); e != nil; e = e.Next() {
		id := e.Value.(int)
		if id == userid {
			p.readyUIdList.Remove(e)
			break
		}
	}
}

//获取是否全部玩家准备，如果返回值为True，既可以开始游戏了
func (p *MatchRoom) IsAllReady() bool {
	return p.readyUIdList.Len() == 3
}

//广播发消息
//Broadcast 和 SendMSg 后期可以优化为接口！！ wangfw
func (p *MatchRoom) Broadcast(opcode uint32, subcode uint32, value interface{}, player *PlayerData) {
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

	for e := p.clientList.Front(); e != nil; e = e.Next() {
		if player == nil {
			continue
		}
		playernow := e.Value.(PlayerData)
		i, err := playernow.Conn.Write(packet)

		if err != nil {
			fmt.Println("SendMsg Error!")
			fmt.Println(err)
		}

		fmt.Println("发送的信息", i)
	}

}

/// 获取玩家所在的房间
func GetRoom(userid int) MatchRoom {
	roomid := UserIdRoomIdDic[userid]
	return RoomIdModelDic[roomid]
}

//*************************************解析interface***********//
//取出来是interface，解析queue的内容
//interface 里面是MatchRoom
func anyclisQueueOfMatchRoom(t interface{}) MatchRoom {

	if f, ok := t.(MatchRoom); ok {
		return MatchRoom{f.RommId, f.clientList, f.readyUIdList}
	}

	return MatchRoom{}
}

//*************************************UserIdRoomIdDic***********//

/// 是否在匹配房间里
func IsMatching(userid int) bool {
	_, ok := UserIdRoomIdDic[userid]

	return ok
}
