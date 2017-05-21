package protocol

// 自定义数据类型定义段

//Foo Foo
type Foo struct {
}

// 回调函数类型定义段

//CallbackSayHello aaa
type CallbackSayHello func(msg string)

//回调函数定义段

//OnSayHello aaa
var OnSayHello CallbackSayHello

//解包响应函数实现段

//ParseSayHello aaa
func ParseSayHello(data []byte) {
	if OnSayHello == nil {
		// error code
		return
	}
	msg := string(data)
	OnSayHello(msg)
}

//代理函数实现段

func SayHello(h interface{}, msg string) {
	// data := []byte(msg)
}

func BroadcastSayHello(hlist []interface{}, msg string) {

}

//消息映射段
const (
	MessageSayHello = 1000
)

type CallbackParsePacket func([]byte)

var messageMap map[int]CallbackParsePacket

func init() {
	messageMap[MessageSayHello] = ParseSayHello
}

//消息解析段

func ParsePacket(data []byte) {
	msgid := 0
	_, ok := messageMap[msgid]
	if !ok {

	}
}
