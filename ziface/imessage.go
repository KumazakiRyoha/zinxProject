package ziface

type IMessage interface {

	//获取消息的id
	GetMsgId() uint32
	//获取消息的长度
	GetMsgLen() uint32
	//获取消息的内容
	GetData() []byte
	//设置消息Id
	SetMsgId(uint32)
	//设置消息的内容
	SetData([]byte)
	//设置消息的长度
	SetDataLen(uint32)
}
