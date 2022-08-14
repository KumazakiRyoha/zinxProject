package znet

/**
将请求的消息封装到一个Message中，定义抽象的接口
*/
type Message struct {
	Id      uint32 // 消息UD
	DataLen uint32 // 消息的长度
	Data    []byte // 消息内容
}

//获取消息的id
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

//获取消息的长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

//获取消息的内容
func (m *Message) GetData() []byte {
	return m.Data
}

//设置消息Id
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

//设置消息的内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}

//设置消息的长度
func (m *Message) SetDataLen(datalen uint32) {
	m.DataLen = datalen
}
