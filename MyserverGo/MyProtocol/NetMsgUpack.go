package MyProtocol

import (
	"encoding/binary"
	"fmt"
)

//解包
func Upack(src []byte) (uint32, uint32, []byte) {

	//获取opcode 和 subcode
	//opcode
	Opcode := binary.BigEndian.Uint32(src[0:4])
	src = src[4:]

	//获取subcode
	SubCode := binary.BigEndian.Uint32(src[0:4])
	src = src[4:]

	//ChooseOfValueOfOpcode(Opcode, SubCode, src, conn)

	return Opcode, SubCode, src
}

//账户模块信息解包
func UpackSubCodeAccountStr(src []byte) (string, string) {

	//myread := make([]byte, 1024)
	lenName := binary.BigEndian.Uint16(src[0:2])
	strName := string(src[2 : 2+lenName])
	src = src[2+lenName:]

	lenPas := binary.BigEndian.Uint16(src[0:2])
	strPassword := string(src[2 : 2+lenPas])
	src = src[2+lenPas:]

	fmt.Println(strName)
	fmt.Println(strPassword)

	return strName, strPassword

}

//int类型解包
func UpackIntger(src []byte) int {

	target := binary.BigEndian.Uint32(src[0:4])
	return int(target)

}

//string类型解包
//uint32的len+具体消息体
func UpackString(src []byte) string {

	len := binary.BigEndian.Uint32(src[0:4])
	src = src[4:]

	target := string(src[:len])

	return target
}
