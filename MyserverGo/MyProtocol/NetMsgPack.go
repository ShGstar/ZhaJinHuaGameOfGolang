package MyProtocol

import (
	"bytes"
	List "container/list"
	"encoding/binary"
	"fmt"
	"reflect"
)

//打包两个code
func PackCode(opcode uint32, subcode uint32) []byte {
	buf := new(bytes.Buffer)

	opcodebuf := make([]byte, 4)
	binary.BigEndian.PutUint32(opcodebuf, uint32(opcode))

	subcodebuf := make([]byte, 4)
	binary.BigEndian.PutUint32(subcodebuf, uint32(subcode))

	buf.Write(opcodebuf)
	buf.Write(subcodebuf)

	return buf.Bytes()
}

//打包消息体
func PackValue(opcode uint32, subcode uint32, t interface{}) []byte {

	//buf := new(bytes.Buffer) //可读写的缓存，用于存放数据，目前没用到
	packet := make([]byte, 4)

	//客户端获取用户信息的请求
	if (opcode == SubCode.Account && subcode == AccountCode.GetUserInfo_SRES) ||
		(opcode == SubCode.Math && subcode == MatchCode.Enter_BRO) {

		buf := new(bytes.Buffer)

		return PackDto(buf, t)
	} else if value, ok := t.(int); ok {
		binary.BigEndian.PutUint32(packet, uint32(value))
		//buf.Write([]byte{4})
		//binary.Write(buf, binary.BigEndian, value)
	} else if opcode == SubCode.Math && subcode == MatchCode.Enter_SRES {

		value := reflect.ValueOf(t)

		buf := new(bytes.Buffer)

		//格式：len+userIdUserDtoDicBuf+ len+enterOrderUserIdListBuf+ len+readyUIdListBuf

		userIdUserDtoDic := value.Index(0).Interface().(map[int]PlayerDataOfDto)
		userIdUserDtoDicBuf := PackMap(userIdUserDtoDic)
		lenuserIdUserDtoDicBuf := len(userIdUserDtoDicBuf)

		binary.Write(buf, binary.BigEndian, uint32(lenuserIdUserDtoDicBuf))
		buf.Write(userIdUserDtoDicBuf)

		enterOrderUserIdList := value.Index(1).Interface().(List.List)
		enterOrderUserIdListBuf := PackList(enterOrderUserIdList)
		lenenterOrderUserIdListBuf := len(enterOrderUserIdListBuf)

		binary.Write(buf, binary.BigEndian, uint32(lenenterOrderUserIdListBuf))
		buf.Write(enterOrderUserIdListBuf)

		readyUIdList := value.Index(2).Interface().(List.List)
		readyUIdListBuf := PackList(readyUIdList)
		lenreadyUIdListBuf := len(readyUIdListBuf)

		binary.Write(buf, binary.BigEndian, uint32(lenreadyUIdListBuf))
		buf.Write(readyUIdListBuf)

		return buf.Bytes()
	} else {
		fmt.Println("PackValue 断言err")
	}

	return packet
}

//打包最终的包
//构造包头和包体
func PackMsg(value []byte) []byte {

	length := len(value)

	//包头
	lengbuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengbuf, uint32(length))

	var packet bytes.Buffer
	packet.Write(lengbuf)
	packet.Write(value)

	return packet.Bytes()
}

//打包UserDto
func PackDto(buf *bytes.Buffer, t interface{}) []byte {

	//len+value格式（id + len+password+ len+InconName+ jinbi）!!!格式
	value := reflect.ValueOf(t)

	idbuf := make([]byte, 4)
	id := value.Index(0).Interface().(int)
	binary.BigEndian.PutUint32(idbuf, uint32(id))

	namebuf := make([]byte, 4)
	name := value.Index(1).Interface().(string)
	lenN := len(name)
	binary.BigEndian.PutUint32(namebuf, uint32(lenN))

	IconNamebuf := make([]byte, 4)
	InconName := value.Index(2).Interface().(string)
	lenI := len(InconName)
	binary.BigEndian.PutUint32(IconNamebuf, uint32(lenI))

	jinbi := value.Index(3).Interface().(int)
	//长度也要打包 最终格式为 len+value格式（id + len+password+ len+InconName+ jinbi）
	//binary.Write(buf, binary.BigEndian, uint32(id))
	buf.Write(idbuf)

	buf.Write(namebuf)
	buf.Write([]byte(name))

	buf.Write(IconNamebuf)
	buf.Write([]byte(InconName))

	binary.Write(buf, binary.BigEndian, uint32(jinbi))

	return buf.Bytes()

}

//打包clientList（List类型）元素值为int
//打包格式为 包长度+元素值
func PackList(listDemo List.List) []byte {

	buf := new(bytes.Buffer)

	for i := listDemo.Front(); i != nil; i = i.Next() {
		playerid := i.Value.(int)
		bufid := make([]byte, 4)
		binary.BigEndian.PutUint32(bufid, uint32(playerid))
		buf.Write(bufid)
	}
	packet := buf.Bytes()

	lenth := len(packet)

	buf2 := new(bytes.Buffer)
	binary.Write(buf2, binary.BigEndian, uint32(lenth))
	buf2.Write(packet)

	return buf2.Bytes()
}

//打包userIdUserDtoDic(map类型)
func PackMap(userIdUserDtoDic map[int]PlayerDataOfDto) []byte {

	buf := new(bytes.Buffer)

	lenmap := len(userIdUserDtoDic)
	binary.Write(buf, binary.BigEndian, uint32(lenmap))

	for _, value := range userIdUserDtoDic {

		//UserDto
		//t := []interface{}{value.State, value.Name, value.IconName, value.Coin}
		bufuser := new(bytes.Buffer)
		userDtobuf := PackDto(bufuser, []interface{}{value.State, value.Name, value.IconName, value.Coin})

		//UserDto 得长度
		UserDtolenbuf := make([]byte, 4)
		lenth := len(userDtobuf)
		binary.BigEndian.PutUint32(UserDtolenbuf, uint32(lenth))

		buf.Write(UserDtolenbuf)
		buf.Write(userDtobuf)

	}

	return buf.Bytes()
}
