package Evemont

import (
	net "net"

	QueueI "MyserverGo/container"
)

/*
	@author wangfw
	采用gdb-go基础库的queue
	线程安全的
*/

type ClientTag struct {
	Conn  *net.TCPConn
	State int //状态
}

//玩家
//玩家
type PlayerData struct {
	Conn     *net.TCPConn
	Name     string
	Password string
	State    int //状态
	IconName string
	Coin     int //jinbi
}

//type ClientLoginSuccsess struct {
//}

//玩家连接，指连接成功并且登陆成功的
var Player *QueueI.Queue

//玩家连接,指连接成功的
var ClientNum *QueueI.Queue

//初始化环境资源
func InitEvE() {
	ClientNum = QueueI.NewQueue(0)
	Player = QueueI.NewQueue(0)
	//ClientNum.Push(1)
}

//建立一个玩家对象模型
func InitPlaterData(name string, password string, conn *net.TCPConn, state int) PlayerData {
	//playernum := PlayerData{conn, name, password, state}
	playernum := PlayerData{}
	playernum.Conn = conn
	playernum.Name = name
	playernum.Password = password
	playernum.State = state

	return playernum
}

//将玩家一个玩家对象模型插入到玩家连接池中
func (p *PlayerData) AddPlayer() {
	Player.Push(p)
}

//建立一个玩家对象模型 并 插入到玩家连接池中
func InitAndAddPlayer(name string, password string, conn *net.TCPConn, state int) PlayerData {
	Player := InitPlaterData(name, password, conn, state)
	Player.AddPlayer()
	return Player
}

//取出来是interface，解析queue的内容
//interface 里面是PlayerData
func anyclisQueueOfPlayerData(t interface{}) (PlayerData, bool) {

	if f, ok := t.(PlayerData); ok {
		return PlayerData{f.Conn, f.Name, f.Password, f.State, f.IconName, f.Coin}, ok
	}

	return PlayerData{}, false
}

//获取队列第一个元素
func GetFirstElem() (PlayerData, bool) {

	t := Player.GetFirst()

	return anyclisQueueOfPlayerData(t)
}
