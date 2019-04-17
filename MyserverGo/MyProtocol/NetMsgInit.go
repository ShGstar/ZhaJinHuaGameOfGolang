package MyProtocol

//操作码和子操作码的初始化

//全局变量 各种操作码
var (
	SubCode     SubCodeType
	AccountCode AccountCodeType
	FightCode   FightCodeType
	ChatCode    ChatCodeType
	MatchCode   MatchCodeType
)

//初始化各种操作码
func init() {
	SubCode = InitOfSubCode()
	AccountCode = InitOfAccountCode()
	FightCode = InitOfFightCode()
	ChatCode = InitOfChatCode()
	MatchCode = InitOfMatchCode()
}

//操作码
func InitOfSubCode() SubCodeType {
	SubCode := SubCodeType{0, 1, 2, 3}

	return SubCode
}

//子操作码-账号模块
func InitOfAccountCode() AccountCodeType {
	AccountCode := AccountCodeType{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	return AccountCode
}

//子操作码-聊天战斗
func InitOfChatCode() ChatCodeType {

	ChatCode := ChatCodeType{0, 1}

	return ChatCode
}

//子操作码-战斗模块
func InitOfFightCode() FightCodeType {
	FightCode := FightCodeType{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}

	return FightCode
}

//子操作码-匹配模块
func InitOfMatchCode() MatchCodeType {
	MatchCode := MatchCodeType{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	return MatchCode
}
